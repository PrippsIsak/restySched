package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/yourusername/restysched/internal/repository"
	"github.com/yourusername/restysched/internal/domain"
	"github.com/yourusername/restysched/web/templates"
)

type CompanyConfigHandler struct {
	repo repository.CompanyConfigRepository
}

func NewCompanyConfigHandler(repo repository.CompanyConfigRepository) *CompanyConfigHandler {
	return &CompanyConfigHandler{
		repo: repo,
	}
}

func (h *CompanyConfigHandler) ShowConfig(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get or create default config
	config, err := h.repo.GetOrCreate(ctx)
	if err != nil {
		http.Error(w, "Failed to load configuration", http.StatusInternalServerError)
		return
	}

	templates.CompanyConfig(config).Render(ctx, w)
}

func (h *CompanyConfigHandler) SaveConfig(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Parse working days
	workingDays := []int{}
	for _, dayStr := range r.Form["working_days"] {
		day, err := strconv.Atoi(dayStr)
		if err == nil && day >= 0 && day <= 6 {
			workingDays = append(workingDays, day)
		}
	}

	// Parse shift requirements
	shiftReqs := []domain.ShiftRequirement{}
	i := 0
	for {
		shiftTypeKey := "shift_type_" + strconv.Itoa(i)
		if r.FormValue(shiftTypeKey) == "" {
			break
		}

		// Parse required skills (comma-separated)
		skillsStr := r.FormValue("required_skills_" + strconv.Itoa(i))
		var skills []string
		if skillsStr != "" {
			for _, skill := range strings.Split(skillsStr, ",") {
				trimmed := strings.TrimSpace(skill)
				if trimmed != "" {
					skills = append(skills, trimmed)
				}
			}
		}

		req := domain.ShiftRequirement{
			ShiftType:      r.FormValue(shiftTypeKey),
			MinEmployees:   parseInt(r.FormValue("min_employees_"+strconv.Itoa(i)), 1),
			MaxEmployees:   parseInt(r.FormValue("max_employees_"+strconv.Itoa(i)), 2),
			RequiredSkills: skills,
			Description:    r.FormValue("description_" + strconv.Itoa(i)),
		}
		shiftReqs = append(shiftReqs, req)
		i++
	}

	// Parse form data
	config := &domain.CompanyConfig{
		CompanyName: r.FormValue("company_name"),
		WorkingHours: domain.WorkingHours{
			WorkingDays: workingDays,
			OpenTime:    r.FormValue("open_time"),
			CloseTime:   r.FormValue("close_time"),
			Timezone:    r.FormValue("timezone"),
		},
		ShiftRequirements: shiftReqs,
		SchedulingPolicies: domain.SchedulingPolicies{
			MaxConsecutiveDays: parseInt(r.FormValue("max_consecutive_days"), 5),
			MinRestHours:       parseInt(r.FormValue("min_rest_hours"), 12),
			MaxHoursPerWeek:    parseInt(r.FormValue("max_hours_per_week"), 40),
			MaxShiftsPerDay:    parseInt(r.FormValue("max_shifts_per_day"), 1),
		},
		AIContext: strings.TrimSpace(r.FormValue("ai_context")),
	}

	// Validate
	if err := config.Validate(); err != nil {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`<div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">` + err.Error() + `</div>`))
		return
	}

	// Update or create
	err := h.repo.Update(ctx, config)
	if err != nil {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`<div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">Failed to save configuration</div>`))
		return
	}

	// Success response
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`
		<div class="bg-green-100 border border-green-400 text-green-700 px-4 py-3 rounded">
			Configuration saved successfully! This will be included in all schedule analysis sent to n8n.
		</div>
	`))
}

func parseInt(s string, defaultVal int) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return val
}
