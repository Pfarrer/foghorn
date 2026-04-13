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
			name:    "invalid global check_container_debug_output mode",
			config:  "check_container_debug_output: noisy\nchecks:\n  - name: test\n    image: test/image:1.0.0\n    schedule:\n      interval: '1m'\n    evaluation: []\n    enabled: true",
			wantErr: true,
			errMsg:  "check_container_debug_output must be one of off, on_failure, always",
		},
		{
			name:    "invalid per-check check_container_debug_output mode",
			config:  "checks:\n  - name: test\n    image: test/image:1.0.0\n    check_container_debug_output: noisy\n    schedule:\n      interval: '1m'\n    evaluation: []\n    enabled: true",
			wantErr: true,
			errMsg:  "check_container_debug_output must be one of off, on_failure, always",
		},
		{
			name:    "negative debug output max chars",
			config:  "debug_output_max_chars: -1\nchecks:\n  - name: test\n    image: test/image:1.0.0\n    schedule:\n      interval: '1m'\n    evaluation: []\n    enabled: true",
			wantErr: true,
			errMsg:  "debug_output_max_chars cannot be negative",
		},
		{
			name:    "valid debug output config",
			config:  "check_container_debug_output: on_failure\ndebug_output_max_chars: 2048\nchecks:\n  - name: test\n    image: test/image:1.0.0\n    check_container_debug_output: always\n    schedule:\n      interval: '1m'\n    evaluation: []\n    enabled: true",
			wantErr: false,
		},
		{
			name:    "auto_update enabled without schedule",
			config:  "auto_update_containers: true\nchecks:\n  - name: test\n    image: test/image:1.0.0\n    schedule:\n      interval: '1m'\n    evaluation: []\n    enabled: true",
			wantErr: true,
			errMsg:  "auto_update_schedule is required when auto_update_containers is true",
		},
		{
			name:    "auto_update with both cron and interval",
			config:  "auto_update_containers: true\nauto_update_schedule:\n  cron: '* * * * *'\n  interval: '1h'\nchecks:\n  - name: test\n    image: test/image:1.0.0\n    schedule:\n      interval: '1m'\n    evaluation: []\n    enabled: true",
			wantErr: true,
			errMsg:  "auto_update_schedule: only one of cron or interval should be specified",
		},
		{
			name:    "auto_update enabled with interval",
			config:  "auto_update_containers: true\nauto_update_schedule:\n  interval: '6h'\nchecks:\n  - name: test\n    image: test/image:1.0.0\n    schedule:\n      interval: '1m'\n    evaluation: []\n    enabled: true",
			wantErr: false,
		},
		{
			name:    "auto_update enabled with cron",
			config:  "auto_update_containers: true\nauto_update_schedule:\n  cron: '0 */6 * * *'\nchecks:\n  - name: test\n    image: test/image:1.0.0\n    schedule:\n      interval: '1m'\n    evaluation: []\n    enabled: true",
			wantErr: false,
		},
		{
			name:    "auto_update disabled without schedule",
			config:  "auto_update_containers: false\nchecks:\n  - name: test\n    image: test/image:1.0.0\n    schedule:\n      interval: '1m'\n    evaluation: []\n    enabled: true",
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

func TestCheckConfigAllFields(t *testing.T) {
	yamlContent := `checks:
  - name: full-check
    image: test/image:1.0.0
    description: "A full check config"
    enabled: true
    tags:
      - production
      - critical
    schedule:
      interval: "5m"
    evaluation:
      - type: json
        condition: "status == 'pass'"
        expected: true
    env:
      HOST: "example.com"
      PORT: "443"
    timeout: "30s"
    metadata:
      team: infra
      priority: high
`
	cfg := &Config{}
	err := yaml.Unmarshal([]byte(yamlContent), cfg)
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if len(cfg.Checks) != 1 {
		t.Fatalf("expected 1 check, got %d", len(cfg.Checks))
	}
	c := cfg.Checks[0]

	if c.Name != "full-check" {
		t.Errorf("name = %q, want full-check", c.Name)
	}
	if c.Image != "test/image:1.0.0" {
		t.Errorf("image = %q", c.Image)
	}
	if c.Description != "A full check config" {
		t.Errorf("description = %q", c.Description)
	}
	if !c.Enabled {
		t.Error("enabled should be true")
	}
	if len(c.Tags) != 2 || c.Tags[0] != "production" || c.Tags[1] != "critical" {
		t.Errorf("tags = %v", c.Tags)
	}
	if c.Schedule.Interval != "5m" {
		t.Errorf("interval = %q", c.Schedule.Interval)
	}
	if len(c.Evaluation) != 1 {
		t.Fatalf("evaluation len = %d", len(c.Evaluation))
	}
	if c.Evaluation[0].Type != "json" {
		t.Errorf("eval type = %q", c.Evaluation[0].Type)
	}
	if c.Env["HOST"] != "example.com" || c.Env["PORT"] != "443" {
		t.Errorf("env = %v", c.Env)
	}
	if c.Timeout != "30s" {
		t.Errorf("timeout = %q", c.Timeout)
	}
	if c.Metadata["team"] != "infra" {
		t.Errorf("metadata = %v", c.Metadata)
	}
}

func TestMinimalCheckConfig(t *testing.T) {
	yamlContent := `checks:
  - name: minimal
    image: test/image:1.0.0
    schedule:
      cron: "*/5 * * * *"
    enabled: true
`
	cfg := &Config{}
	err := yaml.Unmarshal([]byte(yamlContent), cfg)
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if err := validate(cfg); err != nil {
		t.Fatalf("minimal config should be valid: %v", err)
	}
}

func TestScheduleCronAndInterval(t *testing.T) {
	cronYAML := `checks:
  - name: cron-check
    image: test/image:1.0.0
    schedule:
      cron: "*/5 * * * *"
    enabled: true
`
	cfg := &Config{}
	if err := yaml.Unmarshal([]byte(cronYAML), cfg); err != nil {
		t.Fatalf("cron unmarshal: %v", err)
	}
	if cfg.Checks[0].Schedule.Cron != "*/5 * * * *" {
		t.Errorf("cron = %q", cfg.Checks[0].Schedule.Cron)
	}

	intervalYAML := `checks:
  - name: interval-check
    image: test/image:1.0.0
    schedule:
      interval: "5m"
    enabled: true
`
	cfg2 := &Config{}
	if err := yaml.Unmarshal([]byte(intervalYAML), cfg2); err != nil {
		t.Fatalf("interval unmarshal: %v", err)
	}
	if cfg2.Checks[0].Schedule.Interval != "5m" {
		t.Errorf("interval = %q", cfg2.Checks[0].Schedule.Interval)
	}
}

func TestEvaluationRuleParsing(t *testing.T) {
	yamlContent := `checks:
  - name: eval-check
    image: test/image:1.0.0
    schedule:
      interval: "1m"
    enabled: true
    evaluation:
      - type: threshold
        condition: "response_time <"
        threshold: 500.0
        expected: true
        metadata:
          unit: ms
`
	cfg := &Config{}
	if err := yaml.Unmarshal([]byte(yamlContent), cfg); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	r := cfg.Checks[0].Evaluation[0]
	if r.Type != "threshold" {
		t.Errorf("type = %q", r.Type)
	}
	if r.Condition != "response_time <" {
		t.Errorf("condition = %q", r.Condition)
	}
	if r.Threshold != 500.0 {
		t.Errorf("threshold = %f", r.Threshold)
	}
	if r.Expected != true {
		t.Errorf("expected = %v", r.Expected)
	}
	if r.Metadata["unit"] != "ms" {
		t.Errorf("metadata = %v", r.Metadata)
	}
}

func TestGlobalConfigFields(t *testing.T) {
	yamlContent := `version: "1.0"
max_concurrent_checks: 5
state_log_file: "/var/lib/foghorn/state.log"
state_log_period: "24h"
secret_store_file: "/etc/foghorn/secrets.enc"
check_container_debug_output: "on_failure"
debug_output_max_chars: 2048
global:
  timezone: UTC
checks:
  - name: test
    image: test/image:1.0.0
    schedule:
      interval: "1m"
    enabled: true
`
	cfg := &Config{}
	if err := yaml.Unmarshal([]byte(yamlContent), cfg); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if cfg.Version != "1.0" {
		t.Errorf("version = %q", cfg.Version)
	}
	if cfg.MaxConcurrentChecks != 5 {
		t.Errorf("max_concurrent_checks = %d", cfg.MaxConcurrentChecks)
	}
	if cfg.StateLogFile != "/var/lib/foghorn/state.log" {
		t.Errorf("state_log_file = %q", cfg.StateLogFile)
	}
	if cfg.StateLogPeriod != "24h" {
		t.Errorf("state_log_period = %q", cfg.StateLogPeriod)
	}
	if cfg.SecretStoreFile != "/etc/foghorn/secrets.enc" {
		t.Errorf("secret_store_file = %q", cfg.SecretStoreFile)
	}
	if cfg.CheckContainerDebugOutput != "on_failure" {
		t.Errorf("check_container_debug_output = %q", cfg.CheckContainerDebugOutput)
	}
	if cfg.DebugOutputMaxChars != 2048 {
		t.Errorf("debug_output_max_chars = %d", cfg.DebugOutputMaxChars)
	}
	if cfg.Global["timezone"] != "UTC" {
		t.Errorf("global = %v", cfg.Global)
	}
}

func TestEnvMapParsing(t *testing.T) {
	yamlContent := `checks:
  - name: env-check
    image: test/image:1.0.0
    schedule:
      interval: "1m"
    enabled: true
    env:
      HOST: "example.com"
      PORT: "443"
`
	cfg := &Config{}
	if err := yaml.Unmarshal([]byte(yamlContent), cfg); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	env := cfg.Checks[0].Env
	if env["HOST"] != "example.com" {
		t.Errorf("HOST = %q", env["HOST"])
	}
	if env["PORT"] != "443" {
		t.Errorf("PORT = %q", env["PORT"])
	}
}

func TestMixedConfigFormat(t *testing.T) {
	yamlContent := `version: "1.0"
checks:
  - name: list-check
    image: test/image:1.0.0
    schedule:
      interval: "1m"
    enabled: true
---
name: doc-check
image: test/image:2.0.0
schedule:
  cron: "* * * * *"
enabled: true
`
	tmpDir := t.TempDir()
	invalidPath := tmpDir + "/mixed.yaml"
	if writeErr := os.WriteFile(invalidPath, []byte(yamlContent), 0o644); writeErr != nil {
		t.Fatalf("failed to write mixed test YAML: %v", writeErr)
	}

	cfg, err := Load(invalidPath)
	if err != nil {
		t.Fatalf("Load() failed for mixed YAML: %v", err)
	}

	if len(cfg.Checks) != 2 {
		t.Fatalf("Expected 2 checks, got %d", len(cfg.Checks))
	}

	foundList := false
	foundDoc := false
	for _, c := range cfg.Checks {
		if c.Name == "list-check" {
			foundList = true
		}
		if c.Name == "doc-check" {
			foundDoc = true
		}
	}

	if !foundList || !foundDoc {
		t.Errorf("Failed to parse both list and doc checks from mixed format")
	}
}
