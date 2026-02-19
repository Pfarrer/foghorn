package daemon

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/docker/docker/client"
	"github.com/pfarrer/foghorn/config"
	"github.com/pfarrer/foghorn/executor"
	"github.com/pfarrer/foghorn/imageresolver"
	"github.com/pfarrer/foghorn/logger"
	"github.com/pfarrer/foghorn/scheduler"
	"github.com/pfarrer/foghorn/state"
	"github.com/pfarrer/foghorn/tui"
)

func Run() {
	var (
		help                    bool
		configPath              string
		logLevel                string
		verbose                 bool
		dryRun                  bool
		verifyImageAvailability bool
		tuiMode                 bool
		stateLogFile            string
	)

	flag.BoolVar(&help, "h", false, "Show help message")
	flag.BoolVar(&help, "help", false, "Show help message")
	flag.StringVar(&configPath, "c", "", "Path to configuration file")
	flag.StringVar(&configPath, "config", "", "Path to configuration file")
	flag.StringVar(&logLevel, "l", "info", "Log level (debug, info, warn, error)")
	flag.StringVar(&logLevel, "log-level", "info", "Log level (debug, info, warn, error)")
	flag.BoolVar(&verbose, "v", false, "Enable verbose logging")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose logging")
	flag.BoolVar(&dryRun, "d", false, "Validate configuration only")
	flag.BoolVar(&dryRun, "dry-run", false, "Validate configuration only")
	flag.BoolVar(&verifyImageAvailability, "i", false, "Verify all Docker images in config are available locally")
	flag.BoolVar(&verifyImageAvailability, "verify-image-availability", false, "Verify all Docker images in config are available locally")
	flag.BoolVar(&tuiMode, "tui", false, "Enable TUI dashboard mode")
	flag.BoolVar(&tuiMode, "t", false, "Enable TUI dashboard mode")
	flag.StringVar(&stateLogFile, "s", "", "Path to state log file")
	flag.StringVar(&stateLogFile, "state-log-file", "", "Path to state log file")
	flag.StringVar(&stateLogFile, "state_log_file", "", "Path to state log file")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Foghorn - Service Monitoring Tool\n\n")
		fmt.Fprintf(os.Stderr, "Usage: foghorn [OPTIONS]\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "  -c, --config <path>\n")
		fmt.Fprintf(os.Stderr, "      Path to configuration file\n")
		fmt.Fprintf(os.Stderr, "  -l, --log-level <level>\n")
		fmt.Fprintf(os.Stderr, "      Log level (debug, info, warn, error) (default: info)\n")
		fmt.Fprintf(os.Stderr, "  -v, --verbose\n")
		fmt.Fprintf(os.Stderr, "      Enable verbose logging\n")
		fmt.Fprintf(os.Stderr, "  -d, --dry-run\n")
		fmt.Fprintf(os.Stderr, "      Validate configuration only\n")
		fmt.Fprintf(os.Stderr, "  -i, --verify-image-availability\n")
		fmt.Fprintf(os.Stderr, "      Verify all Docker images in config are available locally\n")
		fmt.Fprintf(os.Stderr, "  -t, --tui\n")
		fmt.Fprintf(os.Stderr, "      Enable TUI dashboard mode\n")
		fmt.Fprintf(os.Stderr, "  -s, --state-log-file <path>\n")
		fmt.Fprintf(os.Stderr, "      Path to state log file\n")
		fmt.Fprintf(os.Stderr, "  -h, --help\n")
		fmt.Fprintf(os.Stderr, "      Show help message\n")
	}

	flag.Parse()

	if help {
		flag.Usage()
		os.Exit(0)
	}

	lvl, err := logger.ParseLevel(logLevel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
		flag.Usage()
		os.Exit(1)
	}
	logger.SetGlobal(logger.New(lvl, verbose))

	if configPath == "" {
		fmt.Fprintf(os.Stderr, "Error: configuration file path is required\n\n")
		flag.Usage()
		os.Exit(1)
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	logger.Info("Loaded configuration with %d checks", len(cfg.Checks))

	stateLogPath := stateLogFile
	if stateLogPath == "" {
		stateLogPath = cfg.StateLogFile
	}

	var stateRecords map[string]scheduler.CheckState
	var stateLog *state.StateLog
	if stateLogPath != "" {
		if cfg.StateLogPeriod == "" {
			fmt.Fprintf(os.Stderr, "Error: state_log_period is required when state_log_file is set\n")
			os.Exit(1)
		}
		retention, err := time.ParseDuration(cfg.StateLogPeriod)
		if err != nil || retention <= 0 {
			fmt.Fprintf(os.Stderr, "Error: state_log_period must be a positive duration\n")
			os.Exit(1)
		}

		stateLog, err = state.Open(stateLogPath, retention)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening state log: %v\n", err)
			os.Exit(1)
		}
		defer stateLog.Close()

		records, err := stateLog.Load()
		if err != nil {
			logger.Warn("Failed to load state log: %v", err)
		} else {
			history := buildHistory(records, 10)
			latest := state.LatestByCheck(records)
			stateRecords = make(map[string]scheduler.CheckState, len(latest))
			for name, record := range latest {
				stateRecords[name] = scheduler.CheckState{
					LastStatus:   record.Status,
					LastDuration: time.Duration(record.DurationMs) * time.Millisecond,
					LastRun:      record.CompletedAt,
					History:      history[name],
				}
			}
			for name, entries := range history {
				if _, ok := stateRecords[name]; ok {
					continue
				}
				stateRecords[name] = scheduler.CheckState{
					History: entries,
				}
			}
		}
	}

	if verifyImageAvailability {
		if err := verifyImageAvailabilityFn(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
	}

	if dryRun {
		logger.Info("Configuration validation successful.")
		os.Exit(0)
	}

	dockerExecutor, err := executor.NewDockerExecutor()
	if err != nil {
		logger.Error("Error creating Docker executor: %v", err)
		fmt.Fprintf(os.Stderr, "Error creating Docker executor: %v\n", err)
		os.Exit(1)
	}
	defer dockerExecutor.Close()

	maxConcurrent := cfg.MaxConcurrentChecks
	if maxConcurrent > 0 {
		logger.Info("Maximum concurrent checks: %d", maxConcurrent)
	}
	sched := scheduler.NewScheduler(dockerExecutor, time.UTC, maxConcurrent)
	if stateLog != nil {
		sched.SetResultLogger(stateLog)
	}

	for i := range cfg.Checks {
		check := &cfg.Checks[i]
		adapter := scheduler.NewConfigAdapter(check)
		if err := sched.AddCheck(adapter); err != nil {
			logger.Error("Error adding check %s: %v", check.Name, err)
			fmt.Fprintf(os.Stderr, "Error adding check %s: %v\n", check.Name, err)
		}
	}

	if len(stateRecords) > 0 {
		sched.ApplyState(stateRecords)
	}

	sched.Start(1 * time.Second)

	if tuiMode {
		logger.SetOutput(io.Discard)
		model := tui.NewModel(sched, logLevel)
		p := tea.NewProgram(model)
		if _, err := p.Run(); err != nil {
			logger.SetOutput(os.Stdout)
			logger.Error("Error running TUI: %v", err)
			fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		}
		logger.SetOutput(os.Stdout)
	} else {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
	}

	sched.Stop()
}

func verifyImageAvailabilityFn(cfg *config.Config) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logger.Error("Failed to connect to Docker daemon: %v", err)
		return fmt.Errorf("failed to connect to Docker daemon: %w", err)
	}
	defer cli.Close()

	logger.Info("Validating Docker images...")

	imageChecks := make(map[string][]string)
	unresolvedChecks := make(map[string][]string)
	unresolvedErrors := make(map[string]error)
	enabledChecks := 0

	for _, check := range cfg.Checks {
		if check.Enabled {
			enabledChecks++
			if check.Image != "" {
				resolved, err := imageresolver.Resolve(context.Background(), cli, check.Image)
				if err != nil {
					unresolvedChecks[check.Image] = append(unresolvedChecks[check.Image], check.Name)
					unresolvedErrors[check.Image] = err
					continue
				}
				imageChecks[resolved] = append(imageChecks[resolved], check.Name)
			}
		}
	}

	logger.Debug("Checking %d images for %d enabled checks", len(imageChecks), enabledChecks)

	missingImages := make(map[string][]string)

	for image, checkNames := range imageChecks {
		logger.Debug("Checking image: %s", image)
		_, _, err := cli.ImageInspectWithRaw(context.Background(), image)
		if err != nil {
			if client.IsErrNotFound(err) {
				logger.Warn("Image %s not available locally (required by: %s)", image, strings.Join(checkNames, ", "))
				missingImages[image] = checkNames
			} else {
				logger.Error("Error checking image %s: %v", image, err)
				return fmt.Errorf("error checking image %s: %w", image, err)
			}
		} else {
			logger.Debug("Image %s is available", image)
		}
	}

	if len(unresolvedChecks) > 0 || len(missingImages) > 0 {
		var builder strings.Builder
		if len(unresolvedChecks) > 0 {
			builder.WriteString("Error: The following image selectors could not be resolved from registry tags:\n\n")
			for image, checkNames := range unresolvedChecks {
				if resolveErr, ok := unresolvedErrors[image]; ok {
					fmt.Fprintf(&builder, "- %s (required by: %s, reason: %s)\n", image, strings.Join(checkNames, ", "), resolveErr.Error())
				} else {
					fmt.Fprintf(&builder, "- %s (required by: %s)\n", image, strings.Join(checkNames, ", "))
				}
			}
			builder.WriteString("\n")
		}
		if len(missingImages) > 0 {
			builder.WriteString("Error: The following Docker images are not available locally:\n\n")
			for image, checkNames := range missingImages {
				fmt.Fprintf(&builder, "- %s (required by: %s)\n", image, strings.Join(checkNames, ", "))
			}
			builder.WriteString("\nPlease pull the missing images:\n")
			for image := range missingImages {
				fmt.Fprintf(&builder, "  docker pull %s\n", image)
			}
		}
		return fmt.Errorf("%s", builder.String())
	}

	logger.Info("All Docker images validated successfully:")
	for image := range imageChecks {
		logger.Info("  - %s ✓", image)
	}

	return nil
}

func buildHistory(records []state.Record, maxEntries int) map[string][]scheduler.CheckHistoryEntry {
	if len(records) == 0 || maxEntries <= 0 {
		return nil
	}

	history := make(map[string][]scheduler.CheckHistoryEntry)
	for _, record := range records {
		if record.CheckName == "" {
			continue
		}
		history[record.CheckName] = append(history[record.CheckName], scheduler.CheckHistoryEntry{
			Status:      record.Status,
			CompletedAt: record.CompletedAt,
		})
	}

	for name, entries := range history {
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].CompletedAt.Before(entries[j].CompletedAt)
		})
		if len(entries) > maxEntries {
			entries = entries[len(entries)-maxEntries:]
		}
		history[name] = entries
	}

	return history
}
