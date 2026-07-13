package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the complete server configuration
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	SSL      SSLConfig      `yaml:"ssl"`
	Email    EmailConfig    `yaml:"email"`
}

// ServerConfig holds server-specific settings
type ServerConfig struct {
	Port          int      `yaml:"port"`
	Address       string   `yaml:"address"`
	FQDN          string   `yaml:"fqdn"`
	AdminPath     string   `yaml:"admin_path"`
	SetupToken    string   `yaml:"setup_token"`
	Maintenance   bool     `yaml:"maintenance"`
	Debug         bool     `yaml:"debug"`
	TrustedProxies []string `yaml:"trusted_proxies"`
	
	Limits struct {
		MaxBodySize  string `yaml:"max_body_size"`
		ReadTimeout  string `yaml:"read_timeout"`
		WriteTimeout string `yaml:"write_timeout"`
		IdleTimeout  string `yaml:"idle_timeout"`
	} `yaml:"limits"`
	
	Session struct {
		Admin struct {
			CookieName       string `yaml:"cookie_name"`
			MaxAge           string `yaml:"max_age"`
			IdleTimeout      string `yaml:"idle_timeout"`
			ExtendOnActivity bool   `yaml:"extend_on_activity"`
		} `yaml:"admin"`
		User struct {
			CookieName       string `yaml:"cookie_name"`
			MaxAge           string `yaml:"max_age"`
			IdleTimeout      string `yaml:"idle_timeout"`
			ExtendOnActivity bool   `yaml:"extend_on_activity"`
		} `yaml:"user"`
		Secure   string `yaml:"secure"`
		HTTPOnly bool   `yaml:"http_only"`
		SameSite string `yaml:"same_site"`
	} `yaml:"session"`
	
	RateLimit struct {
		Enabled  bool `yaml:"enabled"`
		Requests int  `yaml:"requests"`
		Window   int  `yaml:"window"`
	} `yaml:"rate_limit"`
	
	I18n struct {
		DefaultLanguage string   `yaml:"default_language"`
		Supported       []string `yaml:"supported"`
	} `yaml:"i18n"`
	
	Tracking struct {
		Type string `yaml:"type"`
		ID   string `yaml:"id"`
		URL  string `yaml:"url"`
	} `yaml:"tracking"`
	
	Privacy struct {
		Data struct {
			Sold            bool `yaml:"sold"`
			StoredOnServer  bool `yaml:"stored_on_server"`
		} `yaml:"data"`
	} `yaml:"privacy"`
}

// DatabaseConfig holds database connection settings
type DatabaseConfig struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Name     string `yaml:"name"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	URL      string `yaml:"url"`
	Token    string `yaml:"token"`
	
	Pool struct {
		MaxOpen      int    `yaml:"max_open"`
		MaxIdle      int    `yaml:"max_idle"`
		MaxLifetime  string `yaml:"max_lifetime"`
		MaxIdleTime  string `yaml:"max_idle_time"`
	} `yaml:"pool"`
}

// SSLConfig holds SSL/TLS and Let's Encrypt settings
type SSLConfig struct {
	Enabled       bool   `yaml:"enabled"`
	CertFile      string `yaml:"cert_file"`
	KeyFile       string `yaml:"key_file"`
	LetsEncrypt   bool   `yaml:"lets_encrypt"`
	Email         string `yaml:"email"`
	Staging       bool   `yaml:"staging"`
}

// EmailConfig holds SMTP settings for outgoing mail
type EmailConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	From     string `yaml:"from"`
	UseTLS   bool   `yaml:"use_tls"`
}

var (
	// Global config instance
	globalConfig *Config
	configPath   string
)

// Load loads configuration from YAML file with environment variable expansion
func Load(path string) (*Config, error) {
	configPath = path
	
	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}
	
	// Expand environment variables
	expanded := os.ExpandEnv(string(data))
	
	// Parse YAML
	cfg := &Config{}
	if err := yaml.Unmarshal([]byte(expanded), cfg); err != nil {
		return nil, fmt.Errorf("parse config YAML: %w", err)
	}
	
	// Validate and set defaults
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}
	
	globalConfig = cfg
	return cfg, nil
}

// validate validates and sets defaults for configuration
func (c *Config) validate() error {
	// Server defaults
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		c.Server.Port = getRandomAvailablePort()
	}
	if c.Server.Address == "" {
		c.Server.Address = "[::]" // IPv6 any (includes IPv4)
	}
	if c.Server.AdminPath == "" {
		c.Server.AdminPath = "admin"
	}
	
	// Normalize admin path
	c.Server.AdminPath = strings.Trim(c.Server.AdminPath, "/")
	
	// Database defaults
	if c.Database.Driver == "" {
		c.Database.Driver = "sqlite"
	}
	
	// Normalize database driver
	c.Database.Driver = normalizeDriver(c.Database.Driver)
	
	// Pool defaults
	if c.Database.Pool.MaxOpen <= 0 {
		c.Database.Pool.MaxOpen = 25
	}
	if c.Database.Pool.MaxIdle <= 0 {
		c.Database.Pool.MaxIdle = 5
	}
	if c.Database.Pool.MaxLifetime == "" {
		c.Database.Pool.MaxLifetime = "5m"
	}
	if c.Database.Pool.MaxIdleTime == "" {
		c.Database.Pool.MaxIdleTime = "1m"
	}
	
	// Session defaults
	if c.Server.Session.Admin.CookieName == "" {
		c.Server.Session.Admin.CookieName = "admin_session"
	}
	if c.Server.Session.Admin.MaxAge == "" {
		c.Server.Session.Admin.MaxAge = "30d"
	}
	if c.Server.Session.Admin.IdleTimeout == "" {
		c.Server.Session.Admin.IdleTimeout = "24h"
	}
	if c.Server.Session.User.CookieName == "" {
		c.Server.Session.User.CookieName = "user_session"
	}
	if c.Server.Session.User.MaxAge == "" {
		c.Server.Session.User.MaxAge = "7d"
	}
	if c.Server.Session.User.IdleTimeout == "" {
		c.Server.Session.User.IdleTimeout = "24h"
	}
	
	// i18n defaults
	if c.Server.I18n.DefaultLanguage == "" {
		c.Server.I18n.DefaultLanguage = "en"
	}
	if len(c.Server.I18n.Supported) == 0 {
		c.Server.I18n.Supported = []string{"en"}
	}
	
	return nil
}

// normalizeDriver normalizes database driver names per AI.md PART 3
func normalizeDriver(driver string) string {
	switch strings.ToLower(driver) {
	case "sqlite", "sqlite2", "sqlite3":
		return "sqlite"
	case "libsql", "turso":
		return "libsql"
	case "postgres", "pgsql", "postgresql":
		return "pgx"
	case "mysql", "mariadb":
		return "mysql"
	case "mssql":
		return "sqlserver"
	case "mongodb", "mongo":
		return "mongodb"
	default:
		return driver
	}
}

// getRandomAvailablePort returns a random port in 64000-64999 range
func getRandomAvailablePort() int {
	// TODO: Implement actual port availability check
	return 64500 // Default for now
}

// Get returns the global config instance
func Get() *Config {
	return globalConfig
}

// ParseDuration parses duration string (e.g., "30s", "5m", "24h")
func ParseDuration(s string) (time.Duration, error) {
	return time.ParseDuration(s)
}

// GetConfigPath returns the current config file path
func GetConfigPath() string {
	return configPath
}

// DSN generates a database connection string based on driver type
func (d *DatabaseConfig) DSN() string {
	switch d.Driver {
	case "sqlite":
		if d.URL != "" {
			return d.URL
		}
		if d.Name != "" {
			return d.Name
		}
		return filepath.Join("/var/lib/mail", "server.db")
	
	case "libsql":
		// libSQL requires URL
		if d.Token != "" && !strings.Contains(d.URL, "authToken=") {
			sep := "?"
			if strings.Contains(d.URL, "?") {
				sep = "&"
			}
			return d.URL + sep + "authToken=" + d.Token
		}
		return d.URL
	
	case "pgx":
		// PostgreSQL connection string
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=prefer",
			d.Host, d.Port, d.Username, d.Password, d.Name)
	
	case "mysql":
		// MySQL connection string
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			d.Username, d.Password, d.Host, d.Port, d.Name)
	
	default:
		return ""
	}
}
