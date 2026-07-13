package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Level represents log severity levels
type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

// String returns the string representation of the level
func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Logger provides structured logging per AI.md PART 11
type Logger struct {
	mu       sync.Mutex
	level    Level
	logDir   string
	
	// Log files per AI.md PART 11
	serverLog   *os.File
	errorLog    *os.File
	accessLog   *os.File
	auditLog    *os.File
	securityLog *os.File
	debugLog    *os.File
	
	// Format settings
	format      string // "text" or "json"
	consoleOut  io.Writer
	
	// Rotation settings (not yet implemented)
	maxSize     int64
	maxAge      int
}

// Config holds logger configuration
type Config struct {
	Level      string
	LogDir     string
	Format     string
	MaxSize    int64
	MaxAge     int
	Console    bool
}

// Global logger instance
var global *Logger
var once sync.Once

// Init initializes the global logger per AI.md PART 8 Step 10
func Init(cfg Config) error {
	var err error
	once.Do(func() {
		global, err = New(cfg)
	})
	return err
}

// New creates a new logger instance
func New(cfg Config) (*Logger, error) {
	level := parseLevel(cfg.Level)
	
	l := &Logger{
		level:      level,
		logDir:     cfg.LogDir,
		format:     cfg.Format,
		maxSize:    cfg.MaxSize,
		maxAge:     cfg.MaxAge,
		consoleOut: os.Stdout,
	}
	
	// Create log directory
	if err := os.MkdirAll(cfg.LogDir, 0755); err != nil {
		return nil, fmt.Errorf("creating log directory: %w", err)
	}
	
	// Open log files per AI.md PART 11
	var openErr error
	
	l.serverLog, openErr = l.openLogFile("server.log")
	if openErr != nil {
		return nil, fmt.Errorf("opening server.log: %w", openErr)
	}
	
	l.errorLog, openErr = l.openLogFile("error.log")
	if openErr != nil {
		return nil, fmt.Errorf("opening error.log: %w", openErr)
	}
	
	l.accessLog, openErr = l.openLogFile("access.log")
	if openErr != nil {
		return nil, fmt.Errorf("opening access.log: %w", openErr)
	}
	
	l.auditLog, openErr = l.openLogFile("audit.log")
	if openErr != nil {
		return nil, fmt.Errorf("opening audit.log: %w", openErr)
	}
	
	l.securityLog, openErr = l.openLogFile("security.log")
	if openErr != nil {
		return nil, fmt.Errorf("opening security.log: %w", openErr)
	}
	
	if level == DebugLevel {
		l.debugLog, openErr = l.openLogFile("debug.log")
		if openErr != nil {
			return nil, fmt.Errorf("opening debug.log: %w", openErr)
		}
	}
	
	return l, nil
}

// openLogFile opens a log file for appending
func (l *Logger) openLogFile(name string) (*os.File, error) {
	path := filepath.Join(l.logDir, name)
	return os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
}

// Close closes all log files
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	
	var errs []error
	
	if l.serverLog != nil {
		if err := l.serverLog.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if l.errorLog != nil {
		if err := l.errorLog.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if l.accessLog != nil {
		if err := l.accessLog.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if l.auditLog != nil {
		if err := l.auditLog.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if l.securityLog != nil {
		if err := l.securityLog.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if l.debugLog != nil {
		if err := l.debugLog.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	
	if len(errs) > 0 {
		return fmt.Errorf("closing log files: %v", errs)
	}
	
	return nil
}

// writeLog writes a log entry per AI.md PART 11
// Log files MUST use raw text only (no emojis, no ANSI codes)
func (l *Logger) writeLog(level Level, msg string, fields map[string]interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	
	// Skip if below log level
	if level < l.level {
		return
	}
	
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	
	var entry string
	if l.format == "json" {
		entry = l.formatJSON(timestamp, level, msg, fields)
	} else {
		entry = l.formatText(timestamp, level, msg, fields)
	}
	
	// Write to appropriate log file
	var logFile *os.File
	switch level {
	case ErrorLevel:
		logFile = l.errorLog
	case DebugLevel:
		logFile = l.debugLog
	default:
		logFile = l.serverLog
	}
	
	if logFile != nil {
		fmt.Fprintln(logFile, entry)
	}
	
	// Also write to server.log for all levels
	if logFile != l.serverLog && l.serverLog != nil {
		fmt.Fprintln(l.serverLog, entry)
	}
}

// formatText formats log entry as plain text per AI.md PART 11
func (l *Logger) formatText(timestamp string, level Level, msg string, fields map[string]interface{}) string {
	entry := fmt.Sprintf("%s [%s] %s", timestamp, level.String(), msg)
	
	if len(fields) > 0 {
		for k, v := range fields {
			entry += fmt.Sprintf(" %s=%v", k, v)
		}
	}
	
	return entry
}

// formatJSON formats log entry as JSON per AI.md PART 11
func (l *Logger) formatJSON(timestamp string, level Level, msg string, fields map[string]interface{}) string {
	// Basic JSON formatting (can be enhanced with proper JSON encoder)
	entry := fmt.Sprintf(`{"time":"%s","level":"%s","msg":"%s"`, timestamp, level.String(), msg)
	
	if len(fields) > 0 {
		for k, v := range fields {
			entry += fmt.Sprintf(`,"%s":"%v"`, k, v)
		}
	}
	
	entry += "}"
	return entry
}

// parseLevel converts string level to Level type
func parseLevel(s string) Level {
	switch s {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn", "warning":
		return WarnLevel
	case "error":
		return ErrorLevel
	default:
		return InfoLevel
	}
}

// Global logging functions

// Debug logs a debug message
func Debug(msg string) {
	if global != nil {
		global.writeLog(DebugLevel, msg, nil)
	}
}

// Debugf logs a formatted debug message
func Debugf(format string, args ...interface{}) {
	if global != nil {
		global.writeLog(DebugLevel, fmt.Sprintf(format, args...), nil)
	}
}

// Info logs an info message
func Info(msg string) {
	if global != nil {
		global.writeLog(InfoLevel, msg, nil)
	}
}

// Infof logs a formatted info message
func Infof(format string, args ...interface{}) {
	if global != nil {
		global.writeLog(InfoLevel, fmt.Sprintf(format, args...), nil)
	}
}

// Warn logs a warning message
func Warn(msg string) {
	if global != nil {
		global.writeLog(WarnLevel, msg, nil)
	}
}

// Warnf logs a formatted warning message
func Warnf(format string, args ...interface{}) {
	if global != nil {
		global.writeLog(WarnLevel, fmt.Sprintf(format, args...), nil)
	}
}

// Error logs an error message
func Error(msg string) {
	if global != nil {
		global.writeLog(ErrorLevel, msg, nil)
	}
}

// Errorf logs a formatted error message
func Errorf(format string, args ...interface{}) {
	if global != nil {
		global.writeLog(ErrorLevel, fmt.Sprintf(format, args...), nil)
	}
}

// Close closes the global logger
func Close() error {
	if global != nil {
		return global.Close()
	}
	return nil
}

// Audit logs an audit event per AI.md PART 11
// Audit log MUST be JSON only (machine-parseable)
func Audit(event string, actor string, fields map[string]interface{}) {
	if global == nil || global.auditLog == nil {
		return
	}
	
	global.mu.Lock()
	defer global.mu.Unlock()
	
	timestamp := time.Now().Format(time.RFC3339)
	
	// Audit log is ALWAYS JSON per AI.md PART 11
	entry := fmt.Sprintf(`{"time":"%s","event":"%s","actor":"%s"`, timestamp, event, actor)
	
	if fields != nil {
		for k, v := range fields {
			entry += fmt.Sprintf(`,"%s":"%v"`, k, v)
		}
	}
	
	entry += "}"
	fmt.Fprintln(global.auditLog, entry)
}

// Security logs a security event per AI.md PART 11
// Default format: fail2ban compatible
func Security(msg string, ip string) {
	if global == nil || global.securityLog == nil {
		return
	}
	
	global.mu.Lock()
	defer global.mu.Unlock()
	
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	
	// Fail2ban format per AI.md PART 11
	entry := fmt.Sprintf("%s [security] %s from %s", timestamp, msg, ip)
	fmt.Fprintln(global.securityLog, entry)
}

// Access logs an HTTP access event per AI.md PART 11
// Default format: Apache Combined Log Format
func Access(ip, method, path string, status, size int, userAgent, referer string) {
	if global == nil || global.accessLog == nil {
		return
	}
	
	global.mu.Lock()
	defer global.mu.Unlock()
	
	timestamp := time.Now().Format("02/Jan/2006:15:04:05 -0700")
	
	// Apache Combined Log Format per AI.md PART 11
	if referer == "" {
		referer = "-"
	}
	if userAgent == "" {
		userAgent = "-"
	}
	
	entry := fmt.Sprintf(`%s - - [%s] "%s %s HTTP/1.1" %d %d "%s" "%s"`,
		ip, timestamp, method, path, status, size, referer, userAgent)
	
	fmt.Fprintln(global.accessLog, entry)
}
