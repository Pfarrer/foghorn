package executor

import (
	"encoding/json"
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
	exec, err := NewDockerExecutor()
	if err != nil {
		t.Fatalf("Failed to create executor: %v", err)
	}

	checkConfig := &config.CheckConfig{
		Name:    "test-check",
		Image:   "test-image",
		Enabled: true,
		Timeout: "30s",
		Env: map[string]string{
			"ENDPOINT":   "https://example.com",
			"CUSTOM_VAR": "custom-value",
			"SECRET_KEY": "secret-value",
		},
		Metadata: map[string]interface{}{
			"priority": "high",
			"owner":    "test-team",
		},
	}

	env := exec.buildEnvVars(checkConfig)

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

	if envMap["FOGHORN_TIMEOUT"] != "30s" {
		t.Errorf("Expected FOGHORN_TIMEOUT=30s, got %s", envMap["FOGHORN_TIMEOUT"])
	}

	if envMap["CUSTOM_VAR"] != "custom-value" {
		t.Errorf("Expected CUSTOM_VAR=custom-value, got %s", envMap["CUSTOM_VAR"])
	}

	if envMap["SECRET_KEY"] == "secret-value" {
		t.Error("SECRET_KEY should not be passed directly")
	}

	if envMap["FOGHORN_SECRETS"] == "" {
		t.Error("FOGHORN_SECRETS should be set when SECRET_KEY is present")
	}

	if envMap["FOGHORN_CHECK_CONFIG"] == "" {
		t.Error("FOGHORN_CHECK_CONFIG should be set when metadata is present")
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
