package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/anomalyco/foghorn/config"
	"github.com/anomalyco/foghorn/executor"
	"github.com/anomalyco/foghorn/scheduler"
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
		help       bool
		configPath string
		logLevel   string
		verbose    bool
	)

	flag.BoolVar(&help, "h", false, "Show help message")
	flag.BoolVar(&help, "help", false, "Show help message")
	flag.StringVar(&configPath, "c", "", "Path to configuration file")
	flag.StringVar(&configPath, "config", "", "Path to configuration file")
	flag.StringVar(&logLevel, "l", "info", "Log level (debug, info, warn, error)")
	flag.StringVar(&logLevel, "log-level", "info", "Log level (debug, info, warn, error)")
	flag.BoolVar(&verbose, "v", false, "Enable verbose logging")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose logging")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Foghorn - Service Monitoring Tool\n\n")
		fmt.Fprintf(os.Stderr, "Usage: foghorn [OPTIONS]\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
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

	dockerExecutor, err := executor.NewDockerExecutor()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating Docker executor: %v\n", err)
		os.Exit(1)
	}
	defer dockerExecutor.Close()

	sched := scheduler.NewScheduler(dockerExecutor, time.UTC)

	for i := range cfg.Checks {
		check := &cfg.Checks[i]
		if check.Schedule.Cron != "" {
			adapter := scheduler.NewConfigAdapter(check)
			if err := sched.AddCheck(adapter); err != nil {
				fmt.Fprintf(os.Stderr, "Error adding check %s: %v\n", check.Name, err)
			} else {
				fmt.Printf("Scheduled check: %s (%s)\n", check.Name, check.Schedule.Cron)
			}
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
