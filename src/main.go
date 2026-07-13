package main

import (
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	
	"github.com/webappsgo/mail/src/config"
	"github.com/webappsgo/mail/src/database"
	"github.com/webappsgo/mail/src/logger"
	"github.com/webappsgo/mail/src/paths"
	"github.com/webappsgo/mail/src/pidfile"
	"github.com/webappsgo/mail/src/scheduler"
	"github.com/webappsgo/mail/src/server"
)

// Build info - set via -ldflags at build time (per AI.md PART 26)
var (
	Version      = "dev"
	CommitID     = "unknown"
	BuildDate    = "unknown"
	OfficialSite = "" // Empty = users must use --server flag
)

// Binary name (for display in help/error messages per AI.md PART 8)
var binaryName = filepath.Base(os.Args[0])

func main() {
	// PHASE 1: Parse flags and handle immediate-exit commands (per AI.md PART 8)
	// These commands print output and exit WITHOUT starting the server
	
	// Parse command line arguments
	if len(os.Args) < 2 {
		// No flags - start server normally (Phase 6)
		startServer()
		return
	}
	
	cmd := os.Args[1]
	
	// Handle immediate-exit flags
	switch cmd {
	case "--help", "-h":
		printHelp()
		os.Exit(0)
		
	case "--version", "-v":
		printVersion()
		os.Exit(0)
		
	case "--status":
		checkStatus()
		
	case "--shell":
		handleShell()
		
	case "--service":
		handleService()
		
	case "--maintenance":
		handleMaintenance()
		
	case "--update":
		handleUpdate()
		
	case "--daemon":
		// TODO: Implement daemon mode per AI.md PART 8
		fmt.Fprintf(os.Stderr, "%s: --daemon not yet implemented\n", binaryName)
		os.Exit(1)
		
	case "--debug":
		// Set debug mode and continue to server start
		// TODO: Set debug flag in config
		startServer()
		
	default:
		fmt.Fprintf(os.Stderr, "%s: unknown flag: %s\n", binaryName, cmd)
		fmt.Fprintf(os.Stderr, "Run '%s --help' for usage.\n", binaryName)
		os.Exit(1)
	}
}

// printHelp displays the main help message per AI.md PART 8
func printHelp() {
	fmt.Printf("%s %s - Email Infrastructure Management Panel\n\n", binaryName, Version)
	fmt.Printf("Usage:\n")
	fmt.Printf("  %s [flags]\n\n", binaryName)
	fmt.Printf("Information:\n")
	fmt.Printf("  -h, --help                        Show help\n")
	fmt.Printf("  -v, --version                     Show version\n")
	fmt.Printf("      --status                      Show server status and health\n\n")
	fmt.Printf("Shell Integration:\n")
	fmt.Printf("      --shell completions [SHELL]   Print shell completions\n")
	fmt.Printf("      --shell init [SHELL]          Print shell init command\n")
	fmt.Printf("      --shell --help                Show shell help\n\n")
	fmt.Printf("Server Configuration:\n")
	fmt.Printf("      --mode {production|development}  Application mode (default: production)\n")
	fmt.Printf("      --config DIR                  Config directory\n")
	fmt.Printf("      --data DIR                    Data directory\n")
	fmt.Printf("      --cache DIR                   Cache directory\n")
	fmt.Printf("      --log DIR                     Log directory\n")
	fmt.Printf("      --backup DIR                  Backup directory\n")
	fmt.Printf("      --pid FILE                    PID file path\n")
	fmt.Printf("      --address ADDR                Listen address (default: 0.0.0.0)\n")
	fmt.Printf("      --port PORT                   Listen port (default: random 64xxx)\n")
	fmt.Printf("      --daemon                      Run as daemon (detach from terminal)\n")
	fmt.Printf("      --debug                       Enable debug mode\n")
	fmt.Printf("      --color {always|never|auto}   Color output (default: auto)\n\n")
	fmt.Printf("Service Management:\n")
	fmt.Printf("      --service CMD                 Service management (--service --help for details)\n")
	fmt.Printf("      --maintenance CMD             Maintenance operations (--maintenance --help for details)\n")
	fmt.Printf("      --update [CMD]                Check/perform updates (--update --help for details)\n\n")
	fmt.Printf("Run '%s <command> --help' for detailed help on any command.\n", binaryName)
}

// printVersion displays version information per AI.md PART 8 and PART 26
func printVersion() {
	fmt.Printf("%s version %s\n", binaryName, Version)
	fmt.Println("Email Infrastructure Management Panel")
	fmt.Printf("Commit: %s\n", CommitID)
	fmt.Printf("Built: %s\n", BuildDate)
	fmt.Printf("Go: %s\n", runtime.Version())
	fmt.Printf("Platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	if OfficialSite != "" {
		fmt.Printf("Site: %s\n", OfficialSite)
	}
}

// checkStatus queries the running server and displays status
// Per AI.md PART 8: Exit 0 if healthy, 1 if unhealthy
func checkStatus() {
	// TODO: Implement status check per AI.md PART 13
	fmt.Fprintf(os.Stderr, "%s: --status not yet implemented\n", binaryName)
	os.Exit(1)
}

// handleShell manages shell integration commands
func handleShell() {
	if len(os.Args) < 3 {
		printShellHelp()
		os.Exit(1)
	}
	
	subcmd := os.Args[2]
	switch subcmd {
	case "--help":
		printShellHelp()
		os.Exit(0)
	case "completions":
		// TODO: Implement shell completions
		fmt.Fprintf(os.Stderr, "%s: shell completions not yet implemented\n", binaryName)
		os.Exit(1)
	case "init":
		// TODO: Implement shell init
		fmt.Fprintf(os.Stderr, "%s: shell init not yet implemented\n", binaryName)
		os.Exit(1)
	default:
		fmt.Fprintf(os.Stderr, "%s: unknown shell command: %s\n", binaryName, subcmd)
		fmt.Fprintf(os.Stderr, "Run '%s --shell --help' for usage.\n", binaryName)
		os.Exit(1)
	}
}

// printShellHelp displays shell integration help
func printShellHelp() {
	fmt.Printf("%s shell integration\n\n", binaryName)
	fmt.Printf("Usage:\n")
	fmt.Printf("  %s --shell <command> [options]\n\n", binaryName)
	fmt.Printf("Commands:\n")
	fmt.Printf("  completions [SHELL]   Generate shell completions (bash, zsh, fish, powershell)\n")
	fmt.Printf("  init [SHELL]          Print shell initialization command\n")
	fmt.Printf("  --help                Show this help\n")
}

// handleService manages service operations
func handleService() {
	if len(os.Args) < 3 {
		printServiceHelp()
		os.Exit(1)
	}
	
	subcmd := os.Args[2]
	switch subcmd {
	case "--help":
		printServiceHelp()
		os.Exit(0)
	default:
		// TODO: Implement service management per AI.md PART 20
		fmt.Fprintf(os.Stderr, "%s: service management not yet implemented\n", binaryName)
		os.Exit(1)
	}
}

// printServiceHelp displays service management help
func printServiceHelp() {
	fmt.Printf("%s service management\n\n", binaryName)
	fmt.Printf("Usage:\n")
	fmt.Printf("  %s --service <command>\n\n", binaryName)
	fmt.Printf("Commands:\n")
	fmt.Printf("  start        Start the service\n")
	fmt.Printf("  stop         Stop the service\n")
	fmt.Printf("  restart      Restart the service\n")
	fmt.Printf("  reload       Reload configuration\n")
	fmt.Printf("  --install    Install and enable service\n")
	fmt.Printf("  --disable    Disable service\n")
	fmt.Printf("  --uninstall  Uninstall service\n")
	fmt.Printf("  --help       Show this help\n")
}

// handleMaintenance manages maintenance operations
func handleMaintenance() {
	if len(os.Args) < 3 {
		printMaintenanceHelp()
		os.Exit(1)
	}
	
	subcmd := os.Args[2]
	switch subcmd {
	case "--help":
		printMaintenanceHelp()
		os.Exit(0)
	default:
		// TODO: Implement maintenance operations per AI.md PART 21
		fmt.Fprintf(os.Stderr, "%s: maintenance operations not yet implemented\n", binaryName)
		os.Exit(1)
	}
}

// printMaintenanceHelp displays maintenance operations help
func printMaintenanceHelp() {
	fmt.Printf("%s maintenance operations\n\n", binaryName)
	fmt.Printf("Usage:\n")
	fmt.Printf("  %s --maintenance <command> [options]\n\n", binaryName)
	fmt.Printf("Commands:\n")
	fmt.Printf("  backup [FILE]     Create backup\n")
	fmt.Printf("  restore FILE      Restore from backup\n")
	fmt.Printf("  update            Update application\n")
	fmt.Printf("  mode MODE         Set application mode\n")
	fmt.Printf("  setup             Run setup wizard\n")
	fmt.Printf("  --help            Show this help\n")
}

// handleUpdate manages update operations
func handleUpdate() {
	subcmd := "check" // Default to check
	if len(os.Args) >= 3 {
		subcmd = os.Args[2]
	}
	
	switch subcmd {
	case "--help":
		printUpdateHelp()
		os.Exit(0)
	default:
		// TODO: Implement update operations per AI.md PART 22
		fmt.Fprintf(os.Stderr, "%s: update operations not yet implemented\n", binaryName)
		os.Exit(1)
	}
}

// printUpdateHelp displays update operations help
func printUpdateHelp() {
	fmt.Printf("%s update operations\n\n", binaryName)
	fmt.Printf("Usage:\n")
	fmt.Printf("  %s --update [command]\n\n", binaryName)
	fmt.Printf("Commands:\n")
	fmt.Printf("  check              Check for updates (default)\n")
	fmt.Printf("  yes                Download and install update\n")
	fmt.Printf("  branch BRANCH      Switch to branch (stable, beta, daily)\n")
	fmt.Printf("  --help             Show this help\n")
}

// startServer initializes and starts the HTTP server
// Per AI.md PART 8: This is PHASE 6 - actual server startup
func startServer() {
	fmt.Printf("%s - Email Infrastructure Management Panel\n", binaryName)
	fmt.Printf("Version: %s (commit: %s)\n\n", Version, CommitID)
	
	// PHASE 6 per AI.md PART 8 Server Startup Sequence
	
	// Step 6: Determine run context
	isRoot := os.Getuid() == 0
	fmt.Printf("[INFO] Running as: ")
	if isRoot {
		fmt.Println("root (privileged)")
	} else {
		fmt.Println("user (non-privileged)")
	}
	
	// Step 7: Resolve all paths
	appPaths := paths.GetDefaultPaths()
	fmt.Printf("[INFO] Config dir: %s\n", appPaths.Config)
	fmt.Printf("[INFO] Data dir: %s\n", appPaths.Data)
	fmt.Printf("[INFO] Log dir: %s\n", appPaths.Log)
	fmt.Printf("[INFO] Database dir: %s\n", appPaths.DBDir)
	
	// Step 8/9: Setup directories
	fmt.Println("[INFO] Creating directories...")
	if err := appPaths.EnsureAllDirs(isRoot); err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Failed to create directories: %v\n", err)
		os.Exit(1)
	}
	
	// Step 10: Initialize logging per AI.md PART 11
	fmt.Println("[INFO] Initializing logging...")
	logConfig := logger.Config{
		Level:   "info",
		LogDir:  appPaths.Log,
		Format:  "text",
		MaxSize: 50 * 1024 * 1024, // 50MB
		MaxAge:  7,                 // 7 days
		Console: true,
	}
	
	if err := logger.Init(logConfig); err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Failed to initialize logging: %v\n", err)
		os.Exit(1)
	}
	defer logger.Close()
	
	logger.Info("Server starting...")
	logger.Infof("Version: %s (commit: %s)", Version, CommitID)
	
	// Step 11-12: Check and write PID file
	pidPath := getPIDPath(appPaths, isRoot)
	logger.Infof("PID file: %s", pidPath)
	
	pid, isRunning, err := pidfile.Check(pidPath)
	if err != nil {
		logger.Warnf("PID file check failed: %v", err)
	} else if isRunning {
		logger.Errorf("Server already running (PID %d)", pid)
		fmt.Fprintf(os.Stderr, "[ERROR] Server already running (PID %d)\n", pid)
		os.Exit(1)
	} else if pid > 0 {
		logger.Infof("Removing stale PID file (process %d not running)", pid)
		pidfile.Remove(pidPath)
	}
	
	if err := pidfile.Write(pidPath); err != nil {
		logger.Errorf("Failed to write PID file: %v", err)
		fmt.Fprintf(os.Stderr, "[ERROR] Failed to write PID file: %v\n", err)
		os.Exit(1)
	}
	defer pidfile.Remove(pidPath)
	
	// Step 13: Load configuration
	cfg, err := loadConfig(appPaths)
	if err != nil {
		logger.Errorf("Failed to load config: %v", err)
		fmt.Fprintf(os.Stderr, "[ERROR] Failed to load config: %v\n", err)
		os.Exit(1)
	}
	
	// Step 14: Reconfigure logging from config
	// TODO: Update log level from cfg.Server.Logs.Level when implemented
	logger.Info("Configuration loaded successfully")
	
	// Step 15: Initialize database
	logger.Info("Connecting to database...")
	db, err := initDatabase(cfg)
	if err != nil {
		logger.Errorf("Database initialization failed: %v", err)
		fmt.Fprintf(os.Stderr, "[ERROR] Database initialization failed: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()
	logger.Info("Database connected successfully")
	
	// Step 16: Start scheduler per AI.md PART 19
	logger.Info("Initializing scheduler...")
	sched, err := scheduler.New(db)
	if err != nil {
		logger.Errorf("Scheduler initialization failed: %v", err)
		fmt.Fprintf(os.Stderr, "[ERROR] Scheduler initialization failed: %v\n", err)
		os.Exit(1)
	}
	sched.Start()
	defer sched.Stop()
	logger.Info("Scheduler initialized successfully")
	
	// Step 17: Start Tor (if available)
	// TODO: Initialize Tor per AI.md PART 32
	
	// Step 18: Start HTTP server
	logger.Info("Starting HTTP server...")
	srv, err := startHTTPServer(cfg, db)
	if err != nil {
		logger.Errorf("Failed to start server: %v", err)
		fmt.Fprintf(os.Stderr, "[ERROR] Failed to start server: %v\n", err)
		os.Exit(1)
	}
	
	// Step 19: Register signal handlers (handled by srv.WaitForShutdown)
	
	// Step 20: Log startup complete per AI.md PART 17
	logger.Infof("Server started successfully on %s:%d", cfg.Server.Address, cfg.Server.Port)
	
	// Display startup banner per AI.md PART 17
	displayStartupBanner(cfg)
	
	// Step 21: Enter main loop
	srv.WaitForShutdown()
	
	logger.Info("Server stopped")
	fmt.Println("[INFO] Server stopped")
}

// displayStartupBanner displays the startup information per AI.md PART 17
func displayStartupBanner(cfg *config.Config) {
	fmt.Println()
	fmt.Println("╭───────────────────────────────────────────────────────────╮")
	fmt.Printf("│  🚀 MAIL · 📦 %s%-39s│\n", Version, "")
	fmt.Println("├───────────────────────────────────────────────────────────┤")
	fmt.Printf("│  📡 Listening on http://%s:%d%-23s│\n", cfg.Server.Address, cfg.Server.Port, "")
	fmt.Println("│  ✅ Server started successfully                           │")
	fmt.Println("╰───────────────────────────────────────────────────────────╯")
	fmt.Println()
}

// getPIDPath returns the PID file path based on privilege level
// Per AI.md PART 8 Step 11: root uses /var/run/, user uses {data_dir}/
func getPIDPath(appPaths *paths.Paths, isRoot bool) string {
	if isRoot {
		// Root: /var/run/webappsgo/mail.pid
		return "/var/run/webappsgo/mail.pid"
	}
	// User: {data_dir}/mail.pid
	return filepath.Join(appPaths.Data, "mail.pid")
}

// loadConfig loads the configuration file
// Per AI.md PART 8: Step 13
func loadConfig(appPaths *paths.Paths) (*config.Config, error) {
	configFile := filepath.Join(appPaths.Config, "server.yml")
	
	// Check if config file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// First run - create default config
		fmt.Println("[INFO] First run detected - creating default configuration")
		cfg := createDefaultConfig(appPaths)
		
		// Generate one-time setup token per AI.md PART 8 Step 13
		setupToken, err := generateSetupToken()
		if err != nil {
			return nil, fmt.Errorf("generating setup token: %w", err)
		}
		cfg.Server.SetupToken = setupToken
		
		// Save config to file per AI.md PART 5
		if err := saveConfig(cfg, configFile); err != nil {
			fmt.Fprintf(os.Stderr, "[WARN] Failed to save config file: %v\n", err)
		}
		
		// Display setup token in console per AI.md PART 17
		displaySetupToken(setupToken, cfg)
		
		return cfg, nil
	}
	
	// Load existing config from YAML file
	fmt.Printf("[INFO] Loading config from %s\n", configFile)
	cfg, err := config.Load(configFile)
	if err != nil {
		return nil, fmt.Errorf("loading config: %w", err)
	}
	
	// Update database path if it's relative or empty
	if cfg.Database.Driver == "sqlite" && cfg.Database.Name == "" {
		cfg.Database.Name = filepath.Join(appPaths.DBDir, "server.db")
	}
	
	return cfg, nil
}

// generateSetupToken generates a one-time setup token per AI.md PART 17
// Format: 32 hexadecimal characters (128-bit random)
func generateSetupToken() (string, error) {
	// Use crypto/rand for secure random bytes per AI.md PART 17
	bytes := make([]byte, 16) // 16 bytes = 32 hex chars
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("reading random bytes: %w", err)
	}
	return fmt.Sprintf("%x", bytes), nil
}

// displaySetupToken displays the setup token banner per AI.md PART 17
func displaySetupToken(token string, cfg *config.Config) {
	fmt.Println()
	fmt.Println("┌───────────────────────────────────────────────────────────┐")
	fmt.Println("│  🔑 SETUP REQUIRED                                        │")
	fmt.Println("├───────────────────────────────────────────────────────────┤")
	fmt.Printf("│  Setup Token: %-42s  │\n", token)
	fmt.Println("│                                                           │")
	fmt.Printf("│  Go to http://%s:%d/admin/server/setup         │\n", cfg.Server.Address, cfg.Server.Port)
	fmt.Println("│  and enter this token to complete setup.                 │")
	fmt.Println("│                                                           │")
	fmt.Println("│  This token will only be shown ONCE.                      │")
	fmt.Println("└───────────────────────────────────────────────────────────┘")
	fmt.Println()
}

// createDefaultConfig creates default configuration
func createDefaultConfig(appPaths *paths.Paths) *config.Config {
	return &config.Config{
		Server: config.ServerConfig{
			Address: "0.0.0.0",
			Port:    64500, // Default port per AI.md PART 5
		},
		Database: config.DatabaseConfig{
			Driver: "sqlite",
			Name:   filepath.Join(appPaths.DBDir, "server.db"),
		},
	}
}

// saveConfig saves configuration to YAML file
// Per AI.md PART 5: Comments MUST be above settings, never inline
func saveConfig(cfg *config.Config, path string) error {
	data := fmt.Sprintf(`# Mail Panel Configuration
# Generated automatically on first run

server:
  # Server listening address (0.0.0.0 = all interfaces)
  address: %s
  
  # Server listening port (64000-64999 range recommended)
  port: %d
  
  # Admin panel path (e.g., /admin)
  admin_path: admin
  
  # One-time setup token (generated on first run)
  setup_token: %s

database:
  # Database driver (sqlite, pgx, mysql)
  driver: %s
  
  # Database file path (SQLite) or connection name
  name: %s
`, cfg.Server.Address, cfg.Server.Port, cfg.Server.SetupToken, cfg.Database.Driver, cfg.Database.Name)
	
	return os.WriteFile(path, []byte(data), 0644)
}

// initDatabase initializes database connection and schema
// Per AI.md PART 8: Step 15
func initDatabase(cfg *config.Config) (*database.DB, error) {
	// Connect to database
	db, err := database.Connect(cfg)
	if err != nil {
		return nil, fmt.Errorf("database connection: %w", err)
	}
	
	// Create schema (idempotent per AI.md PART 10)
	if err := db.EnsureSchema(); err != nil {
		db.Close()
		return nil, fmt.Errorf("schema creation: %w", err)
	}
	
	return db, nil
}

// startHTTPServer creates and starts the HTTP server
// Per AI.md PART 8: Step 18
func startHTTPServer(cfg *config.Config, db *database.DB) (*server.Server, error) {
	srv := server.New(cfg, db)
	
	if err := srv.Start(); err != nil {
		return nil, err
	}
	
	return srv, nil
}
