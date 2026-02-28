package config

import (
	"os"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestConfigLoad(t *testing.T) {
	data, err := os.ReadFile("../example.yaml")
	if err != nil {
		t.Fatalf("Failed to read example.yaml: %v", err)
	}

	if len(data) == 0 {
		t.Fatal("example.yaml is empty")
	}

	content := string(data)
	expectedFields := []string{"name:", "image:", "schedule:", "evaluation:", "enabled:"}

	for _, field := range expectedFields {
		if !strings.Contains(content, field) {
			t.Errorf("example.yaml missing expected field: %s", field)
		}
	}
}

func TestLoadValidConfig(t *testing.T) {
	cfg, err := Load("../example.yaml")
	if err != nil {
		t.Fatalf("Failed to load valid config: %v", err)
	}

	if len(cfg.Checks) == 0 {
		t.Error("Expected at least one check in config")
	}

	for _, check := range cfg.Checks {
		if check.Name == "" {
			t.Error("Check name should not be empty")
		}
		if check.Image == "" {
			t.Error("Check image should not be empty")
		}
		if check.Schedule.Cron == "" && check.Schedule.Interval == "" {
			t.Errorf("Check %s: should have cron or interval", check.Name)
		}
	}
}

func TestValidateRequiredFields(t *testing.T) {
	tests := []struct {
		name    string
		config  string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "missing name",
			config:  "checks:\n  - image: test/image:1.0.0\n    schedule:\n      cron: '* * * * *'\n    evaluation: []\n    enabled: true",
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name:    "missing image",
			config:  "checks:\n  - name: test\n    schedule:\n      cron: '* * * * *'\n    evaluation: []\n    enabled: true",
			wantErr: true,
			errMsg:  "image is required",
		},
		{
			name:    "missing schedule",
			config:  "checks:\n  - name: test\n    image: test/image:1.0.0\n    evaluation: []\n    enabled: true",
			wantErr: true,
			errMsg:  "schedule (cron or interval) is required",
		},
		{
			name:    "both cron and interval",
			config:  "checks:\n  - name: test\n    image: test/image:1.0.0\n    schedule:\n      cron: '* * * * *'\n      interval: '1m'\n    evaluation: []\n    enabled: true",
			wantErr: true,
			errMsg:  "only one of cron or interval should be specified",
		},
		{
			name:    "valid config with cron",
			config:  "checks:\n  - name: test\n    image: test/image:1.0.0\n    schedule:\n      cron: '* * * * *'\n    evaluation: []\n    enabled: true",
			wantErr: false,
		},
		{
			name:    "valid config with interval",
			config:  "checks:\n  - name: test\n    image: test/image:1.0.0\n    schedule:\n      interval: '1m'\n    evaluation: []\n    enabled: true",
			wantErr: false,
		},
		{
			name:    "invalid image tag",
			config:  "checks:\n  - name: test\n    image: test/image:latest\n    schedule:\n      interval: '1m'\n    evaluation: []\n    enabled: true",
			wantErr: true,
			errMsg:  "invalid image tag",
		},
		{
			name:    "missing image tag",
			config:  "checks:\n  - name: test\n    image: test/image\n    schedule:\n      interval: '1m'\n    evaluation: []\n    enabled: true",
			wantErr: true,
			errMsg:  "invalid image tag",
		},
		{
			name:    "invalid global debug_output mode",
			config:  "debug_output: noisy\nchecks:\n  - name: test\n    image: test/image:1.0.0\n    schedule:\n      interval: '1m'\n    evaluation: []\n    enabled: true",
			wantErr: true,
			errMsg:  "debug_output must be one of off, on_failure, always",
		},
		{
			name:    "invalid per-check debug_output mode",
			config:  "checks:\n  - name: test\n    image: test/image:1.0.0\n    debug_output: noisy\n    schedule:\n      interval: '1m'\n    evaluation: []\n    enabled: true",
			wantErr: true,
			errMsg:  "debug_output must be one of off, on_failure, always",
		},
		{
			name:    "negative debug output max chars",
			config:  "debug_output_max_chars: -1\nchecks:\n  - name: test\n    image: test/image:1.0.0\n    schedule:\n      interval: '1m'\n    evaluation: []\n    enabled: true",
			wantErr: true,
			errMsg:  "debug_output_max_chars cannot be negative",
		},
		{
			name:    "valid debug output config",
			config:  "debug_output: on_failure\ndebug_output_max_chars: 2048\nchecks:\n  - name: test\n    image: test/image:1.0.0\n    debug_output: always\n    schedule:\n      interval: '1m'\n    evaluation: []\n    enabled: true",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{}
			err := yaml.Unmarshal([]byte(tt.config), cfg)
			if err != nil {
				t.Fatalf("Failed to unmarshal test config: %v", err)
			}

			err = validate(cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("Error message should contain %q, got %q", tt.errMsg, err.Error())
			}
		})
	}
}

func TestHelpfulErrorMessages(t *testing.T) {
	t.Run("non-existent file", func(t *testing.T) {
		_, err := Load("../non-existent.yaml")
		if err == nil {
			t.Fatal("Load() should fail for non-existent file")
		}
		if err.Error() == "" {
			t.Error("Error should not be empty")
		}
	})

	t.Run("invalid YAML", func(t *testing.T) {
		tmpDir := t.TempDir()
		invalidPath := tmpDir + "/invalid.yaml"
		invalidContent := "this is not valid yaml:\n  - missing\n  colons\n    and\n  broken indentation\n"
		if writeErr := os.WriteFile(invalidPath, []byte(invalidContent), 0o644); writeErr != nil {
			t.Fatalf("failed to write invalid test YAML: %v", writeErr)
		}

		_, err := Load(invalidPath)
		if err == nil {
			t.Fatal("Load() should fail for invalid YAML")
		}
		if err.Error() == "" {
			t.Error("Error should not be empty")
		}
	})
}

func TestPrintSummary(t *testing.T) {
	cfg := &Config{
		Version: "1.0",
		Checks: []CheckConfig{
			{Name: "check1", Image: "img1:1.0.0", Schedule: Schedule{Cron: "* * * * *"}, Evaluation: []EvaluationRule{}, Enabled: true},
			{Name: "check2", Image: "img2:1.0.0", Schedule: Schedule{Interval: "1m"}, Evaluation: []EvaluationRule{}, Enabled: true},
			{Name: "check3", Image: "img3:1.0.0", Schedule: Schedule{Cron: "* * * * *"}, Evaluation: []EvaluationRule{}, Enabled: false},
		},
	}

	PrintSummary(cfg)
}
