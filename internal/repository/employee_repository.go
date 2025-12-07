package repository

import (
	"context"

	"github.com/isak/restySched/internal/domain"
)

// EmployeeRepository defines the interface for employee data operations
type EmployeeRepository interface {
	// Create creates a new employee
	Create(ctx context.Context, employee *domain.Employee) error

	// GetByID retrieves an employee by ID
	GetByID(ctx context.Context, id string) (*domain.Employee, error)

	// GetAll retrieves all employees
	GetAll(ctx context.Context) ([]domain.Employee, error)

	// GetActive retrieves all active employees
	GetActive(ctx context.Context) ([]domain.Employee, error)

	// Update updates an existing employee
	Update(ctx context.Context, employee *domain.Employee) error

	// Delete soft deletes an employee (sets active to false)
	Delete(ctx context.Context, id string) error

	// GetByEmail retrieves an employee by email
	GetByEmail(ctx context.Context, email string) (*domain.Employee, error)
}
