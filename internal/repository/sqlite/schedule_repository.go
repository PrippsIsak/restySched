package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/isak/restySched/internal/domain"
	"github.com/isak/restySched/internal/repository"
)

type scheduleRepository struct {
	db *sql.DB
}

// NewScheduleRepository creates a new SQLite schedule repository
func NewScheduleRepository(db *sql.DB) repository.ScheduleRepository {
	return &scheduleRepository{db: db}
}

func (r *scheduleRepository) Create(ctx context.Context, schedule *domain.Schedule) error {
	if schedule.ID == "" {
		schedule.ID = uuid.New().String()
	}

	now := time.Now()
	schedule.CreatedAt = now
	schedule.UpdatedAt = now

	employeesJSON, err := json.Marshal(schedule.Employees)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO schedules (id, period_start, period_end, employees, status, sent_to_n8n, sent_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = r.db.ExecContext(ctx, query,
		schedule.ID,
		schedule.PeriodStart,
		schedule.PeriodEnd,
		employeesJSON,
		schedule.Status,
		schedule.SentToN8N,
		schedule.SentAt,
		schedule.CreatedAt,
		schedule.UpdatedAt,
	)

	return err
}

func (r *scheduleRepository) GetByID(ctx context.Context, id string) (*domain.Schedule, error) {
	query := `
		SELECT id, period_start, period_end, employees, status, sent_to_n8n, sent_at, created_at, updated_at
		FROM schedules
		WHERE id = ?
	`

	schedule := &domain.Schedule{}
	var employeesJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&schedule.ID,
		&schedule.PeriodStart,
		&schedule.PeriodEnd,
		&employeesJSON,
		&schedule.Status,
		&schedule.SentToN8N,
		&schedule.SentAt,
		&schedule.CreatedAt,
		&schedule.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrScheduleNotFound
	}
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(employeesJSON, &schedule.Employees); err != nil {
		return nil, err
	}

	return schedule, nil
}

func (r *scheduleRepository) GetAll(ctx context.Context) ([]domain.Schedule, error) {
	query := `
		SELECT id, period_start, period_end, employees, status, sent_to_n8n, sent_at, created_at, updated_at
		FROM schedules
		ORDER BY period_start DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanSchedules(rows)
}

func (r *scheduleRepository) GetByPeriod(ctx context.Context, start, end time.Time) ([]domain.Schedule, error) {
	query := `
		SELECT id, period_start, period_end, employees, status, sent_to_n8n, sent_at, created_at, updated_at
		FROM schedules
		WHERE period_start >= ? AND period_end <= ?
		ORDER BY period_start DESC
	`

	rows, err := r.db.QueryContext(ctx, query, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanSchedules(rows)
}

func (r *scheduleRepository) Update(ctx context.Context, schedule *domain.Schedule) error {
	schedule.UpdatedAt = time.Now()

	employeesJSON, err := json.Marshal(schedule.Employees)
	if err != nil {
		return err
	}

	query := `
		UPDATE schedules
		SET period_start = ?, period_end = ?, employees = ?, status = ?, sent_to_n8n = ?, sent_at = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.ExecContext(ctx, query,
		schedule.PeriodStart,
		schedule.PeriodEnd,
		employeesJSON,
		schedule.Status,
		schedule.SentToN8N,
		schedule.SentAt,
		schedule.UpdatedAt,
		schedule.ID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrScheduleNotFound
	}

	return nil
}

func (r *scheduleRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM schedules WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrScheduleNotFound
	}

	return nil
}

func (r *scheduleRepository) MarkAsSent(ctx context.Context, id string) error {
	now := time.Now()
	query := `
		UPDATE schedules
		SET sent_to_n8n = 1, sent_at = ?, status = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.ExecContext(ctx, query, now, domain.ScheduleStatusSent, now, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrScheduleNotFound
	}

	return nil
}

func (r *scheduleRepository) scanSchedules(rows *sql.Rows) ([]domain.Schedule, error) {
	var schedules []domain.Schedule

	for rows.Next() {
		var schedule domain.Schedule
		var employeesJSON []byte

		err := rows.Scan(
			&schedule.ID,
			&schedule.PeriodStart,
			&schedule.PeriodEnd,
			&employeesJSON,
			&schedule.Status,
			&schedule.SentToN8N,
			&schedule.SentAt,
			&schedule.CreatedAt,
			&schedule.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(employeesJSON, &schedule.Employees); err != nil {
			return nil, err
		}

		schedules = append(schedules, schedule)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return schedules, nil
}
