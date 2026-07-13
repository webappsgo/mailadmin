package pidfile

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

// Check checks if a PID file exists and if the process is running
// Per AI.md PART 8 Step 11: Check PID file
// Returns (pid, isRunning, error)
func Check(pidPath string) (int, bool, error) {
	// Check if PID file exists
	data, err := os.ReadFile(pidPath)
	if err != nil {
		if os.IsNotExist(err) {
			// No PID file - not running
			return 0, false, nil
		}
		return 0, false, fmt.Errorf("reading PID file: %w", err)
	}
	
	// Parse PID from file
	pidStr := strings.TrimSpace(string(data))
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		// Invalid PID file - consider stale
		return 0, false, fmt.Errorf("invalid PID in file: %s", pidStr)
	}
	
	// Check if process is running
	// Per AI.md PART 8: Send signal 0 to check if process exists
	process, err := os.FindProcess(pid)
	if err != nil {
		// Process not found - stale PID
		return pid, false, nil
	}
	
	// Try to signal process (signal 0 = check existence)
	err = process.Signal(syscall.Signal(0))
	if err != nil {
		// Process not running - stale PID
		return pid, false, nil
	}
	
	// Process exists and is running
	return pid, true, nil
}

// Write writes the current process ID to the PID file
// Per AI.md PART 8 Step 12: Write PID file
func Write(pidPath string) error {
	// Ensure directory exists
	dir := filepath.Dir(pidPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating PID directory: %w", err)
	}
	
	// Write PID to file
	pid := os.Getpid()
	data := fmt.Sprintf("%d\n", pid)
	
	// Write with 0644 permissions
	if err := os.WriteFile(pidPath, []byte(data), 0644); err != nil {
		return fmt.Errorf("writing PID file: %w", err)
	}
	
	return nil
}

// Remove removes the PID file
// Per AI.md PART 8: Remove PID file during graceful shutdown
func Remove(pidPath string) error {
	if err := os.Remove(pidPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("removing PID file: %w", err)
	}
	return nil
}
