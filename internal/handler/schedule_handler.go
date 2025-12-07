package handler

import (
	"log"
	"net/http"

	"github.com/isak/restySched/internal/service"
	"github.com/isak/restySched/web/templates"
)

type ScheduleHandler struct {
	service *service.ScheduleService
}

func NewScheduleHandler(service *service.ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{service: service}
}

func (h *ScheduleHandler) ListSchedules(w http.ResponseWriter, r *http.Request) {
	schedules, err := h.service.GetAllSchedules(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch schedules", http.StatusInternalServerError)
		return
	}

	templates.ScheduleList(schedules).Render(r.Context(), w)
}

func (h *ScheduleHandler) GenerateBiweeklySchedule(w http.ResponseWriter, r *http.Request) {
	schedule, err := h.service.GenerateBiweeklySchedule(r.Context())
	if err != nil {
		log.Printf("Failed to generate schedule: %v", err)
		http.Error(w, "Failed to generate schedule", http.StatusInternalServerError)
		return
	}

	templates.ScheduleCard(*schedule).Render(r.Context(), w)
}

func (h *ScheduleHandler) SendToN8N(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.service.SendScheduleToN8N(r.Context(), id); err != nil {
		log.Printf("Failed to send schedule to n8n: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch updated schedule
	schedule, err := h.service.GetSchedule(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to fetch updated schedule", http.StatusInternalServerError)
		return
	}

	templates.ScheduleCard(*schedule).Render(r.Context(), w)
}

func (h *ScheduleHandler) DeleteSchedule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.service.DeleteSchedule(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
