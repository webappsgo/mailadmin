package server

import (
	"encoding/json"
	"net/http"
	"runtime"
)

// handleHealth handles /healthz endpoint
// Per AI.md PART 13: Health endpoint for monitoring
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement full health check per AI.md PART 13
	// For now, basic response
	
	health := map[string]interface{}{
		"status": "healthy",
		"version": "dev", // TODO: Get from build info
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(health)
}

// handleAbout handles /api/v1/server/about endpoint
// Per AI.md PART 14: Server info routes
func (s *Server) handleAbout(w http.ResponseWriter, r *http.Request) {
	// TODO: Read from IDEA.md per AI.md PART 0
	about := map[string]interface{}{
		"name":        "mail",
		"description": "Email Infrastructure Management Panel",
		"version":     "dev", // TODO: Get from build info
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(about)
}

// handleVersion handles /api/v1/server/version endpoint
// Per AI.md PART 14: Server info routes
func (s *Server) handleVersion(w http.ResponseWriter, r *http.Request) {
	// TODO: Get build info from main package
	version := map[string]interface{}{
		"version":    "dev",
		"commit":     "unknown",
		"build_date": "unknown",
		"go_version": runtime.Version(),
		"platform":   runtime.GOOS + "/" + runtime.GOARCH,
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(version)
}
