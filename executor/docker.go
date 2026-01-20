package executor

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/pfarrer/foghorn/config"
	"github.com/pfarrer/foghorn/scheduler"
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

	return &DockerExecutor{
		cli:            cli,
		defaultTimeout: 30 * time.Second,
		outputLocation: "stdout",
	}, nil
}

func (e *DockerExecutor) Execute(check scheduler.CheckConfig) error {
	adapter, ok := check.(*scheduler.ConfigAdapter)
	if !ok {
		return fmt.Errorf("invalid check config type")
	}

	checkConfig := adapter.Config

	timeout := e.defaultTimeout
	if checkConfig.Timeout != "" {
		parsedTimeout, err := time.ParseDuration(checkConfig.Timeout)
		if err == nil {
			timeout = parsedTimeout
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	env := e.buildEnvVars(checkConfig)

	containerConfig := &container.Config{
		Image: checkConfig.Image,
		Env:   env,
	}

	hostConfig := &container.HostConfig{
		AutoRemove: false,
	}

	resp, err := e.cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, "")
	if err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}
	defer e.cli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{Force: true})

	if err := e.cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	statusCh, errCh := e.cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)

	select {
	case statusResult := <-statusCh:
		if statusResult.StatusCode != 0 {
			return fmt.Errorf("check failed with exit code %d", statusResult.StatusCode)
		}
		_, err := e.readResult(ctx, resp.ID)
		if err != nil {
			return fmt.Errorf("failed to read check result: %w", err)
		}
		return nil
	case err := <-errCh:
		return fmt.Errorf("error waiting for container: %w", err)
	case <-ctx.Done():
		e.cli.ContainerKill(ctx, resp.ID, "SIGKILL")
		return fmt.Errorf("check execution timed out after %v", timeout)
	}
}

func (e *DockerExecutor) buildEnvVars(check *config.CheckConfig) []string {
	env := []string{
		fmt.Sprintf("FOGHORN_CHECK_NAME=%s", check.Name),
	}

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

	for k, v := range check.Env {
		if strings.HasPrefix(k, "SECRET_") {
			secrets := map[string]string{k: v}
			secretsJSON, _ := json.Marshal(secrets)
			env = append(env, fmt.Sprintf("FOGHORN_SECRETS=%s", string(secretsJSON)))
			break
		}
	}

	for k, v := range check.Env {
		if !strings.HasPrefix(k, "FOGHORN_") && k != "ENDPOINT" && !strings.HasPrefix(k, "SECRET_") {
			env = append(env, fmt.Sprintf("%s=%s", k, v))
		}
	}

	return env
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

func (e *DockerExecutor) readResult(ctx context.Context, containerID string) (*CheckResult, error) {
	reader, err := e.cli.ContainerLogs(ctx, containerID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: false,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to read container logs: %w", err)
	}
	defer reader.Close()

	logs, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read logs: %w", err)
	}

	logStr := string(demultiplexLogs(logs))
	logStr = strings.TrimSpace(logStr)

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
