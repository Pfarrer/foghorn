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
		AutoRemove: true,
	}

	resp, err := e.cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, "")
	if err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}

	defer e.cli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{Force: true})

	startTime := time.Now()

	if err := e.cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	statusCh, errCh := e.cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)

	select {
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("error waiting for container: %w", err)
		}
	case <-ctx.Done():
		e.cli.ContainerKill(ctx, resp.ID, "SIGKILL")
		return fmt.Errorf("check execution timed out after %v", timeout)
	}

	var status container.WaitResponse
	select {
	case status = <-statusCh:
	case <-ctx.Done():
		return fmt.Errorf("check execution timed out after %v", timeout)
	}
	duration := time.Since(startTime)

	if status.StatusCode != 0 {
		return fmt.Errorf("check failed with exit code %d", status.StatusCode)
	}

	result, err := e.readResult(ctx, resp.ID)
	if err != nil {
		return fmt.Errorf("failed to read check result: %w", err)
	}

	e.printResult(result, duration)

	return nil
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

	logStr := string(logs)
	logStr = strings.TrimPrefix(logStr, "\x00\x00")
	logStr = strings.TrimSuffix(logStr, "\n")

	var result CheckResult
	if err := json.Unmarshal([]byte(logStr), &result); err != nil {
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

func (e *DockerExecutor) printResult(result *CheckResult, duration time.Duration) {
	fmt.Printf("Status: %s\n", result.Status)
	fmt.Printf("Message: %s\n", result.Message)
	if result.Data != nil {
		fmt.Printf("Data: %v\n", result.Data)
	}
	fmt.Printf("Duration: %v\n", duration.Round(time.Millisecond))
}

func (e *DockerExecutor) Close() error {
	if e.cli != nil {
		return e.cli.Close()
	}
	return nil
}
