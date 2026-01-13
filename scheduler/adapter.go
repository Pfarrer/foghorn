package scheduler

import (
	"github.com/anomalyco/foghorn/config"
)

type ConfigAdapter struct {
	cfg *config.CheckConfig
}

func NewConfigAdapter(cfg *config.CheckConfig) *ConfigAdapter {
	return &ConfigAdapter{cfg: cfg}
}

func (a *ConfigAdapter) GetName() string {
	return a.cfg.Name
}

func (a *ConfigAdapter) GetSchedule() string {
	if a.cfg.Schedule.Cron != "" {
		return a.cfg.Schedule.Cron
	}
	return ""
}

func (a *ConfigAdapter) IsEnabled() bool {
	return a.cfg.Enabled
}
