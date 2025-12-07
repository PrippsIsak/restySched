package service

import (
	"context"
	"fmt"
	"time"

	"github.com/isak/restySched/internal/domain"
	"github.com/isak/restySched/internal/n8n"
	"github.com/isak/restySched/internal/repository"
)

// ScheduleService handles business logic for schedules
type ScheduleService struct {
	scheduleRepo repository.ScheduleRepository
	employeeRepo repository.EmployeeRepository
	n8nClient    n8n.Client
}

// NewScheduleService creates a new schedule service
func NewScheduleService(
	scheduleRepo repository.ScheduleRepository,
	employeeRepo repository.EmployeeRepository,
	n8nClient n8n.Client,
) *ScheduleService {
	return &ScheduleService{
		scheduleRepo: scheduleRepo,
		employeeRepo: employeeRepo,
		n8nClient:    n8nClient,
	}
}

// GenerateSchedule generates a new schedule for the given period
func (s *ScheduleService) GenerateSchedule(ctx context.Context, periodStart, periodEnd time.Time) (*domain.Schedule, error) {
	if periodEnd.Before(periodStart) {
		return nil, domain.ErrInvalidSchedulePeriod
	}

	// Get all active employees
	employees, err := s.employeeRepo.GetActive(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get active employees: %w", err)
	}

	schedule := &domain.Schedule{
		PeriodStart: periodStart,
		PeriodEnd:   periodEnd,
		Employees:   employees,
		Status:      domain.ScheduleStatusDraft,
		SentToN8N:   false,
	}

	if err := s.scheduleRepo.Create(ctx, schedule); err != nil {
		return nil, fmt.Errorf("failed to create schedule: %w", err)
	}

	return schedule, nil
}

// GenerateBiweeklySchedule generates a schedule for the next 2 weeks
func (s *ScheduleService) GenerateBiweeklySchedule(ctx context.Context) (*domain.Schedule, error) {
	now := time.Now()
	periodStart := now
	periodEnd := now.AddDate(0, 0, 14) // 2 weeks

	return s.GenerateSchedule(ctx, periodStart, periodEnd)
}

// SendScheduleToN8N sends a schedule to the n8n webhook
func (s *ScheduleService) SendScheduleToN8N(ctx context.Context, scheduleID string) error {
	schedule, err := s.scheduleRepo.GetByID(ctx, scheduleID)
	if err != nil {
		return err
	}

	if schedule.SentToN8N {
		return domain.ErrScheduleAlreadySent
	}

	// Convert to n8n payload
	payload := s.buildN8NPayload(schedule)

	// Send to n8n
	if err := s.n8nClient.SendSchedule(ctx, payload); err != nil {
		return fmt.Errorf("failed to send schedule to n8n: %w", err)
	}

	// Mark as sent
	if err := s.scheduleRepo.MarkAsSent(ctx, scheduleID); err != nil {
		return fmt.Errorf("failed to mark schedule as sent: %w", err)
	}

	return nil
}

// GetSchedule retrieves a schedule by ID
func (s *ScheduleService) GetSchedule(ctx context.Context, id string) (*domain.Schedule, error) {
	return s.scheduleRepo.GetByID(ctx, id)
}

// GetAllSchedules retrieves all schedules
func (s *ScheduleService) GetAllSchedules(ctx context.Context) ([]domain.Schedule, error) {
	return s.scheduleRepo.GetAll(ctx)
}

// GetSchedulesByPeriod retrieves schedules for a specific period
func (s *ScheduleService) GetSchedulesByPeriod(ctx context.Context, start, end time.Time) ([]domain.Schedule, error) {
	return s.scheduleRepo.GetByPeriod(ctx, start, end)
}

// DeleteSchedule deletes a schedule
func (s *ScheduleService) DeleteSchedule(ctx context.Context, id string) error {
	return s.scheduleRepo.Delete(ctx, id)
}

func (s *ScheduleService) buildN8NPayload(schedule *domain.Schedule) domain.N8NSchedulePayload {
	employees := make([]domain.N8NEmployeeData, len(schedule.Employees))
	for i, emp := range schedule.Employees {
		employees[i] = domain.N8NEmployeeData{
			ID:              emp.ID,
			Name:            emp.Name,
			Email:           emp.Email,
			Role:            emp.Role,
			RoleDescription: emp.RoleDescription,
			MonthlyHours:    emp.MonthlyHours,
		}
	}

	return domain.N8NSchedulePayload{
		ScheduleID:  schedule.ID,
		PeriodStart: schedule.PeriodStart.Format(time.RFC3339),
		PeriodEnd:   schedule.PeriodEnd.Format(time.RFC3339),
		Employees:   employees,
		GeneratedAt: time.Now().Format(time.RFC3339),
	}
}
