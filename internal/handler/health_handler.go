package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/isak/restySched/internal/repository"
	"github.com/rs/zerolog/log"
)

// HealthHandler handles health check endpoints
type HealthHandler struct {
	employeeRepo repository.EmployeeRepository
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(employeeRepo repository.EmployeeRepository) *HealthHandler {
	return &HealthHandler{
		employeeRepo: employeeRepo,
	}
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp string            `json:"timestamp"`
	Checks    map[string]string `json:"checks,omitempty"`
}

// Health returns basic health status (liveness probe)
// This endpoint should return 200 if the application is running
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Error().Err(err).Msg("Failed to encode health response")
	}
}

// Ready returns readiness status (readiness probe)
// This endpoint checks if the application can serve traffic
// It verifies database connectivity and other critical dependencies
func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	checks := make(map[string]string)
	allHealthy := true

	// Check database connectivity
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	// Try a simple database query to verify connectivity
	_, err := h.employeeRepo.GetAll(ctx)
	if err != nil {
		checks["database"] = "unhealthy"
		allHealthy = false
		log.Warn().Err(err).Msg("Database health check failed")
	} else {
		checks["database"] = "healthy"
	}

	status := "ready"
	statusCode := http.StatusOK

	if !allHealthy {
		status = "not_ready"
		statusCode = http.StatusServiceUnavailable
	}

	response := HealthResponse{
		Status:    status,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Checks:    checks,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Error().Err(err).Msg("Failed to encode readiness response")
	}
}
