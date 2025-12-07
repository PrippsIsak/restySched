package handler

import (
	"net/http"
	"strconv"

	"github.com/isak/restySched/internal/domain"
	"github.com/isak/restySched/internal/service"
	"github.com/isak/restySched/web/templates"
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
		http.Error(w, "Failed to fetch employees", http.StatusInternalServerError)
		return
	}

	templates.EmployeeList(employees).Render(r.Context(), w)
}

func (h *EmployeeHandler) ShowNewForm(w http.ResponseWriter, r *http.Request) {
	templates.EmployeeForm(nil, false).Render(r.Context(), w)
}

func (h *EmployeeHandler) ShowEditForm(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	employee, err := h.service.GetEmployee(r.Context(), id)
	if err != nil {
		http.Error(w, "Employee not found", http.StatusNotFound)
		return
	}

	templates.EmployeeForm(employee, true).Render(r.Context(), w)
}

func (h *EmployeeHandler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	monthlyHours, err := strconv.Atoi(r.FormValue("monthly_hours"))
	if err != nil {
		http.Error(w, "Invalid monthly hours", http.StatusBadRequest)
		return
	}

	input := domain.EmployeeCreateInput{
		Name:            r.FormValue("name"),
		Email:           r.FormValue("email"),
		Role:            r.FormValue("role"),
		RoleDescription: r.FormValue("role_description"),
		MonthlyHours:    monthlyHours,
	}

	_, err = h.service.CreateEmployee(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("HX-Redirect", "/employees")
	w.WriteHeader(http.StatusOK)
}

func (h *EmployeeHandler) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	employee, err := h.service.GetEmployee(r.Context(), id)
	if err != nil {
		http.Error(w, "Employee not found", http.StatusNotFound)
		return
	}

	monthlyHours, err := strconv.Atoi(r.FormValue("monthly_hours"))
	if err != nil {
		http.Error(w, "Invalid monthly hours", http.StatusBadRequest)
		return
	}

	employee.Name = r.FormValue("name")
	employee.Email = r.FormValue("email")
	employee.Role = r.FormValue("role")
	employee.RoleDescription = r.FormValue("role_description")
	employee.MonthlyHours = monthlyHours

	if err := h.service.UpdateEmployee(r.Context(), employee); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("HX-Redirect", "/employees")
	w.WriteHeader(http.StatusOK)
}

func (h *EmployeeHandler) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.service.DeleteEmployee(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
