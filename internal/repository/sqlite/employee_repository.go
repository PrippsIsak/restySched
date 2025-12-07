package sqlite

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/isak/restySched/internal/domain"
	"github.com/isak/restySched/internal/repository"
)

type employeeRepository struct {
	db *sql.DB
}

// NewEmployeeRepository creates a new SQLite employee repository
func NewEmployeeRepository(db *sql.DB) repository.EmployeeRepository {
	return &employeeRepository{db: db}
}

func (r *employeeRepository) Create(ctx context.Context, employee *domain.Employee) error {
	if employee.ID == "" {
		employee.ID = uuid.New().String()
	}

	now := time.Now()
	employee.CreatedAt = now
	employee.UpdatedAt = now
	employee.Active = true

	query := `
		INSERT INTO employees (id, name, email, role, role_description, monthly_hours, active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		employee.ID,
		employee.Name,
		employee.Email,
		employee.Role,
		employee.RoleDescription,
		employee.MonthlyHours,
		employee.Active,
		employee.CreatedAt,
		employee.UpdatedAt,
	)

	return err
}

func (r *employeeRepository) GetByID(ctx context.Context, id string) (*domain.Employee, error) {
	query := `
		SELECT id, name, email, role, role_description, monthly_hours, active, created_at, updated_at
		FROM employees
		WHERE id = ?
	`

	employee := &domain.Employee{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&employee.ID,
		&employee.Name,
		&employee.Email,
		&employee.Role,
		&employee.RoleDescription,
		&employee.MonthlyHours,
		&employee.Active,
		&employee.CreatedAt,
		&employee.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrEmployeeNotFound
	}
	if err != nil {
		return nil, err
	}

	return employee, nil
}

func (r *employeeRepository) GetAll(ctx context.Context) ([]domain.Employee, error) {
	query := `
		SELECT id, name, email, role, role_description, monthly_hours, active, created_at, updated_at
		FROM employees
		ORDER BY name
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanEmployees(rows)
}

func (r *employeeRepository) GetActive(ctx context.Context) ([]domain.Employee, error) {
	query := `
		SELECT id, name, email, role, role_description, monthly_hours, active, created_at, updated_at
		FROM employees
		WHERE active = 1
		ORDER BY name
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanEmployees(rows)
}

func (r *employeeRepository) Update(ctx context.Context, employee *domain.Employee) error {
	employee.UpdatedAt = time.Now()

	query := `
		UPDATE employees
		SET name = ?, email = ?, role = ?, role_description = ?, monthly_hours = ?, active = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.ExecContext(ctx, query,
		employee.Name,
		employee.Email,
		employee.Role,
		employee.RoleDescription,
		employee.MonthlyHours,
		employee.Active,
		employee.UpdatedAt,
		employee.ID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrEmployeeNotFound
	}

	return nil
}

func (r *employeeRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE employees SET active = 0, updated_at = ? WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrEmployeeNotFound
	}

	return nil
}

func (r *employeeRepository) GetByEmail(ctx context.Context, email string) (*domain.Employee, error) {
	query := `
		SELECT id, name, email, role, role_description, monthly_hours, active, created_at, updated_at
		FROM employees
		WHERE email = ?
	`

	employee := &domain.Employee{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&employee.ID,
		&employee.Name,
		&employee.Email,
		&employee.Role,
		&employee.RoleDescription,
		&employee.MonthlyHours,
		&employee.Active,
		&employee.CreatedAt,
		&employee.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrEmployeeNotFound
	}
	if err != nil {
		return nil, err
	}

	return employee, nil
}

func (r *employeeRepository) scanEmployees(rows *sql.Rows) ([]domain.Employee, error) {
	var employees []domain.Employee

	for rows.Next() {
		var employee domain.Employee
		err := rows.Scan(
			&employee.ID,
			&employee.Name,
			&employee.Email,
			&employee.Role,
			&employee.RoleDescription,
			&employee.MonthlyHours,
			&employee.Active,
			&employee.CreatedAt,
			&employee.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return employees, nil
}
