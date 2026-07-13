package database

import (
	"context"
	"strings"
	"time"
)

// timeoutContext creates a context with timeout for database operations
func timeoutContext(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}

// isColumnExistsError checks if error is "column already exists"
// Per AI.md PART 10: Handle cross-database compatibility
func isColumnExistsError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "duplicate column") || // SQLite
		strings.Contains(msg, "already exists") || // PostgreSQL
		strings.Contains(msg, "Duplicate column name") // MySQL
}

// isSerializationError checks if error is serialization failure (for retry logic)
// Per AI.md PART 10: PostgreSQL 40001, MySQL 1213
func isSerializationError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "40001") || // PostgreSQL serialization_failure
		strings.Contains(msg, "1213") // MySQL Deadlock found
}
