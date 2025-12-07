package repository

import (
	"context"
	"time"

	"github.com/isak/restySched/internal/domain"
)

// ScheduleRepository defines the interface for schedule data operations
type ScheduleRepository interface {
	// Create creates a new schedule
	Create(ctx context.Context, schedule *domain.Schedule) error

	// GetByID retrieves a schedule by ID
	GetByID(ctx context.Context, id string) (*domain.Schedule, error)

	// GetAll retrieves all schedules
	GetAll(ctx context.Context) ([]domain.Schedule, error)

	// GetByPeriod retrieves schedules for a specific period
	GetByPeriod(ctx context.Context, start, end time.Time) ([]domain.Schedule, error)

	// Update updates an existing schedule
	Update(ctx context.Context, schedule *domain.Schedule) error

	// Delete deletes a schedule
	Delete(ctx context.Context, id string) error

	// MarkAsSent marks a schedule as sent to n8n
	MarkAsSent(ctx context.Context, id string) error
}
