package handler

import (
	"net/http"

	"github.com/isak/restySched/internal/service"
	"github.com/isak/restySched/web/templates"
	"github.com/rs/zerolog/log"
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
		log.Error().Err(err).Msg("Failed to fetch schedules")
		handleInternalError(w, err, "fetch schedules")
		return
	}

	if err := templates.ScheduleList(schedules).Render(r.Context(), w); err != nil {
		log.Error().Err(err).Msg("Failed to render schedule list")
		handleInternalError(w, err, "render template")
	}
}

func (h *ScheduleHandler) GenerateBiweeklySchedule(w http.ResponseWriter, r *http.Request) {
	schedule, err := h.service.GenerateBiweeklySchedule(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate schedule")
		respondWithError(w, err, http.StatusInternalServerError)
		return
	}

	log.Info().
		Str("schedule_id", schedule.ID).
		Time("period_start", schedule.PeriodStart).
		Time("period_end", schedule.PeriodEnd).
		Int("employees", len(schedule.Employees)).
		Msg("Schedule generated successfully")

	if err := templates.ScheduleCard(*schedule).Render(r.Context(), w); err != nil {
		log.Error().Err(err).Msg("Failed to render schedule card")
		handleInternalError(w, err, "render template")
	}
}

func (h *ScheduleHandler) SendToN8N(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.service.SendScheduleToN8N(r.Context(), id); err != nil {
		log.Warn().
			Err(err).
			Str("schedule_id", id).
			Msg("Failed to send schedule to n8n")
		respondWithError(w, err, http.StatusInternalServerError)
		return
	}

	log.Info().
		Str("schedule_id", id).
		Msg("Schedule sent to n8n successfully")

	// Fetch updated schedule
	schedule, err := h.service.GetSchedule(r.Context(), id)
	if err != nil {
		log.Error().Err(err).Str("schedule_id", id).Msg("Failed to fetch updated schedule")
		handleInternalError(w, err, "fetch schedule")
		return
	}

	if err := templates.ScheduleCard(*schedule).Render(r.Context(), w); err != nil {
		log.Error().Err(err).Msg("Failed to render schedule card")
		handleInternalError(w, err, "render template")
	}
}

func (h *ScheduleHandler) DeleteSchedule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.service.DeleteSchedule(r.Context(), id); err != nil {
		log.Warn().Err(err).Str("schedule_id", id).Msg("Failed to delete schedule")
		respondWithError(w, err, http.StatusInternalServerError)
		return
	}

	log.Info().Str("schedule_id", id).Msg("Schedule deleted successfully")
	w.WriteHeader(http.StatusOK)
}
