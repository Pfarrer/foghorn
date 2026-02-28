package executor

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/pfarrer/foghorn/config"
	"github.com/pfarrer/foghorn/imageresolver"
	"github.com/pfarrer/foghorn/logger"
	"github.com/pfarrer/foghorn/scheduler"
	"github.com/pfarrer/foghorn/secretstore"
)

type CheckResult struct {
	Status     string                 `json:"status"`
	Message    string                 `json:"message"`
	Data       map[string]interface{} `json:"data,omitempty"`
	Timestamp  string                 `json:"timestamp"`
	DurationMs int64                  `json:"duration_ms"`
}

type DockerExecutor struct {
	cli            *client.Client
	defaultTimeout time.Duration
	outputLocation string
	resultCallback func(checkName string, status string, duration time.Duration)
	resolveMu      sync.Mutex
	resolvedImages map[string]string
	secretResolver SecretResolver
	secretBaseDir  string
	debugOutput    string
	debugMaxChars  int
}

const (
	debugOutputModeOff       = "off"
	debugOutputModeOnFailure = "on_failure"
	debugOutputModeAlways    = "always"
	defaultDebugOutputMode   = debugOutputModeOff
	defaultDebugOutputMax    = 4096
)

type SecretResolver interface {
	Resolve(ref string) (string, error)
}

type ExecuteOptions struct {
	Timeout        time.Duration
	OutputLocation string
}

func NewDockerExecutor() (*DockerExecutor, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %w", err)
	}

	secretBaseDir := filepath.Join(os.TempDir(), "foghorn-secrets")
	if err := os.MkdirAll(secretBaseDir, 0o700); err != nil {
		return nil, fmt.Errorf("failed to create secret base directory: %w", err)
	}

	if err := cleanupOldSecretDirs(secretBaseDir); err != nil {
		logger.Warn("Failed to cleanup old secret directories: %v", err)
	}

	return &DockerExecutor{
		cli:            cli,
		defaultTimeout: 30 * time.Second,
		outputLocation: "stdout",
		resolvedImages: make(map[string]string),
		secretBaseDir:  secretBaseDir,
		debugOutput:    defaultDebugOutputMode,
		debugMaxChars:  defaultDebugOutputMax,
	}, nil
}

func (e *DockerExecutor) Execute(check scheduler.CheckConfig) error {
	adapter, ok := check.(*scheduler.ConfigAdapter)
	if !ok {
		return fmt.Errorf("invalid check config type")
	}

	checkConfig := adapter.Config
	checkName := checkConfig.Name

	timeout := e.defaultTimeout
	if checkConfig.Timeout != "" {
		parsedTimeout, err := time.ParseDuration(checkConfig.Timeout)
		if err == nil {
			timeout = parsedTimeout
		}
	}

	startTime := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	image, err := e.resolveImage(ctx, checkConfig.Image)
	if err != nil {
		duration := time.Since(startTime)
		if e.resultCallback != nil {
			e.resultCallback(checkName, "error", duration)
		}
		logger.Error("Check %s: Failed to resolve image: %v", checkName, err)
		return err
	}
	if err := e.ensureImageAvailable(ctx, image, checkName); err != nil {
		duration := time.Since(startTime)
		if e.resultCallback != nil {
			e.resultCallback(checkName, "error", duration)
		}
		logger.Error("Check %s: Failed to prepare image: %v", checkName, err)
		return err
	}

	logger.Debug("Check %s: Creating container with image %s (timeout: %v)", checkName, image, timeout)

	env, secretDir, secretsToRedact, err := e.buildEnvVars(checkConfig)
	if err != nil {
		duration := time.Since(startTime)
		if e.resultCallback != nil {
			e.resultCallback(checkName, "error", duration)
		}
		logger.Error("Check %s: Failed to prepare environment: %v", checkName, err)
		return err
	}
	if secretDir != "" {
		defer cleanupSecretDir(secretDir)
	}

	debugMode := normalizeDebugOutputMode(checkConfig.DebugOutput)
	if debugMode == "" {
		debugMode = e.debugOutput
	}

	containerConfig := &container.Config{
		Image: image,
		Env:   env,
	}

	hostConfig := &container.HostConfig{
		AutoRemove: false,
	}
	if secretDir != "" {
		hostConfig.Binds = append(hostConfig.Binds, fmt.Sprintf("%s:/run/foghorn/secrets:ro", secretDir))
	}

	resp, err := e.cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, "")
	if err != nil {
		duration := time.Since(startTime)
		if e.resultCallback != nil {
			e.resultCallback(checkName, "error", duration)
		}
		logger.Error("Check %s: Failed to create container: %v", checkName, err)
		return fmt.Errorf("failed to create container: %w", err)
	}
	logger.Debug("Check %s: Container created (ID: %s)", checkName, resp.ID)
	defer e.cli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{Force: true})

	if err := e.cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		duration := time.Since(startTime)
		if e.resultCallback != nil {
			e.resultCallback(checkName, "error", duration)
		}
		logger.Error("Check %s: Failed to start container: %v", checkName, err)
		return fmt.Errorf("failed to start container: %w", err)
	}
	logger.Debug("Check %s: Container started (ID: %s)", checkName, resp.ID)

	statusCh, errCh := e.cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)

	select {
	case statusResult := <-statusCh:
		if statusResult.StatusCode != 0 {
			if shouldLogContainerDebugOutput(debugMode, true) {
				if err := e.logContainerDebugOutput(checkName, resp.ID, "failure", secretsToRedact); err != nil {
					logger.Debug("Check %s: Failed to read container output after failure: %v", checkName, err)
				}
			}

			duration := time.Since(startTime)
			if e.resultCallback != nil {
				e.resultCallback(checkName, "error", duration)
			}
			logger.Error("Check %s: Failed with exit code %d", checkName, statusResult.StatusCode)
			return fmt.Errorf("check failed with exit code %d", statusResult.StatusCode)
		}
		result, err := e.readResult(ctx, resp.ID)
		if err != nil {
			duration := time.Since(startTime)
			if e.resultCallback != nil {
				e.resultCallback(checkName, "error", duration)
			}
			logger.Error("Check %s: Failed to read result: %v", checkName, err)
			return fmt.Errorf("failed to read check result: %w", err)
		}
		duration := time.Since(startTime)
		if e.resultCallback != nil {
			e.resultCallback(checkName, result.Status, duration)
		}
		if shouldLogContainerDebugOutput(debugMode, false) {
			if err := e.logContainerDebugOutput(checkName, resp.ID, "success", secretsToRedact); err != nil {
				logger.Debug("Check %s: Failed to read container output after success: %v", checkName, err)
			}
		}
		logger.Info("Check %s: Completed with status %s (duration: %dms) - %s", checkName, result.Status, result.DurationMs, result.Message)
		return nil
	case err := <-errCh:
		duration := time.Since(startTime)
		if e.resultCallback != nil {
			e.resultCallback(checkName, "error", duration)
		}
		logger.Error("Check %s: Error waiting for container: %v", checkName, err)
		return fmt.Errorf("error waiting for container: %w", err)
	case <-ctx.Done():
		duration := time.Since(startTime)
		if e.resultCallback != nil {
			e.resultCallback(checkName, "error", duration)
		}
		logger.Warn("Check %s: Execution timed out after %v", checkName, timeout)
		e.cli.ContainerKill(ctx, resp.ID, "SIGKILL")
		return fmt.Errorf("check execution timed out after %v", timeout)
	}
}

func (e *DockerExecutor) buildEnvVars(check *config.CheckConfig) ([]string, string, []string, error) {
	env := []string{
		fmt.Sprintf("FOGHORN_CHECK_NAME=%s", check.Name),
	}
	secretsToRedact := make([]string, 0)

	if check.Metadata != nil {
		configJSON, err := json.Marshal(check.Metadata)
		if err == nil {
			env = append(env, fmt.Sprintf("FOGHORN_CHECK_CONFIG=%s", string(configJSON)))
		}
	}

	if endpoint, ok := check.Env["ENDPOINT"]; ok {
		env = append(env, fmt.Sprintf("FOGHORN_ENDPOINT=%s", endpoint))
		env = append(env, fmt.Sprintf("ENDPOINT=%s", endpoint))
	}

	if timeout := check.Timeout; timeout != "" {
		env = append(env, fmt.Sprintf("FOGHORN_TIMEOUT=%s", timeout))
	}

	secretDir := ""
	for k, v := range check.Env {
		if refKey, ok := secretstore.ParseRef(v); ok {
			if e.secretResolver == nil {
				return nil, "", nil, fmt.Errorf("check %s requires secret %q, but secret store is not configured", check.Name, refKey)
			}

			secretValue, err := e.secretResolver.Resolve(v)
			if err != nil {
				return nil, "", nil, fmt.Errorf("check %s failed resolving secret: %w", check.Name, err)
			}
			if secretValue == "" {
				return nil, "", nil, fmt.Errorf("check %s secret %q resolved to an empty value", check.Name, refKey)
			}
			if secretDir == "" {
				dir, err := e.createSecretDir()
				if err != nil {
					return nil, "", nil, fmt.Errorf("failed to create temporary secret directory: %w", err)
				}
				secretDir = dir
			}

			filename := sanitizeSecretFilename(k)
			secretPath := filepath.Join(secretDir, filename)
			if err := os.WriteFile(secretPath, []byte(secretValue), 0o644); err != nil {
				return nil, "", nil, fmt.Errorf("failed to write secret file for %s: %w", k, err)
			}
			logger.Debug("Check %s: Injected secret reference %q into %s_FILE", check.Name, refKey, k)
			env = append(env, fmt.Sprintf("%s_FILE=/run/foghorn/secrets/%s", k, filename))
			secretsToRedact = append(secretsToRedact, secretValue)
		}
	}

	for k, v := range check.Env {
		if _, ok := secretstore.ParseRef(v); ok {
			continue
		}
		if !strings.HasPrefix(k, "FOGHORN_") && k != "ENDPOINT" {
			env = append(env, fmt.Sprintf("%s=%s", k, v))
		}
	}

	return env, secretDir, secretsToRedact, nil
}

func demultiplexLogs(data []byte) []byte {
	var result []byte
	for len(data) >= 8 {
		streamType := data[0]
		frameSize := int(data[4])<<24 | int(data[5])<<16 | int(data[6])<<8 | int(data[7])
		data = data[8:]
		if frameSize <= len(data) {
			result = append(result, data[:frameSize]...)
			data = data[frameSize:]
		} else {
			break
		}
		_ = streamType
	}
	return result
}

func truncateLogOutput(output string, maxChars int) string {
	if maxChars <= 0 || len(output) <= maxChars {
		return output
	}
	return "... (truncated, showing tail)\n" + output[len(output)-maxChars:]
}

func normalizeDebugOutputMode(mode string) string {
	switch strings.TrimSpace(mode) {
	case debugOutputModeOff, debugOutputModeOnFailure, debugOutputModeAlways:
		return strings.TrimSpace(mode)
	default:
		return ""
	}
}

func shouldLogContainerDebugOutput(mode string, failed bool) bool {
	switch normalizeDebugOutputMode(mode) {
	case debugOutputModeAlways:
		return true
	case debugOutputModeOnFailure:
		return failed
	default:
		return false
	}
}

var (
	authHeaderPattern  = regexp.MustCompile(`(?im)(authorization\s*[:=]\s*)([^\r\n]+)`)
	credentialPattern  = regexp.MustCompile(`(?im)(\"?(?:password|passwd|token|secret|api[_-]?key|authorization)\"?\s*[:=]\s*)(\"[^\"]*\"|'[^']*'|[^\s,}]+)`)
	bearerTokenPattern = regexp.MustCompile(`(?i)\bbearer\s+[A-Za-z0-9\-._~+/]+=*`)
)

func redactContainerOutput(output string, secrets []string) string {
	redacted := output
	uniqueSecrets := make([]string, 0, len(secrets))
	for _, secret := range secrets {
		if secret == "" {
			continue
		}
		if !slices.Contains(uniqueSecrets, secret) {
			uniqueSecrets = append(uniqueSecrets, secret)
		}
	}
	for _, secret := range uniqueSecrets {
		redacted = strings.ReplaceAll(redacted, secret, "[REDACTED]")
	}
	redacted = authHeaderPattern.ReplaceAllString(redacted, "${1}[REDACTED]")
	redacted = credentialPattern.ReplaceAllString(redacted, "${1}[REDACTED]")
	redacted = bearerTokenPattern.ReplaceAllString(redacted, "Bearer [REDACTED]")
	return redacted
}

func (e *DockerExecutor) logContainerDebugOutput(checkName string, containerID string, reason string, secretsToRedact []string) error {
	debugCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	output, err := e.readContainerOutput(debugCtx, containerID, true, true)
	cancel()
	if err != nil {
		return err
	}
	if output == "" {
		logger.Debug("Check %s: Container output on %s was empty", checkName, reason)
		return nil
	}
	redacted := redactContainerOutput(output, secretsToRedact)
	logger.Debug("Check %s: Container output on %s:\n%s", checkName, reason, truncateLogOutput(redacted, e.debugMaxChars))
	return nil
}

func (e *DockerExecutor) readContainerOutput(ctx context.Context, containerID string, showStdout bool, showStderr bool) (string, error) {
	reader, err := e.cli.ContainerLogs(ctx, containerID, container.LogsOptions{
		ShowStdout: showStdout,
		ShowStderr: showStderr,
	})
	if err != nil {
		return "", fmt.Errorf("failed to read container logs: %w", err)
	}
	defer reader.Close()

	logs, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read logs: %w", err)
	}

	return strings.TrimSpace(string(demultiplexLogs(logs))), nil
}

func (e *DockerExecutor) readResult(ctx context.Context, containerID string) (*CheckResult, error) {
	logStr, err := e.readContainerOutput(ctx, containerID, true, false)
	if err != nil {
		return nil, err
	}

	var result CheckResult
	if err := json.Unmarshal([]byte(logStr), &result); err != nil {
		openBrace := strings.LastIndex(logStr, "{")
		if openBrace != -1 {
			if err := json.Unmarshal([]byte(logStr[openBrace:]), &result); err == nil {
				return &result, nil
			}
		}

		reader, _, err := e.cli.CopyFromContainer(ctx, containerID, "/output/result.json")
		if err == nil {
			defer reader.Close()

			fileContent, err := io.ReadAll(reader)
			if err == nil {
				if err := json.Unmarshal(fileContent, &result); err == nil {
					return &result, nil
				}
			}
		}

		return nil, fmt.Errorf("failed to parse JSON output: %w", err)
	}

	return &result, nil
}

func (e *DockerExecutor) Close() error {
	if e.cli != nil {
		return e.cli.Close()
	}
	return nil
}

func (e *DockerExecutor) SetResultCallback(callback func(checkName string, status string, duration time.Duration)) {
	e.resultCallback = callback
}

func (e *DockerExecutor) SetSecretResolver(resolver SecretResolver) {
	e.secretResolver = resolver
}

func (e *DockerExecutor) SetDebugOutput(mode string, maxChars int) {
	normalized := normalizeDebugOutputMode(mode)
	if normalized == "" {
		normalized = defaultDebugOutputMode
	}
	e.debugOutput = normalized
	if maxChars > 0 {
		e.debugMaxChars = maxChars
		return
	}
	e.debugMaxChars = defaultDebugOutputMax
}

var nonFileSafeChars = regexp.MustCompile(`[^A-Za-z0-9._-]`)

func sanitizeSecretFilename(input string) string {
	sanitized := nonFileSafeChars.ReplaceAllString(input, "_")
	if sanitized == "" {
		return "secret"
	}
	return sanitized
}

func (e *DockerExecutor) resolveImage(ctx context.Context, image string) (string, error) {
	e.resolveMu.Lock()
	if resolved, ok := e.resolvedImages[image]; ok {
		e.resolveMu.Unlock()
		return resolved, nil
	}
	e.resolveMu.Unlock()

	resolved, err := imageresolver.Resolve(ctx, e.cli, image)
	if err != nil {
		return "", err
	}

	e.resolveMu.Lock()
	e.resolvedImages[image] = resolved
	e.resolveMu.Unlock()

	return resolved, nil
}

func (e *DockerExecutor) createSecretDir() (string, error) {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random suffix: %w", err)
	}
	suffix := hex.EncodeToString(b)

	dir := filepath.Join(e.secretBaseDir, suffix)
	if err := os.Mkdir(dir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create secret directory: %w", err)
	}

	tsFile := filepath.Join(dir, ".timestamp")
	if err := os.WriteFile(tsFile, []byte(time.Now().Format(time.RFC3339)), 0o600); err != nil {
		_ = os.RemoveAll(dir)
		return "", fmt.Errorf("failed to write timestamp file: %w", err)
	}

	return dir, nil
}

func cleanupSecretDir(secretDir string) error {
	return os.RemoveAll(secretDir)
}

func cleanupOldSecretDirs(baseDir string) error {
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	cutoff := time.Now().Add(-24 * time.Hour)

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		dir := filepath.Join(baseDir, entry.Name())
		tsFile := filepath.Join(dir, ".timestamp")

		data, err := os.ReadFile(tsFile)
		if err != nil {
			continue
		}

		ts, err := time.Parse(time.RFC3339, string(data))
		if err != nil {
			continue
		}

		if ts.Before(cutoff) {
			if err := os.RemoveAll(dir); err == nil {
				logger.Debug("Cleanup old secret directory: %s", dir)
			}
		}
	}

	return nil
}

func (e *DockerExecutor) ensureImageAvailable(ctx context.Context, imageRef string, checkName string) error {
	_, _, err := e.cli.ImageInspectWithRaw(ctx, imageRef)
	if err == nil {
		return nil
	}
	if !client.IsErrNotFound(err) {
		return fmt.Errorf("failed to inspect image %s: %w", imageRef, err)
	}

	logger.Info("Check %s: Pulling image %s", checkName, imageRef)

	pullCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	reader, err := e.cli.ImagePull(pullCtx, imageRef, image.PullOptions{})
	if err != nil {
		return fmt.Errorf("failed to pull image %s: %w", imageRef, err)
	}
	defer reader.Close()

	if _, err := io.Copy(io.Discard, reader); err != nil {
		return fmt.Errorf("failed to complete pull for image %s: %w", imageRef, err)
	}

	return nil
}
