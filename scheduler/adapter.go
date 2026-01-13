package scheduler

import (
	"github.com/anomalyco/foghorn/config"
)

type ConfigAdapter struct {
	Config *config.CheckConfig
}

func NewConfigAdapter(cfg *config.CheckConfig) *ConfigAdapter {
	return &ConfigAdapter{Config: cfg}
}

func (a *ConfigAdapter) GetName() string {
	return a.Config.Name
}

func (a *ConfigAdapter) GetSchedule() string {
	if a.Config.Schedule.Cron != "" {
		return a.Config.Schedule.Cron
	}
	return ""
}

func (a *ConfigAdapter) IsEnabled() bool {
	return a.Config.Enabled
}
