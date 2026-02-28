package config

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/pfarrer/foghorn/containerimage"
	"gopkg.in/yaml.v3"
)

func Load(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	cfg := &Config{}

	for {
		var raw map[string]interface{}
		if err := decoder.Decode(&raw); err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("failed to parse YAML: %w", err)
		}
		if len(raw) == 0 {
			continue
		}

		if _, ok := raw["name"]; ok {
			var check CheckConfig
			if err := decodeInto(raw, &check); err != nil {
				return nil, fmt.Errorf("failed to parse check config: %w", err)
			}
			if check.Name != "" {
				cfg.Checks = append(cfg.Checks, check)
			}
			continue
		}

		var docCfg Config
		if err := decodeInto(raw, &docCfg); err != nil {
			return nil, fmt.Errorf("failed to parse config: %w", err)
		}
		mergeConfig(cfg, &docCfg)
	}

	if err := validate(cfg); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

func PrintSummary(cfg *Config) {
	fmt.Printf("Configuration loaded successfully\n")
	fmt.Printf("Version: %s\n", cfg.Version)
	fmt.Printf("Checks: %d\n", len(cfg.Checks))
	enabledCount := 0
	for _, check := range cfg.Checks {
		if check.Enabled {
			enabledCount++
		}
	}
	fmt.Printf("Enabled checks: %d\n", enabledCount)
	fmt.Printf("Disabled checks: %d\n", len(cfg.Checks)-enabledCount)
	if cfg.MaxConcurrentChecks > 0 {
		fmt.Printf("Max concurrent checks: %d\n", cfg.MaxConcurrentChecks)
	} else {
		fmt.Printf("Max concurrent checks: unlimited\n")
	}
}

func validate(cfg *Config) error {
	if cfg.MaxConcurrentChecks < 0 {
		return fmt.Errorf("max_concurrent_checks cannot be negative")
	}
	if cfg.StateLogFile != "" && cfg.StateLogPeriod == "" {
		return fmt.Errorf("state_log_period is required when state_log_file is set")
	}
	if cfg.StateLogPeriod != "" {
		period, err := time.ParseDuration(cfg.StateLogPeriod)
		if err != nil || period <= 0 {
			return fmt.Errorf("state_log_period must be a positive duration")
		}
	}
	if err := validateDebugOutputMode("config", cfg.DebugOutput); err != nil {
		return err
	}
	if cfg.DebugOutputMaxChars < 0 {
		return fmt.Errorf("debug_output_max_chars cannot be negative")
	}

	for i, check := range cfg.Checks {
		if check.Name == "" {
			return fmt.Errorf("check %d: name is required", i+1)
		}
		if check.Image == "" {
			return fmt.Errorf("check %s: image is required", check.Name)
		}
		if _, err := containerimage.ParseReference(check.Image); err != nil {
			return fmt.Errorf("check %s: invalid image tag: %w", check.Name, err)
		}
		if check.Schedule.Cron == "" && check.Schedule.Interval == "" {
			return fmt.Errorf("check %s: schedule (cron or interval) is required", check.Name)
		}
		if check.Schedule.Cron != "" && check.Schedule.Interval != "" {
			return fmt.Errorf("check %s: only one of cron or interval should be specified", check.Name)
		}
		if err := validateDebugOutputMode(fmt.Sprintf("check %s", check.Name), check.DebugOutput); err != nil {
			return err
		}
	}
	return nil
}

func validateDebugOutputMode(subject string, mode string) error {
	switch strings.TrimSpace(mode) {
	case "", "off", "on_failure", "always":
		return nil
	default:
		return fmt.Errorf("%s: debug_output must be one of off, on_failure, always", subject)
	}
}

func decodeInto(raw map[string]interface{}, dest interface{}) error {
	data, err := yaml.Marshal(raw)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, dest)
}

func mergeConfig(dst *Config, src *Config) {
	if src.Version != "" {
		dst.Version = src.Version
	}
	if src.MaxConcurrentChecks != 0 {
		dst.MaxConcurrentChecks = src.MaxConcurrentChecks
	}
	if src.StateLogFile != "" {
		dst.StateLogFile = src.StateLogFile
	}
	if src.StateLogPeriod != "" {
		dst.StateLogPeriod = src.StateLogPeriod
	}
	if src.SecretStoreFile != "" {
		dst.SecretStoreFile = src.SecretStoreFile
	}
	if src.DebugOutput != "" {
		dst.DebugOutput = src.DebugOutput
	}
	if src.DebugOutputMaxChars != 0 {
		dst.DebugOutputMaxChars = src.DebugOutputMaxChars
	}
	if len(src.Global) > 0 {
		if dst.Global == nil {
			dst.Global = make(map[string]interface{}, len(src.Global))
		}
		for k, v := range src.Global {
			dst.Global[k] = v
		}
	}
	if len(src.Checks) > 0 {
		dst.Checks = append(dst.Checks, src.Checks...)
	}
}
