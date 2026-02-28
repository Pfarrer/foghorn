package executor

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/pfarrer/foghorn/config"
)

func TestCheckResultJSON(t *testing.T) {
	result := CheckResult{
		Status:     "pass",
		Message:    "Check passed successfully",
		Data:       map[string]interface{}{"value": 42},
		Timestamp:  time.Now().Format(time.RFC3339),
		DurationMs: 150,
	}

	jsonBytes, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal CheckResult: %v", err)
	}

	var decoded CheckResult
	if err := json.Unmarshal(jsonBytes, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal CheckResult: %v", err)
	}

	if decoded.Status != result.Status {
		t.Errorf("Expected status %s, got %s", result.Status, decoded.Status)
	}

	if decoded.Message != result.Message {
		t.Errorf("Expected message %s, got %s", result.Message, decoded.Message)
	}

	if decoded.Data["value"] != float64(42) {
		t.Errorf("Expected data value 42, got %v", decoded.Data["value"])
	}
}

func TestBuildEnvVars(t *testing.T) {
	exec := &DockerExecutor{}

	checkConfig := &config.CheckConfig{
		Name:    "test-check",
		Image:   "test-image",
		Enabled: true,
		Timeout: "30s",
		Env: map[string]string{
			"ENDPOINT":   "https://example.com",
			"CUSTOM_VAR": "custom-value",
		},
		Metadata: map[string]interface{}{
			"priority": "high",
			"owner":    "test-team",
		},
	}

	env, secretDir, err := exec.buildEnvVars(checkConfig)
	if err != nil {
		t.Fatalf("buildEnvVars failed: %v", err)
	}
	if secretDir != "" {
		t.Fatalf("did not expect secret directory, got %s", secretDir)
	}

	envMap := make(map[string]string)
	for _, e := range env {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			envMap[parts[0]] = parts[1]
		}
	}

	if envMap["FOGHORN_CHECK_NAME"] != "test-check" {
		t.Errorf("Expected FOGHORN_CHECK_NAME=test-check, got %s", envMap["FOGHORN_CHECK_NAME"])
	}

	if envMap["FOGHORN_ENDPOINT"] != "https://example.com" {
		t.Errorf("Expected FOGHORN_ENDPOINT=https://example.com, got %s", envMap["FOGHORN_ENDPOINT"])
	}

	if envMap["ENDPOINT"] != "https://example.com" {
		t.Errorf("Expected ENDPOINT=https://example.com, got %s", envMap["ENDPOINT"])
	}

	if envMap["FOGHORN_TIMEOUT"] != "30s" {
		t.Errorf("Expected FOGHORN_TIMEOUT=30s, got %s", envMap["FOGHORN_TIMEOUT"])
	}

	if envMap["CUSTOM_VAR"] != "custom-value" {
		t.Errorf("Expected CUSTOM_VAR=custom-value, got %s", envMap["CUSTOM_VAR"])
	}

	if envMap["FOGHORN_CHECK_CONFIG"] == "" {
		t.Error("FOGHORN_CHECK_CONFIG should be set when metadata is present")
	}
}

type testSecretResolver struct {
	values map[string]string
}

func (r *testSecretResolver) Resolve(ref string) (string, error) {
	value, ok := r.values[ref]
	if !ok {
		return "", fmt.Errorf("not found")
	}
	return value, nil
}

func TestBuildEnvVarsWithSecretReferences(t *testing.T) {
	exec := &DockerExecutor{
		secretResolver: &testSecretResolver{
			values: map[string]string{
				"secret://smtp/password": "smtp-secret",
			},
		},
	}

	checkConfig := &config.CheckConfig{
		Name:    "mail-check",
		Image:   "test-image",
		Enabled: true,
		Env: map[string]string{
			"SMTP_PASSWORD": "secret://smtp/password",
			"SMTP_HOST":     "smtp.example.com",
		},
	}

	env, secretDir, err := exec.buildEnvVars(checkConfig)
	if err != nil {
		t.Fatalf("buildEnvVars failed: %v", err)
	}
	if secretDir == "" {
		t.Fatal("expected secret directory to be created")
	}
	defer os.RemoveAll(secretDir)

	envMap := make(map[string]string)
	for _, e := range env {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			envMap[parts[0]] = parts[1]
		}
	}

	if envMap["SMTP_PASSWORD"] != "" {
		t.Error("expected SMTP_PASSWORD not to be passed directly")
	}
	if envMap["SMTP_PASSWORD_FILE"] != "/run/foghorn/secrets/SMTP_PASSWORD" {
		t.Errorf("unexpected SMTP_PASSWORD_FILE value: %q", envMap["SMTP_PASSWORD_FILE"])
	}
	if envMap["SMTP_HOST"] != "smtp.example.com" {
		t.Errorf("expected SMTP_HOST to be passed through, got %q", envMap["SMTP_HOST"])
	}

	secretBytes, err := os.ReadFile(secretDir + "/SMTP_PASSWORD")
	if err != nil {
		t.Fatalf("failed to read secret file: %v", err)
	}
	if string(secretBytes) != "smtp-secret" {
		t.Errorf("unexpected secret file content: %q", string(secretBytes))
	}
}

func TestBuildEnvVarsWithSecretReferenceMissingResolver(t *testing.T) {
	exec := &DockerExecutor{}
	checkConfig := &config.CheckConfig{
		Name:    "mail-check",
		Image:   "test-image",
		Enabled: true,
		Env: map[string]string{
			"SMTP_PASSWORD": "secret://smtp/password",
		},
	}

	_, _, err := exec.buildEnvVars(checkConfig)
	if err == nil {
		t.Fatal("expected error for missing secret resolver")
	}
}

func TestBuildEnvVarsWithEmptyResolvedSecret(t *testing.T) {
	exec := &DockerExecutor{
		secretResolver: &testSecretResolver{
			values: map[string]string{
				"secret://smtp/password": "",
			},
		},
	}
	checkConfig := &config.CheckConfig{
		Name:    "mail-check",
		Image:   "test-image",
		Enabled: true,
		Env: map[string]string{
			"SMTP_PASSWORD": "secret://smtp/password",
		},
	}

	_, _, err := exec.buildEnvVars(checkConfig)
	if err == nil {
		t.Fatal("expected error for empty resolved secret")
	}
	if !strings.Contains(err.Error(), "resolved to an empty value") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewDockerExecutor(t *testing.T) {
	exec, err := NewDockerExecutor()
	if err != nil {
		t.Fatalf("Failed to create Docker executor: %v", err)
	}

	if exec.cli == nil {
		t.Error("Expected cli to be initialized")
	}

	if exec.defaultTimeout != 30*time.Second {
		t.Errorf("Expected default timeout 30s, got %v", exec.defaultTimeout)
	}

	if exec.outputLocation != "stdout" {
		t.Errorf("Expected default output location 'stdout', got %s", exec.outputLocation)
	}

	if err := exec.Close(); err != nil {
		t.Errorf("Failed to close executor: %v", err)
	}
}

func TestDemultiplexLogs(t *testing.T) {
	jsonStr := `{"status":"pass","message":"test"}`

	header := []byte{0x01, 0x00, 0x00, 0x00}
	size := make([]byte, 4)
	size[0] = byte(len(jsonStr) >> 24)
	size[1] = byte(len(jsonStr) >> 16)
	size[2] = byte(len(jsonStr) >> 8)
	size[3] = byte(len(jsonStr))

	multiplexed := append(append(header, size...), jsonStr...)

	result := demultiplexLogs(multiplexed)

	if string(result) != jsonStr {
		t.Errorf("Expected %q, got %q", jsonStr, string(result))
	}

	multipleFrames := append(multiplexed, multiplexed...)
	result = demultiplexLogs(multipleFrames)

	expected := jsonStr + jsonStr
	if string(result) != expected {
		t.Errorf("Expected %q, got %q", expected, string(result))
	}
}

func TestReadResultWithMixedOutput(t *testing.T) {
	exec, err := NewDockerExecutor()
	if err != nil {
		t.Fatalf("Failed to create executor: %v", err)
	}
	defer exec.Close()

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "Pure JSON",
			input:   `{"status":"pass","message":"test","timestamp":"2026-01-20T00:00:00Z","duration_ms":100}`,
			wantErr: false,
		},
		{
			name:    "Mixed text and JSON",
			input:   `Testing HTTP to https://example.com\n{"status":"pass","message":"test","timestamp":"2026-01-20T00:00:00Z","duration_ms":100}`,
			wantErr: false,
		},
		{
			name:    "Multiple lines before JSON",
			input:   `Line 1\nLine 2\nLine 3\n{"status":"pass","message":"test","timestamp":"2026-01-20T00:00:00Z","duration_ms":100}`,
			wantErr: false,
		},
		{
			name:    "Invalid JSON",
			input:   `This is not JSON`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result CheckResult
			err := json.Unmarshal([]byte(tt.input), &result)

			if tt.wantErr && err == nil {
				t.Error("Expected error but got none")
			}

			if !tt.wantErr && err != nil {
				openBrace := strings.LastIndex(tt.input, "{")
				if openBrace != -1 {
					err = json.Unmarshal([]byte(tt.input[openBrace:]), &result)
				}
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestCheckResultValidation(t *testing.T) {
	tests := []struct {
		name    string
		result  CheckResult
		wantErr bool
	}{
		{
			name: "Valid pass status",
			result: CheckResult{
				Status:     "pass",
				Message:    "Test passed",
				Timestamp:  time.Now().Format(time.RFC3339),
				DurationMs: 100,
			},
			wantErr: false,
		},
		{
			name: "Valid fail status",
			result: CheckResult{
				Status:     "fail",
				Message:    "Test failed",
				Timestamp:  time.Now().Format(time.RFC3339),
				DurationMs: 100,
			},
			wantErr: false,
		},
		{
			name: "Valid warn status",
			result: CheckResult{
				Status:     "warn",
				Message:    "Test warning",
				Timestamp:  time.Now().Format(time.RFC3339),
				DurationMs: 100,
			},
			wantErr: false,
		},
		{
			name: "Valid unknown status",
			result: CheckResult{
				Status:     "unknown",
				Message:    "Test unknown",
				Timestamp:  time.Now().Format(time.RFC3339),
				DurationMs: 100,
			},
			wantErr: false,
		},
		{
			name: "Empty status",
			result: CheckResult{
				Status:     "",
				Message:    "Test",
				Timestamp:  time.Now().Format(time.RFC3339),
				DurationMs: 100,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBytes, err := json.Marshal(tt.result)
			if err != nil {
				t.Fatalf("Failed to marshal: %v", err)
			}

			var decoded CheckResult
			if err := json.Unmarshal(jsonBytes, &decoded); err != nil {
				if !tt.wantErr {
					t.Errorf("Unexpected error: %v", err)
				}
			} else if tt.wantErr {
				t.Error("Expected error but got none")
			}
		})
	}
}

func TestTruncateLogOutput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxChars int
		want     string
	}{
		{
			name:     "no truncation when shorter",
			input:    "hello",
			maxChars: 10,
			want:     "hello",
		},
		{
			name:     "no truncation when disabled",
			input:    "hello",
			maxChars: 0,
			want:     "hello",
		},
		{
			name:     "truncate with marker",
			input:    "abcdefghij",
			maxChars: 5,
			want:     "abcde\n... (truncated)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := truncateLogOutput(tt.input, tt.maxChars)
			if got != tt.want {
				t.Errorf("truncateLogOutput() = %q, want %q", got, tt.want)
			}
		})
	}
}
