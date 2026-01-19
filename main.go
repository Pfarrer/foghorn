package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/docker/docker/client"
	"github.com/pfarrer/foghorn/config"
	"github.com/pfarrer/foghorn/executor"
	"github.com/pfarrer/foghorn/scheduler"
)

type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

func main() {
	var (
		help                    bool
		configPath              string
		logLevel                string
		verbose                 bool
		dryRun                  bool
		verifyImageAvailability bool
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
		fmt.Fprintf(os.Stderr, "  -h, --help\n")
		fmt.Fprintf(os.Stderr, "      Show help message\n")
	}

	flag.Parse()

	if help {
		flag.Usage()
		os.Exit(0)
	}

	if configPath == "" {
		fmt.Fprintf(os.Stderr, "Error: configuration file path is required\n\n")
		flag.Usage()
		os.Exit(1)
	}

	if err := validateLogLevel(LogLevel(logLevel)); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
		flag.Usage()
		os.Exit(1)
	}

	if verbose {
		log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
	} else {
		log.SetFlags(log.Ldate | log.Ltime)
	}

	if verbose {
		log.Printf("Loading configuration from: %s", configPath)
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	config.PrintSummary(cfg)

	if verifyImageAvailability {
		fmt.Println("Validating Docker images...")
		if err := verifyImageAvailabilityFn(cfg, verbose); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
	}

	if dryRun {
		fmt.Println("Configuration validation successful.")
		os.Exit(0)
	}

	dockerExecutor, err := executor.NewDockerExecutor()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating Docker executor: %v\n", err)
		os.Exit(1)
	}
	defer dockerExecutor.Close()

	maxConcurrent := cfg.MaxConcurrentChecks
	sched := scheduler.NewScheduler(dockerExecutor, time.UTC, maxConcurrent)

	for i := range cfg.Checks {
		check := &cfg.Checks[i]
		adapter := scheduler.NewConfigAdapter(check)
		if err := sched.AddCheck(adapter); err != nil {
			fmt.Fprintf(os.Stderr, "Error adding check %s: %v\n", check.Name, err)
		} else {
			scheduleType := check.Schedule.Cron
			if check.Schedule.Interval != "" {
				scheduleType = fmt.Sprintf("interval %s", check.Schedule.Interval)
			}
			fmt.Printf("Scheduled check: %s (%s)\n", check.Name, scheduleType)
		}
	}

	fmt.Println("\nScheduler started. Press Ctrl+C to stop.")
	sched.Start(1 * time.Second)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	sched.Stop()
	fmt.Println("\nScheduler stopped.")
}

func validateLogLevel(level LogLevel) error {
	switch level {
	case LogLevelDebug, LogLevelInfo, LogLevelWarn, LogLevelError:
		return nil
	default:
		return fmt.Errorf("invalid log level '%s', must be one of: debug, info, warn, error", level)
	}
}

func verifyImageAvailabilityFn(cfg *config.Config, verbose bool) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("failed to connect to Docker daemon: %w", err)
	}
	defer cli.Close()

	imageChecks := make(map[string][]string)
	enabledChecks := 0

	for _, check := range cfg.Checks {
		if check.Enabled {
			enabledChecks++
			if check.Image != "" {
				imageChecks[check.Image] = append(imageChecks[check.Image], check.Name)
			}
		}
	}

	if verbose {
		log.Printf("Checking %d unique images across %d enabled checks", len(imageChecks), enabledChecks)
	}

	missingImages := make(map[string][]string)

	for image, checkNames := range imageChecks {
		_, _, err := cli.ImageInspectWithRaw(context.Background(), image)
		if err != nil {
			if client.IsErrNotFound(err) {
				missingImages[image] = checkNames
				if verbose {
					log.Printf("Image not found locally: %s", image)
				}
			} else {
				return fmt.Errorf("error checking image %s: %w", image, err)
			}
		} else {
			if verbose {
				log.Printf("Image found locally: %s", image)
			}
		}
	}

	if len(missingImages) > 0 {
		var builder strings.Builder
		builder.WriteString("Error: The following Docker images are not available locally:\n\n")
		for image, checkNames := range missingImages {
			fmt.Fprintf(&builder, "- %s (required by: %s)\n", image, strings.Join(checkNames, ", "))
		}
		builder.WriteString("\nPlease pull the missing images:\n")
		for image := range missingImages {
			fmt.Fprintf(&builder, "  docker pull %s\n", image)
		}
		return fmt.Errorf("%s", builder.String())
	}

	fmt.Println("\nAll Docker images validated successfully:")
	for image, checkNames := range imageChecks {
		fmt.Printf("  - %s ✓", image)
		if verbose && len(checkNames) > 1 {
			fmt.Printf(" (used by: %s)", strings.Join(checkNames, ", "))
		}
		fmt.Println()
	}

	return nil
}
