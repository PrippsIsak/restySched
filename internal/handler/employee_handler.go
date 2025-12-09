package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/isak/restySched/internal/domain"
	"github.com/isak/restySched/internal/service"
	"github.com/isak/restySched/web/templates"
	"github.com/rs/zerolog/log"
)

type EmployeeHandler struct {
	service *service.EmployeeService
}

func NewEmployeeHandler(service *service.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{service: service}
}

func (h *EmployeeHandler) ListEmployees(w http.ResponseWriter, r *http.Request) {
	employees, err := h.service.GetAllEmployees(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch employees")
		handleInternalError(w, err, "fetch employees")
		return
	}

	if err := templates.EmployeeList(employees).Render(r.Context(), w); err != nil {
		log.Error().Err(err).Msg("Failed to render employee list")
		handleInternalError(w, err, "render template")
	}
}

func (h *EmployeeHandler) ShowNewForm(w http.ResponseWriter, r *http.Request) {
	if err := templates.EmployeeForm(nil, false).Render(r.Context(), w); err != nil {
		log.Error().Err(err).Msg("Failed to render employee form")
		handleInternalError(w, err, "render template")
	}
}

func (h *EmployeeHandler) ShowEditForm(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	employee, err := h.service.GetEmployee(r.Context(), id)
	if err != nil {
		log.Warn().Err(err).Str("id", id).Msg("Employee not found")
		respondWithError(w, err, http.StatusNotFound)
		return
	}

	if err := templates.EmployeeForm(employee, true).Render(r.Context(), w); err != nil {
		log.Error().Err(err).Msg("Failed to render employee form")
		handleInternalError(w, err, "render template")
	}
}

func (h *EmployeeHandler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Warn().Err(err).Msg("Invalid form data")
		respondWithError(w, domain.ErrInvalidEmployeeName, http.StatusBadRequest)
		return
	}

	monthlyHours, err := strconv.Atoi(r.FormValue("monthly_hours"))
	if err != nil {
		log.Warn().Err(err).Msg("Invalid monthly hours format")
		respondWithError(w, domain.ErrInvalidMonthlyHours, http.StatusBadRequest)
		return
	}

	input := domain.EmployeeCreateInput{
		Name:            r.FormValue("name"),
		Email:           r.FormValue("email"),
		Role:            r.FormValue("role"),
		RoleDescription: r.FormValue("role_description"),
		MonthlyHours:    monthlyHours,
	}

	employee, err := h.service.CreateEmployee(r.Context(), input)
	if err != nil {
		log.Warn().
			Err(err).
			Str("email", input.Email).
			Msg("Failed to create employee")
		respondWithError(w, err, http.StatusBadRequest)
		return
	}

	log.Info().
		Str("id", employee.ID).
		Str("name", employee.Name).
		Str("email", employee.Email).
		Msg("Employee created successfully")

	w.Header().Set("HX-Redirect", "/employees")
	w.WriteHeader(http.StatusOK)
}

func (h *EmployeeHandler) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := r.ParseForm(); err != nil {
		log.Warn().Err(err).Msg("Invalid form data")
		respondWithError(w, domain.ErrInvalidEmployeeName, http.StatusBadRequest)
		return
	}

	employee, err := h.service.GetEmployee(r.Context(), id)
	if err != nil {
		log.Warn().Err(err).Str("id", id).Msg("Employee not found")
		respondWithError(w, err, http.StatusNotFound)
		return
	}

	monthlyHours, err := strconv.Atoi(r.FormValue("monthly_hours"))
	if err != nil {
		log.Warn().Err(err).Msg("Invalid monthly hours format")
		respondWithError(w, domain.ErrInvalidMonthlyHours, http.StatusBadRequest)
		return
	}

	employee.Name = r.FormValue("name")
	employee.Email = r.FormValue("email")
	employee.Role = r.FormValue("role")
	employee.RoleDescription = r.FormValue("role_description")
	employee.MonthlyHours = monthlyHours

	if err := h.service.UpdateEmployee(r.Context(), employee); err != nil {
		log.Warn().
			Err(err).
			Str("id", id).
			Msg("Failed to update employee")
		respondWithError(w, err, http.StatusBadRequest)
		return
	}

	log.Info().
		Str("id", employee.ID).
		Str("name", employee.Name).
		Msg("Employee updated successfully")

	w.Header().Set("HX-Redirect", "/employees")
	w.WriteHeader(http.StatusOK)
}

func (h *EmployeeHandler) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.service.DeleteEmployee(r.Context(), id); err != nil {
		log.Warn().Err(err).Str("id", id).Msg("Failed to delete employee")
		respondWithError(w, err, http.StatusInternalServerError)
		return
	}

	log.Info().Str("id", id).Msg("Employee deleted successfully")
	w.WriteHeader(http.StatusOK)
}

// ShowAvailabilityManager shows the availability management UI for an employee
func (h *EmployeeHandler) ShowAvailabilityManager(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	employee, err := h.service.GetEmployee(r.Context(), id)
	if err != nil {
		log.Warn().Err(err).Str("id", id).Msg("Employee not found")
		respondWithError(w, err, http.StatusNotFound)
		return
	}

	if err := templates.AvailabilityManager(*employee).Render(r.Context(), w); err != nil {
		log.Error().Err(err).Msg("Failed to render availability manager")
		handleInternalError(w, err, "render template")
	}
}

// AddAvailability adds a new availability period for an employee
func (h *EmployeeHandler) AddAvailability(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := r.ParseForm(); err != nil {
		log.Warn().Err(err).Msg("Invalid form data")
		respondWithError(w, domain.ErrInvalidEmployeeName, http.StatusBadRequest)
		return
	}

	// Parse dates
	startDate, err := time.Parse("2006-01-02", r.FormValue("start_date"))
	if err != nil {
		log.Warn().Err(err).Msg("Invalid start date format")
		http.Error(w, "Invalid start date format", http.StatusBadRequest)
		return
	}

	endDate, err := time.Parse("2006-01-02", r.FormValue("end_date"))
	if err != nil {
		log.Warn().Err(err).Msg("Invalid end date format")
		http.Error(w, "Invalid end date format", http.StatusBadRequest)
		return
	}

	// Validate date range
	if endDate.Before(startDate) {
		log.Warn().Msg("End date is before start date")
		http.Error(w, "End date must be after start date", http.StatusBadRequest)
		return
	}

	// Parse shift types (multiple select)
	shiftTypes := r.Form["shift_types"]

	availability := domain.Availability{
		StartDate:  startDate,
		EndDate:    endDate,
		Type:       r.FormValue("type"),
		Reason:     r.FormValue("reason"),
		ShiftTypes: shiftTypes,
	}

	// Validate availability type
	if availability.Type != domain.AvailabilityTypeUnavailable &&
		availability.Type != domain.AvailabilityTypePreferred &&
		availability.Type != domain.AvailabilityTypeAvailable {
		log.Warn().Str("type", availability.Type).Msg("Invalid availability type")
		http.Error(w, "Invalid availability type", http.StatusBadRequest)
		return
	}

	employee, err := h.service.AddEmployeeAvailability(r.Context(), id, availability)
	if err != nil {
		log.Warn().
			Err(err).
			Str("id", id).
			Msg("Failed to add availability")
		respondWithError(w, err, http.StatusBadRequest)
		return
	}

	log.Info().
		Str("id", id).
		Str("type", availability.Type).
		Time("start", startDate).
		Time("end", endDate).
		Msg("Availability added successfully")

	// Return updated availability list
	if err := templates.AvailabilityList(*employee).Render(r.Context(), w); err != nil {
		log.Error().Err(err).Msg("Failed to render availability list")
		handleInternalError(w, err, "render template")
	}
}

// DeleteAvailability removes an availability period from an employee
func (h *EmployeeHandler) DeleteAvailability(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	indexStr := r.PathValue("index")

	index, err := strconv.Atoi(indexStr)
	if err != nil {
		log.Warn().Err(err).Msg("Invalid availability index")
		http.Error(w, "Invalid availability index", http.StatusBadRequest)
		return
	}

	employee, err := h.service.RemoveEmployeeAvailability(r.Context(), id, index)
	if err != nil {
		log.Warn().
			Err(err).
			Str("id", id).
			Int("index", index).
			Msg("Failed to remove availability")
		respondWithError(w, err, http.StatusBadRequest)
		return
	}

	log.Info().
		Str("id", id).
		Int("index", index).
		Msg("Availability removed successfully")

	// Return updated availability list
	if err := templates.AvailabilityList(*employee).Render(r.Context(), w); err != nil {
		log.Error().Err(err).Msg("Failed to render availability list")
		handleInternalError(w, err, "render template")
	}
}
