package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pfarrer/foghorn/internal/statusapi"
	"github.com/pfarrer/foghorn/tui"
)

func main() {
	var (
		help      bool
		statusURL string
		logLevel  string
	)

	flag.BoolVar(&help, "h", false, "Show help message")
	flag.BoolVar(&help, "help", false, "Show help message")
	flag.StringVar(&statusURL, "status-url", statusapi.DefaultBaseURL, "Daemon status API base URL")
	flag.StringVar(&statusURL, "u", statusapi.DefaultBaseURL, "Daemon status API base URL")
	flag.StringVar(&logLevel, "l", "info", "Log level label for display")
	flag.StringVar(&logLevel, "log-level", "info", "Log level label for display")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Foghorn TUI Client\n\n")
		fmt.Fprintf(os.Stderr, "Usage: foghorn-tui [OPTIONS]\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "  -u, --status-url <url>\n")
		fmt.Fprintf(os.Stderr, "      Daemon status API base URL (default: %s)\n", statusapi.DefaultBaseURL)
		fmt.Fprintf(os.Stderr, "  -l, --log-level <level>\n")
		fmt.Fprintf(os.Stderr, "      Log level label for display (default: info)\n")
		fmt.Fprintf(os.Stderr, "  -h, --help\n")
		fmt.Fprintf(os.Stderr, "      Show help message\n")
	}
	flag.Parse()

	if help {
		flag.Usage()
		return
	}

	model, err := tui.NewRemoteModel(statusURL, logLevel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to daemon status API: %v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		os.Exit(1)
	}
}
