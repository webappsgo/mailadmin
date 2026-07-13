package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "modernc.org/sqlite"

	"github.com/webappsgo/mail/src/config"
)

// DB wraps database connections per AI.md PART 10
// Single instance: server.db + users.db (SQLite)
// Cluster mode: shared PostgreSQL/MySQL connection
type DB struct {
	ServerDB *sql.DB // server.db (config, sessions, audit, cluster)
	UsersDB  *sql.DB // users.db (admins, users, mail, calendar, contacts)
	driver   string
}

// Connect establishes database connections based on configuration
// Per AI.md PART 10: Connection pooling with timeouts
func Connect(cfg *config.Config) (*DB, error) {
	db := &DB{driver: cfg.Database.Driver}

	// Connect to server database
	serverDB, err := openDB(cfg.Database.DSN(), cfg.Database.Driver)
	if err != nil {
		return nil, fmt.Errorf("server database connection failed: %w", err)
	}
	db.ServerDB = serverDB

	// Per AI.md PART 10: SQLite uses separate files, PostgreSQL/MySQL share connection
	if cfg.Database.Driver == "sqlite" {
		// Use separate users.db file for SQLite
		serverDSN := cfg.Database.DSN()
		usersPath := strings.Replace(serverDSN, "server.db", "users.db", 1)
		usersDB, err := openDB(usersPath, "sqlite")
		if err != nil {
			serverDB.Close()
			return nil, fmt.Errorf("users database connection failed: %w", err)
		}
		db.UsersDB = usersDB
	} else {
		// PostgreSQL/MySQL: Share connection (different tables)
		db.UsersDB = serverDB
	}

	// Configure connection pools per AI.md PART 10
	configurePool(db.ServerDB, cfg.Database.Driver)
	if cfg.Database.Driver == "sqlite" && db.UsersDB != db.ServerDB {
		configurePool(db.UsersDB, cfg.Database.Driver)
	}

	return db, nil
}

// openDB establishes a single database connection with timeout
// Per AI.md PART 10: All queries MUST have timeouts
func openDB(dsn, driver string) (*sql.DB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	// Test connection with 5 second timeout per AI.md PART 10
	ctx, cancel := timeoutContext(5 * time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	return db, nil
}

// configurePool sets connection pool settings per driver type
// Per AI.md PART 10: Pool Configuration section
func configurePool(db *sql.DB, driver string) {
	switch driver {
	case "sqlite":
		// SQLite: Single connection (WAL mode allows multiple readers)
		db.SetMaxOpenConns(1)
		db.SetMaxIdleConns(1)
		db.SetConnMaxLifetime(0)

	case "pgx", "postgres", "postgresql":
		// PostgreSQL: Default pool per AI.md PART 10
		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(5)
		db.SetConnMaxLifetime(5 * time.Minute)
		db.SetConnMaxIdleTime(1 * time.Minute)

	case "mysql":
		// MySQL: Default pool per AI.md PART 10
		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(5)
		db.SetConnMaxLifetime(5 * time.Minute)
		db.SetConnMaxIdleTime(1 * time.Minute)
	}
}

// Close closes all database connections
func (db *DB) Close() error {
	var errs []string

	if db.ServerDB != nil {
		if err := db.ServerDB.Close(); err != nil {
			errs = append(errs, fmt.Sprintf("server db: %v", err))
		}
	}

	// Only close users DB if it's separate (SQLite)
	if db.UsersDB != nil && db.UsersDB != db.ServerDB {
		if err := db.UsersDB.Close(); err != nil {
			errs = append(errs, fmt.Sprintf("users db: %v", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("close errors: %s", strings.Join(errs, "; "))
	}
	return nil
}

// EnsureSchema creates all database tables and applies updates
// Per AI.md PART 10: "All schema changes are idempotent and run on every startup"
// "CREATE TABLE IF NOT EXISTS for self-creating schema"
// "No migration files. No version tracking."
func (db *DB) EnsureSchema() error {
	// Create server database schema
	if err := db.ensureServerSchema(); err != nil {
		return fmt.Errorf("server db schema: %w", err)
	}

	// Create users database schema
	if err := db.ensureUsersSchema(); err != nil {
		return fmt.Errorf("users db schema: %w", err)
	}

	return nil
}

// ensureServerSchema creates server.db tables
// Per AI.md PART 10: Idempotent - safe to run multiple times
func (db *DB) ensureServerSchema() error {
	// 1. Create tables (idempotent with CREATE TABLE IF NOT EXISTS)
	for _, stmt := range serverCreateStatements {
		if err := db.execSchema(db.ServerDB, stmt); err != nil {
			return fmt.Errorf("create table: %w", err)
		}
	}

	// 2. Apply schema updates (idempotent ALTER TABLE)
	for _, stmt := range serverSchemaUpdates {
		if err := db.execSchema(db.ServerDB, stmt); err != nil {
			// Per AI.md PART 10: Ignore "column already exists" errors
			if !isColumnExistsError(err) {
				return fmt.Errorf("schema update: %w", err)
			}
		}
	}

	return nil
}

// ensureUsersSchema creates users.db tables
// Per AI.md PART 10: Idempotent - safe to run multiple times
func (db *DB) ensureUsersSchema() error {
	// 1. Create tables (idempotent with CREATE TABLE IF NOT EXISTS)
	for _, stmt := range usersCreateStatements {
		if err := db.execSchema(db.UsersDB, stmt); err != nil {
			return fmt.Errorf("create table: %w", err)
		}
	}

	// 2. Apply schema updates (idempotent ALTER TABLE)
	for _, stmt := range usersSchemaUpdates {
		if err := db.execSchema(db.UsersDB, stmt); err != nil {
			// Per AI.md PART 10: Ignore "column already exists" errors
			if !isColumnExistsError(err) {
				return fmt.Errorf("schema update: %w", err)
			}
		}
	}

	return nil
}

// execSchema executes a schema statement with timeout
// Per AI.md PART 10: Migrations have 5 minute timeout
func (db *DB) execSchema(conn *sql.DB, statement string) error {
	ctx, cancel := timeoutContext(5 * time.Minute)
	defer cancel()

	_, err := conn.ExecContext(ctx, statement)
	if err == context.DeadlineExceeded {
		return errors.New("TIMEOUT: Schema statement timed out")
	}
	return err
}

// WithTransaction executes a function within a transaction
// Per AI.md PART 10: Transaction Patterns section
func (db *DB) WithTransaction(ctx context.Context, conn *sql.DB, fn func(*sql.Tx) error) error {
	// Transaction timeout: 30 seconds per AI.md PART 10
	ctx, cancel := timeoutContext(30 * time.Second)
	defer cancel()

	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// WithSerializableRetry executes a function with serializable isolation and retry logic
// Per AI.md PART 10: Retry on Serialization Failure section
func (db *DB) WithSerializableRetry(ctx context.Context, conn *sql.DB, maxRetries int, fn func(*sql.Tx) error) error {
	for attempt := 0; attempt < maxRetries; attempt++ {
		tx, err := conn.BeginTx(ctx, &sql.TxOptions{
			Isolation: sql.LevelSerializable,
		})
		if err != nil {
			return err
		}

		err = fn(tx)
		if err != nil {
			tx.Rollback()
			// Retry on serialization failure per AI.md PART 10
			if isSerializationError(err) && attempt < maxRetries-1 {
				time.Sleep(time.Duration(attempt*10) * time.Millisecond)
				continue
			}
			return err
		}

		if err := tx.Commit(); err != nil {
			if isSerializationError(err) && attempt < maxRetries-1 {
				continue
			}
			return err
		}
		return nil
	}
	return errors.New("max retries exceeded")
}
