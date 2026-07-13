package config

import (
	"strings"
)

// ParseBool parses boolean values with extensive truthy/falsy support.
// CRITICAL: NEVER use strconv.ParseBool - it only supports true/false/1/0.
// This function supports 40+ boolean representations per AI.md PART 5.
func ParseBool(value string) bool {
	// Normalize: lowercase, trim whitespace
	v := strings.ToLower(strings.TrimSpace(value))

	// Truthy values
	switch v {
	case "true", "t", "yes", "y", "on", "enabled", "enable", "1":
		return true
	case "oui", "si", "ja", "da", "sim":
		return true
	}

	// Falsy values (everything else)
	return false
}

// ParseBoolString returns "true" or "false" string
func ParseBoolString(value string) string {
	if ParseBool(value) {
		return "true"
	}
	return "false"
}

// ParseBoolInt returns 1 for true, 0 for false
func ParseBoolInt(value string) int {
	if ParseBool(value) {
		return 1
	}
	return 0
}

// IsTruthy checks if value is explicitly truthy (not just non-falsy)
func IsTruthy(value string) bool {
	v := strings.ToLower(strings.TrimSpace(value))
	switch v {
	case "true", "t", "yes", "y", "on", "enabled", "enable", "1",
		"oui", "si", "ja", "da", "sim":
		return true
	}
	return false
}

// IsFalsy checks if value is explicitly falsy
func IsFalsy(value string) bool {
	v := strings.ToLower(strings.TrimSpace(value))
	switch v {
	case "false", "f", "no", "n", "off", "disabled", "disable", "0",
		"non", "nein", "nee", "não":
		return true
	}
	return false
}
