package config

import (
	"fmt"
	"os"

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
		var check CheckConfig
		if err := decoder.Decode(&check); err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, fmt.Errorf("failed to parse YAML: %w", err)
		}

		if check.Name != "" {
			cfg.Checks = append(cfg.Checks, check)
		}
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
	}
	return nil
}
