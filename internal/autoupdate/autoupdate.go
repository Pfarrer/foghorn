package autoupdate

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/pfarrer/foghorn/config"
	"github.com/pfarrer/foghorn/containerimage"
	"github.com/pfarrer/foghorn/imageresolver"
	"github.com/pfarrer/foghorn/logger"
	"github.com/pfarrer/foghorn/scheduler"
)

type AutoUpdater struct {
	cli    *client.Client
	checks []config.CheckConfig
}

func NewAutoUpdater(cli *client.Client, checks []config.CheckConfig) *AutoUpdater {
	return &AutoUpdater{
		cli:    cli,
		checks: checks,
	}
}

func (u *AutoUpdater) StartScheduled(ctx context.Context, schedule config.Schedule) {
	interval := parseScheduleInterval(schedule)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	u.Run(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			u.Run(ctx)
		}
	}
}

func parseScheduleInterval(schedule config.Schedule) time.Duration {
	if schedule.Interval != "" {
		if d, err := time.ParseDuration(schedule.Interval); err == nil {
			return d
		}
	}
	if schedule.Cron != "" {
		cronExpr, err := scheduler.ParseCronExpression(schedule.Cron)
		if err == nil {
			next := cronExpr.Next(time.Now())
			if !next.IsZero() {
				return time.Until(next)
			}
		}
	}
	return 6 * time.Hour
}

func (u *AutoUpdater) Run(ctx context.Context) {
	logger.Info("Auto-update: starting check container image updates")

	for i := range u.checks {
		check := &u.checks[i]
		if !check.Enabled {
			continue
		}

		if err := u.updateCheck(ctx, check); err != nil {
			logger.Warn("Auto-update: check %s: %v", check.Name, err)
			continue
		}
	}

	logger.Info("Auto-update: completed")
}

func (u *AutoUpdater) updateCheck(ctx context.Context, check *config.CheckConfig) error {
	resolved, err := imageresolver.Resolve(ctx, u.cli, check.Image)
	if err != nil {
		return fmt.Errorf("failed to resolve image: %w", err)
	}

	currentID := ""
	currentInspect, _, err := u.cli.ImageInspectWithRaw(ctx, resolved)
	if err != nil {
		if !client.IsErrNotFound(err) {
			return fmt.Errorf("failed to inspect image: %w", err)
		}
	} else {
		currentID = currentInspect.ID
	}

	if checkNeedsResolution(check.Image) {
		resolved, err = imageresolver.Resolve(ctx, u.cli, check.Image)
		if err != nil {
			return fmt.Errorf("failed to resolve latest image: %w", err)
		}
	}

	pullCtx, cancel := context.WithTimeout(ctx, 5*60*1e9)
	defer cancel()

	reader, err := u.cli.ImagePull(pullCtx, resolved, image.PullOptions{})
	if err != nil {
		return fmt.Errorf("failed to pull image %s: %w", resolved, err)
	}
	defer reader.Close()

	if _, err := io.Copy(io.Discard, reader); err != nil {
		return fmt.Errorf("failed to complete pull for %s: %w", resolved, err)
	}

	newInspect, _, inspectErr := u.cli.ImageInspectWithRaw(ctx, resolved)
	if inspectErr != nil {
		logger.Info("Auto-update: check %s: pulled image %s", check.Name, resolved)
		return nil
	}

	if currentID != "" && currentID == newInspect.ID {
		logger.Debug("Auto-update: check %s: image %s is already up to date", check.Name, resolved)
		return nil
	}

	logger.Info("Auto-update: check %s: updated image to %s", check.Name, resolved)
	return nil
}

func checkNeedsResolution(imageRef string) bool {
	ref, err := containerimage.ParseReference(imageRef)
	if err != nil {
		return false
	}
	return ref.Selector.Kind != containerimage.SelectorFull
}
