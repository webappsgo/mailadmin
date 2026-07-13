package paths

import (
	"os"
	"path/filepath"
	"runtime"
)

// Paths holds all application paths
type Paths struct {
	Config  string
	Data    string
	Cache   string
	Log     string
	Backup  string
	PID     string
	SSL     string
	Security string
	DBDir   string
}

// GetDefaultPaths returns OS and privilege-specific default paths
// Per AI.md PART 4: OS-SPECIFIC PATHS
func GetDefaultPaths() *Paths {
	isRoot := os.Getuid() == 0
	isDocker := isInDocker()
	
	if isDocker {
		return getDockerPaths()
	}
	
	switch runtime.GOOS {
	case "linux":
		return getLinuxPaths(isRoot)
	case "darwin":
		return getDarwinPaths(isRoot)
	case "freebsd", "openbsd", "netbsd":
		return getBSDPaths(isRoot)
	case "windows":
		return getWindowsPaths(isRoot)
	default:
		// Fallback to Linux paths
		return getLinuxPaths(isRoot)
	}
}

// isInDocker detects if running inside Docker container
// Per AI.md PART 8: Detect container vs local
func isInDocker() bool {
	// Check for /.dockerenv file
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}
	
	// Check for container environment variable
	if os.Getenv("container") != "" {
		return true
	}
	
	return false
}

// getDockerPaths returns Docker container paths
// Per AI.md PART 4: Docker/Container section
func getDockerPaths() *Paths {
	return &Paths{
		Config:   "/config/mail",
		Data:     "/data/mail",
		Cache:    "/data/mail/cache",
		Log:      "/data/log/mail",
		Backup:   "/data/backups/mail",
		PID:      "/var/run/apimgr/mail.pid",
		SSL:      "/config/mail/ssl",
		Security: "/config/mail/security",
		DBDir:    "/data/db/sqlite",
	}
}

// getLinuxPaths returns Linux-specific paths
// Per AI.md PART 4: Linux section
func getLinuxPaths(isRoot bool) *Paths {
	if isRoot {
		return &Paths{
			Config:   "/etc/apimgr/mail",
			Data:     "/var/lib/apimgr/mail",
			Cache:    "/var/cache/apimgr/mail",
			Log:      "/var/log/apimgr/mail",
			Backup:   "/mnt/Backups/apimgr/mail",
			PID:      "/var/run/apimgr/mail.pid",
			SSL:      "/etc/apimgr/mail/ssl",
			Security: "/etc/apimgr/mail/security",
			DBDir:    "/var/lib/apimgr/mail/db",
		}
	}
	
	home := os.Getenv("HOME")
	return &Paths{
		Config:   filepath.Join(home, ".config/apimgr/mail"),
		Data:     filepath.Join(home, ".local/share/apimgr/mail"),
		Cache:    filepath.Join(home, ".cache/apimgr/mail"),
		Log:      filepath.Join(home, ".local/log/apimgr/mail"),
		Backup:   filepath.Join(home, ".local/share/Backups/apimgr/mail"),
		PID:      filepath.Join(home, ".local/share/apimgr/mail/mail.pid"),
		SSL:      filepath.Join(home, ".config/apimgr/mail/ssl"),
		Security: filepath.Join(home, ".config/apimgr/mail/security"),
		DBDir:    filepath.Join(home, ".local/share/apimgr/mail/db"),
	}
}

// getDarwinPaths returns macOS-specific paths
// Per AI.md PART 4: macOS section
func getDarwinPaths(isRoot bool) *Paths {
	if isRoot {
		return &Paths{
			Config:   "/Library/Application Support/apimgr/mail",
			Data:     "/Library/Application Support/apimgr/mail/data",
			Cache:    "/Library/Caches/apimgr/mail",
			Log:      "/Library/Logs/apimgr/mail",
			Backup:   "/Library/Backups/apimgr/mail",
			PID:      "/var/run/apimgr/mail.pid",
			SSL:      "/Library/Application Support/apimgr/mail/ssl",
			Security: "/Library/Application Support/apimgr/mail/security",
			DBDir:    "/Library/Application Support/apimgr/mail/db",
		}
	}
	
	home := os.Getenv("HOME")
	return &Paths{
		Config:   filepath.Join(home, "Library/Application Support/apimgr/mail"),
		Data:     filepath.Join(home, "Library/Application Support/apimgr/mail"),
		Cache:    filepath.Join(home, "Library/Caches/apimgr/mail"),
		Log:      filepath.Join(home, "Library/Logs/apimgr/mail"),
		Backup:   filepath.Join(home, "Library/Backups/apimgr/mail"),
		PID:      filepath.Join(home, "Library/Application Support/apimgr/mail/mail.pid"),
		SSL:      filepath.Join(home, "Library/Application Support/apimgr/mail/ssl"),
		Security: filepath.Join(home, "Library/Application Support/apimgr/mail/security"),
		DBDir:    filepath.Join(home, "Library/Application Support/apimgr/mail/db"),
	}
}

// getBSDPaths returns BSD-specific paths
// Per AI.md PART 4: BSD section
func getBSDPaths(isRoot bool) *Paths {
	if isRoot {
		return &Paths{
			Config:   "/usr/local/etc/apimgr/mail",
			Data:     "/var/db/apimgr/mail",
			Cache:    "/var/cache/apimgr/mail",
			Log:      "/var/log/apimgr/mail",
			Backup:   "/var/backups/apimgr/mail",
			PID:      "/var/run/apimgr/mail.pid",
			SSL:      "/usr/local/etc/apimgr/mail/ssl",
			Security: "/usr/local/etc/apimgr/mail/security",
			DBDir:    "/var/db/apimgr/mail/db",
		}
	}
	
	home := os.Getenv("HOME")
	return &Paths{
		Config:   filepath.Join(home, ".config/apimgr/mail"),
		Data:     filepath.Join(home, ".local/share/apimgr/mail"),
		Cache:    filepath.Join(home, ".cache/apimgr/mail"),
		Log:      filepath.Join(home, ".local/log/apimgr/mail"),
		Backup:   filepath.Join(home, ".local/share/Backups/apimgr/mail"),
		PID:      filepath.Join(home, ".local/share/apimgr/mail/mail.pid"),
		SSL:      filepath.Join(home, ".config/apimgr/mail/ssl"),
		Security: filepath.Join(home, ".config/apimgr/mail/security"),
		DBDir:    filepath.Join(home, ".local/share/apimgr/mail/db"),
	}
}

// getWindowsPaths returns Windows-specific paths
// Per AI.md PART 4: Windows section
func getWindowsPaths(isRoot bool) *Paths {
	if isRoot {
		programData := os.Getenv("ProgramData")
		return &Paths{
			Config:   filepath.Join(programData, "apimgr/mail"),
			Data:     filepath.Join(programData, "apimgr/mail/data"),
			Cache:    filepath.Join(programData, "apimgr/mail/cache"),
			Log:      filepath.Join(programData, "apimgr/mail/logs"),
			Backup:   filepath.Join(programData, "Backups/apimgr/mail"),
			PID:      filepath.Join(programData, "apimgr/mail/mail.pid"),
			SSL:      filepath.Join(programData, "apimgr/mail/ssl"),
			Security: filepath.Join(programData, "apimgr/mail/security"),
			DBDir:    filepath.Join(programData, "apimgr/mail/db"),
		}
	}
	
	appData := os.Getenv("AppData")
	localAppData := os.Getenv("LocalAppData")
	return &Paths{
		Config:   filepath.Join(appData, "apimgr/mail"),
		Data:     filepath.Join(localAppData, "apimgr/mail"),
		Cache:    filepath.Join(localAppData, "apimgr/mail/cache"),
		Log:      filepath.Join(localAppData, "apimgr/mail/logs"),
		Backup:   filepath.Join(localAppData, "Backups/apimgr/mail"),
		PID:      filepath.Join(localAppData, "apimgr/mail/mail.pid"),
		SSL:      filepath.Join(appData, "apimgr/mail/ssl"),
		Security: filepath.Join(appData, "apimgr/mail/security"),
		DBDir:    filepath.Join(localAppData, "apimgr/mail/db"),
	}
}

// EnsureDir creates directory with proper permissions if it doesn't exist
// Per AI.md PART 8: Directory Validation Rules
func EnsureDir(path string, isRoot bool) error {
	perm := os.FileMode(0700)
	if isRoot {
		perm = 0755
	}
	
	if err := os.MkdirAll(path, perm); err != nil {
		return err
	}
	
	// Verify writable
	testFile := filepath.Join(path, ".write-test")
	if err := os.WriteFile(testFile, []byte{}, 0600); err != nil {
		return err
	}
	os.Remove(testFile)
	
	return nil
}

// EnsureAllDirs creates all application directories
// Per AI.md PART 8: Step 8b/9 - Create ALL directories
func (p *Paths) EnsureAllDirs(isRoot bool) error {
	dirs := []string{
		p.Config,
		p.Data,
		p.Cache,
		p.Log,
		p.Backup,
		p.SSL,
		p.Security,
		p.DBDir,
		filepath.Join(p.Config, "tor"),
		filepath.Join(p.Data, "tor"),
		filepath.Join(p.Data, "tor/site"),
	}
	
	for _, dir := range dirs {
		if err := EnsureDir(dir, isRoot); err != nil {
			return err
		}
	}
	
	// Ensure PID file directory
	pidDir := filepath.Dir(p.PID)
	if err := EnsureDir(pidDir, isRoot); err != nil {
		return err
	}
	
	return nil
}
