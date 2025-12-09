package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/isak/restySched/internal/domain"
	"github.com/rs/zerolog/log"
)

// ErrorResponse represents a JSON error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// respondWithError sends an error response with appropriate status code
func respondWithError(w http.ResponseWriter, err error, defaultStatus int) {
	status := defaultStatus
	message := err.Error()

	// Map domain errors to HTTP status codes
	switch {
	case errors.Is(err, domain.ErrEmployeeNotFound),
		errors.Is(err, domain.ErrScheduleNotFound):
		status = http.StatusNotFound

	case errors.Is(err, domain.ErrInvalidEmployeeName),
		errors.Is(err, domain.ErrInvalidEmployeeEmail),
		errors.Is(err, domain.ErrInvalidEmployeeRole),
		errors.Is(err, domain.ErrInvalidMonthlyHours),
		errors.Is(err, domain.ErrInvalidSchedulePeriod):
		status = http.StatusBadRequest

	case errors.Is(err, domain.ErrScheduleAlreadySent):
		status = http.StatusConflict
	}

	// Log the error with context
	log.Error().
		Err(err).
		Int("status", status).
		Msg("Request error")

	// Send HTML error response (since we're using HTMX)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)

	errorHTML := fmt.Sprintf(`
		<div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative" role="alert">
			<strong class="font-bold">Error!</strong>
			<span class="block sm:inline">%s</span>
		</div>
	`, message)

	w.Write([]byte(errorHTML))
}

// respondWithSuccess sends a success message
func respondWithSuccess(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	successHTML := fmt.Sprintf(`
		<div class="bg-green-100 border border-green-400 text-green-700 px-4 py-3 rounded relative" role="alert">
			<strong class="font-bold">Success!</strong>
			<span class="block sm:inline">%s</span>
		</div>
	`, message)

	w.Write([]byte(successHTML))
}

// handleInternalError logs and responds with a generic internal error
func handleInternalError(w http.ResponseWriter, err error, context string) {
	log.Error().
		Err(err).
		Str("context", context).
		Msg("Internal server error")

	respondWithError(w, errors.New("An internal error occurred. Please try again later."), http.StatusInternalServerError)
}
