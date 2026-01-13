package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/anomalyco/foghorn/config"
	"github.com/anomalyco/foghorn/executor"
	"github.com/anomalyco/foghorn/scheduler"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: foghorn <config-file>")
		os.Exit(1)
	}

	configPath := os.Args[1]

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
