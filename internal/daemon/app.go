package daemon

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/docker/docker/client"
	"github.com/pfarrer/foghorn/config"
	"github.com/pfarrer/foghorn/executor"
	"github.com/pfarrer/foghorn/imageresolver"
	"github.com/pfarrer/foghorn/internal/autoupdate"
	"github.com/pfarrer/foghorn/internal/statusapi"
	"github.com/pfarrer/foghorn/logger"
	"github.com/pfarrer/foghorn/scheduler"
	"github.com/pfarrer/foghorn/secretstore"
	"github.com/pfarrer/foghorn/state"
)

func Run() {
	if len(os.Args) > 1 && os.Args[1] == "secret" {
		exitCode := runSecretCLI(os.Args[2:])
		os.Exit(exitCode)
	}

	var (
		help                    bool
		configPath              string
		logLevel                string
		verbose                 bool
		dryRun                  bool
		oneShot                 bool
		verifyImageAvailability bool
		statusListen            string
		stateLogFile            string
		secretStoreFile         string
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
	flag.BoolVar(&oneShot, "one-shot", false, "Run all enabled checks once and exit")
	flag.BoolVar(&verifyImageAvailability, "i", false, "Verify all Docker images in config are available locally")
	flag.BoolVar(&verifyImageAvailability, "verify-image-availability", false, "Verify all Docker images in config are available locally")
	flag.StringVar(&statusListen, "status-listen", statusapi.DefaultListenAddr, "Status API listen address")
	flag.StringVar(&stateLogFile, "s", "", "Path to state log file")
	flag.StringVar(&stateLogFile, "state-log-file", "", "Path to state log file")
	flag.StringVar(&stateLogFile, "state_log_file", "", "Path to state log file")
	flag.StringVar(&secretStoreFile, "secret-store-file", "", "Path to encrypted secret store file")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Foghorn Daemon - Service Monitoring Tool\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  foghorn-daemon [OPTIONS]\n")
		fmt.Fprintf(os.Stderr, "  foghorn-daemon secret [OPTIONS] <COMMAND> [ARGS]\n\n")
		fmt.Fprintf(os.Stderr, "Commands:\n")
		fmt.Fprintf(os.Stderr, "  secret\n")
		fmt.Fprintf(os.Stderr, "      Manage encrypted secrets (run `foghorn-daemon secret --help`)\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "  -c, --config <path>\n")
		fmt.Fprintf(os.Stderr, "      Path to configuration file\n")
		fmt.Fprintf(os.Stderr, "  -l, --log-level <level>\n")
		fmt.Fprintf(os.Stderr, "      Log level (debug, info, warn, error) (default: info)\n")
		fmt.Fprintf(os.Stderr, "  -v, --verbose\n")
		fmt.Fprintf(os.Stderr, "      Enable verbose logging\n")
		fmt.Fprintf(os.Stderr, "  -d, --dry-run\n")
		fmt.Fprintf(os.Stderr, "      Validate configuration only\n")
		fmt.Fprintf(os.Stderr, "  --one-shot\n")
		fmt.Fprintf(os.Stderr, "      Run all enabled checks once and exit\n")
		fmt.Fprintf(os.Stderr, "  -i, --verify-image-availability\n")
		fmt.Fprintf(os.Stderr, "      Verify all Docker images in config are available locally\n")
		fmt.Fprintf(os.Stderr, "  --status-listen <addr>\n")
		fmt.Fprintf(os.Stderr, "      Status API listen address (default: %s)\n", statusapi.DefaultListenAddr)
		fmt.Fprintf(os.Stderr, "  -s, --state-log-file <path>\n")
		fmt.Fprintf(os.Stderr, "      Path to state log file\n")
		fmt.Fprintf(os.Stderr, "  --secret-store-file <path>\n")
		fmt.Fprintf(os.Stderr, "      Path to encrypted secret store file\n")
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
	dockerExecutor.SetDebugOutput(cfg.CheckContainerDebugOutput, cfg.DebugOutputMaxChars)

	if configUsesSecrets(cfg) {
		storePath := resolveSecretStorePath(secretStoreFile, cfg.SecretStoreFile)
		store, err := loadSecretStore(storePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading secret store: %v\n", err)
			os.Exit(1)
		}
		dockerExecutor.SetSecretResolver(store)
		logger.Info("Secret store enabled: %s", storePath)
	}

	if oneShot {
		maxConcurrent := cfg.MaxConcurrentChecks
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

		exitCode := runOneShot(sched, dockerExecutor, maxConcurrent)
		os.Exit(exitCode)
	}

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

	autoUpdateCancel := func() {}
	if cfg.AutoUpdateContainers {
		dockerCli, dockerErr := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if dockerErr != nil {
			logger.Warn("Auto-update: failed to create Docker client: %v", dockerErr)
		} else {
			updater := autoupdate.NewAutoUpdater(dockerCli, cfg.Checks)
			auCtx, auCancel := context.WithCancel(context.Background())
			autoUpdateCancel = func() {
				auCancel()
				dockerCli.Close()
			}
			go updater.StartScheduled(auCtx, cfg.AutoUpdateSchedule)
			logger.Info("Auto-update: enabled with schedule (cron=%q interval=%q)", cfg.AutoUpdateSchedule.Cron, cfg.AutoUpdateSchedule.Interval)
		}
	}

	statusSrv := statusapi.StartServer(statusListen, sched.Snapshot, stateLog)
	statusErr := make(chan error, 1)
	go func() {
		logger.Info("Status API listening on http://%s%s", statusListen, statusapi.StatusPath)
		if err := statusSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			statusErr <- err
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-sigChan:
	case err := <-statusErr:
		logger.Error("Status API server error: %v", err)
		fmt.Fprintf(os.Stderr, "Status API server error: %v\n", err)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := statusSrv.Shutdown(shutdownCtx); err != nil {
		logger.Warn("Status API shutdown error: %v", err)
	}
	autoUpdateCancel()
	sched.Stop()
}

func runOneShot(sched *scheduler.Scheduler, exec scheduler.CheckExecutor, maxConcurrent int) int {
	allChecks := sched.GetAllChecks()

	var enabledChecks []string
	for name, check := range allChecks {
		if check.Config.IsEnabled() {
			enabledChecks = append(enabledChecks, name)
		}
	}

	if len(enabledChecks) == 0 {
		logger.Info("One-shot mode: no enabled checks to run")
		return 0
	}

	logger.Info("One-shot mode: running %d checks", len(enabledChecks))

	var wg sync.WaitGroup

	limit := maxConcurrent
	if limit <= 0 {
		limit = len(enabledChecks)
	}
	sem := make(chan struct{}, limit)

	for _, name := range enabledChecks {
		check := allChecks[name]
		wg.Add(1)
		sem <- struct{}{}
		go func(name string, check *scheduler.ScheduledCheck) {
			defer wg.Done()
			defer func() { <-sem }()
			if err := exec.Execute(check.Config); err != nil {
				logger.Error("One-shot: check %s failed: %v", name, err)
			}
		}(name, check)
	}

	wg.Wait()

	passCount, warnCount, failCount, errorCount, unknownCount := 0, 0, 0, 0, 0
	for _, name := range enabledChecks {
		if check, ok := allChecks[name]; ok {
			switch check.LastStatus {
			case "pass":
				passCount++
			case "warn":
				warnCount++
			case "fail":
				failCount++
			case "error":
				errorCount++
			default:
				unknownCount++
			}
		}
	}

	logger.Info("One-shot complete: %d pass, %d warn, %d fail, %d error, %d unknown",
		passCount, warnCount, failCount, errorCount, unknownCount)

	if failCount > 0 || errorCount > 0 || unknownCount > 0 {
		return 1
	}
	return 0
}

func runSecretCLI(args []string) int {
	fs := flag.NewFlagSet("secret", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	var (
		storePathArg  string
		configPathArg string
		valueArg      string
		help          bool
	)
	fs.StringVar(&storePathArg, "s", "", "Path to encrypted secret store file")
	fs.StringVar(&storePathArg, "secret-store-file", "", "Path to encrypted secret store file")
	fs.StringVar(&configPathArg, "c", "", "Path to configuration file")
	fs.StringVar(&configPathArg, "config", "", "Path to configuration file")
	fs.StringVar(&valueArg, "value", "", "Secret value (avoid this flag in shared environments)")
	fs.BoolVar(&help, "h", false, "Show help message")
	fs.BoolVar(&help, "help", false, "Show help message")
	fs.Usage = printSecretUsage

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		printSecretUsage()
		return 1
	}
	if help {
		printSecretUsage()
		return 0
	}

	parts := fs.Args()
	if len(parts) == 0 {
		fmt.Fprintf(os.Stderr, "Error: secret command is required\n")
		printSecretUsage()
		return 1
	}

	storePath := resolveSecretStorePath(storePathArg, configSecretStorePath(configPathArg))
	store, err := loadSecretStore(storePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	cmd := parts[0]
	switch cmd {
	case "list":
		keys, err := store.ListKeys()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing secrets: %v\n", err)
			return 1
		}
		for _, key := range keys {
			fmt.Println(key)
		}
		return 0
	case "set", "rotate":
		if len(parts) < 2 {
			fmt.Fprintf(os.Stderr, "Error: secret key is required\n")
			printSecretUsage()
			return 1
		}
		key := parts[1]
		value := valueArg
		if value == "" {
			v, err := readSecretValueFromStdin()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading secret from stdin: %v\n", err)
				return 1
			}
			value = v
		}
		if value == "" {
			fmt.Fprintf(os.Stderr, "Error: secret value cannot be empty\n")
			return 1
		}
		if err := store.Set(key, value); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing secret: %v\n", err)
			return 1
		}
		fmt.Printf("stored secret key: %s\n", key)
		return 0
	case "delete":
		if len(parts) < 2 {
			fmt.Fprintf(os.Stderr, "Error: secret key is required\n")
			printSecretUsage()
			return 1
		}
		key := parts[1]
		deleted, err := store.Delete(key)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error deleting secret: %v\n", err)
			return 1
		}
		if !deleted {
			fmt.Printf("secret key not found: %s\n", key)
			return 0
		}
		fmt.Printf("deleted secret key: %s\n", key)
		return 0
	default:
		fmt.Fprintf(os.Stderr, "Error: unknown secret command %q\n", cmd)
		printSecretUsage()
		return 1
	}
}

func printSecretUsage() {
	fmt.Fprintf(os.Stderr, "Foghorn Secret Management\n\n")
	fmt.Fprintf(os.Stderr, "Usage:\n")
	fmt.Fprintf(os.Stderr, "  foghorn-daemon secret [OPTIONS] list\n")
	fmt.Fprintf(os.Stderr, "  foghorn-daemon secret [OPTIONS] [--value <val>] set <key>\n")
	fmt.Fprintf(os.Stderr, "  foghorn-daemon secret [OPTIONS] [--value <val>] rotate <key>\n")
	fmt.Fprintf(os.Stderr, "  foghorn-daemon secret [OPTIONS] delete <key>\n\n")
	fmt.Fprintf(os.Stderr, "Commands:\n")
	fmt.Fprintf(os.Stderr, "  list            List secret keys\n")
	fmt.Fprintf(os.Stderr, "  set <key>       Create or update a secret value\n")
	fmt.Fprintf(os.Stderr, "  rotate <key>    Rotate a secret value (same behavior as set)\n")
	fmt.Fprintf(os.Stderr, "  delete <key>    Delete a secret key\n\n")
	fmt.Fprintf(os.Stderr, "Options:\n")
	fmt.Fprintf(os.Stderr, "  -c, --config <path>\n")
	fmt.Fprintf(os.Stderr, "      Path to configuration file (used to resolve secret_store_file)\n")
	fmt.Fprintf(os.Stderr, "  -s, --secret-store-file <path>\n")
	fmt.Fprintf(os.Stderr, "      Path to encrypted secret store file\n")
	fmt.Fprintf(os.Stderr, "  --value <val>\n")
	fmt.Fprintf(os.Stderr, "      Secret value (avoid this flag in shared environments)\n")
	fmt.Fprintf(os.Stderr, "  -h, --help\n")
	fmt.Fprintf(os.Stderr, "      Show help message\n")
	fmt.Fprintf(os.Stderr, "Notes:\n")
	fmt.Fprintf(os.Stderr, "  - Set/rotate reads value from stdin when --value is omitted.\n")
	fmt.Fprintf(os.Stderr, "  - Secret store path precedence: --secret-store-file, config, FOGHORN_SECRET_STORE_FILE, default path.\n")
	fmt.Fprintf(os.Stderr, "  - Requires FOGHORN_SECRET_MASTER_KEY in environment.\n")
}

func readSecretValueFromStdin() (string, error) {
	info, err := os.Stdin.Stat()
	if err != nil {
		return "", err
	}
	if (info.Mode() & os.ModeCharDevice) != 0 {
		return "", errors.New("no stdin provided; pipe a secret value or use --value")
	}

	reader := bufio.NewReader(os.Stdin)
	data, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

func configSecretStorePath(configPath string) string {
	if configPath == "" {
		return ""
	}
	cfg, err := config.Load(configPath)
	if err != nil {
		return ""
	}
	return cfg.SecretStoreFile
}

func configUsesSecrets(cfg *config.Config) bool {
	for _, check := range cfg.Checks {
		for _, value := range check.Env {
			if _, ok := secretstore.ParseRef(value); ok {
				return true
			}
		}
	}
	return false
}

func resolveSecretStorePath(cliPath string, cfgPath string) string {
	if cliPath != "" {
		return cliPath
	}
	if cfgPath != "" {
		return cfgPath
	}
	if envPath := strings.TrimSpace(os.Getenv("FOGHORN_SECRET_STORE_FILE")); envPath != "" {
		return envPath
	}
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return ".foghorn-secrets.enc"
	}
	return filepath.Join(home, ".config", "foghorn", "secrets.enc")
}

func loadSecretStore(path string) (*secretstore.Store, error) {
	masterKey, err := secretstore.MasterKeyFromEnv()
	if err != nil {
		return nil, err
	}
	return secretstore.New(path, masterKey)
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
