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
			config:  "checks:\n  - image: test\n    schedule:\n      cron: '* * * * *'\n    evaluation: []\n    enabled: true",
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
			config:  "checks:\n  - name: test\n    image: test\n    evaluation: []\n    enabled: true",
			wantErr: true,
			errMsg:  "schedule (cron or interval) is required",
		},
		{
			name:    "both cron and interval",
			config:  "checks:\n  - name: test\n    image: test\n    schedule:\n      cron: '* * * * *'\n      interval: '1m'\n    evaluation: []\n    enabled: true",
			wantErr: true,
			errMsg:  "only one of cron or interval should be specified",
		},
		{
			name:    "valid config with cron",
			config:  "checks:\n  - name: test\n    image: test\n    schedule:\n      cron: '* * * * *'\n    evaluation: []\n    enabled: true",
			wantErr: false,
		},
		{
			name:    "valid config with interval",
			config:  "checks:\n  - name: test\n    image: test\n    schedule:\n      interval: '1m'\n    evaluation: []\n    enabled: true",
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
	tests := []struct {
		name     string
		yamlPath string
		wantErr  bool
	}{
		{
			name:     "non-existent file",
			yamlPath: "../non-existent.yaml",
			wantErr:  true,
		},
		{
			name:     "invalid YAML",
			yamlPath: "../invalid.yaml",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Load(tt.yamlPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && err.Error() == "" {
				t.Error("Error should not be empty")
			}
		})
	}
}

func TestPrintSummary(t *testing.T) {
	cfg := &Config{
		Version: "1.0",
		Checks: []CheckConfig{
			{Name: "check1", Image: "img1", Schedule: Schedule{Cron: "* * * * *"}, Evaluation: []EvaluationRule{}, Enabled: true},
			{Name: "check2", Image: "img2", Schedule: Schedule{Interval: "1m"}, Evaluation: []EvaluationRule{}, Enabled: true},
			{Name: "check3", Image: "img3", Schedule: Schedule{Cron: "* * * * *"}, Evaluation: []EvaluationRule{}, Enabled: false},
		},
	}

	PrintSummary(cfg)
}
