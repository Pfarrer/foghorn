package main

import (
	"fmt"
	"os"
	"time"

	"github.com/anomalyco/foghorn/config"
	"github.com/anomalyco/foghorn/scheduler"
)

type SimpleExecutor struct{}

func (e *SimpleExecutor) Execute(check scheduler.CheckConfig) error {
	fmt.Printf("[%s] Executing check: %s\n", time.Now().Format(time.RFC3339), check.GetName())
	return nil
}

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

	executor := &SimpleExecutor{}
	sched := scheduler.NewScheduler(executor, time.UTC)

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

	select {}
}
