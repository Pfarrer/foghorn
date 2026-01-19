package scheduler

import (
	"github.com/pfarrer/foghorn/config"
)

type ScheduleType string

const (
	ScheduleTypeCron     ScheduleType = "cron"
	ScheduleTypeInterval ScheduleType = "interval"
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
	return a.Config.Schedule.Interval
}

func (a *ConfigAdapter) GetScheduleType() ScheduleType {
	if a.Config.Schedule.Cron != "" {
		return ScheduleTypeCron
	}
	return ScheduleTypeInterval
}

func (a *ConfigAdapter) GetInterval() string {
	return a.Config.Schedule.Interval
}

func (a *ConfigAdapter) IsEnabled() bool {
	return a.Config.Enabled
}
