package config

type Schedule struct {
	Cron     string `yaml:"cron,omitempty"`
	Interval string `yaml:"interval,omitempty"`
}

type EvaluationRule struct {
	Type      string                 `yaml:"type"`
	Condition string                 `yaml:"condition"`
	Threshold float64                `yaml:"threshold,omitempty"`
	Expected  interface{}            `yaml:"expected,omitempty"`
	Metadata  map[string]interface{} `yaml:"metadata,omitempty"`
}

type CheckConfig struct {
	Name                      string                 `yaml:"name"`
	Image                     string                 `yaml:"image"`
	Schedule                  Schedule               `yaml:"schedule"`
	Evaluation                []EvaluationRule       `yaml:"evaluation"`
	Description               string                 `yaml:"description,omitempty"`
	Tags                      []string               `yaml:"tags,omitempty"`
	Enabled                   bool                   `yaml:"enabled"`
	Env                       map[string]string      `yaml:"env,omitempty"`
	Timeout                   string                 `yaml:"timeout,omitempty"`
	CheckContainerDebugOutput string                 `yaml:"check_container_debug_output,omitempty"`
	Metadata                  map[string]interface{} `yaml:"metadata,omitempty"`
}

type Config struct {
	Checks                    []CheckConfig          `yaml:"checks"`
	Global                    map[string]interface{} `yaml:"global,omitempty"`
	Version                   string                 `yaml:"version,omitempty"`
	MaxConcurrentChecks       int                    `yaml:"max_concurrent_checks,omitempty"`
	StateLogFile              string                 `yaml:"state_log_file,omitempty"`
	StateLogPeriod            string                 `yaml:"state_log_period,omitempty"`
	SecretStoreFile           string                 `yaml:"secret_store_file,omitempty"`
	CheckContainerDebugOutput string                 `yaml:"check_container_debug_output,omitempty"`
	DebugOutputMaxChars       int                    `yaml:"debug_output_max_chars,omitempty"`
}
