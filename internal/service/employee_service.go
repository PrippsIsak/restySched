package service

import (
	"context"

	"github.com/isak/restySched/internal/domain"
	"github.com/isak/restySched/internal/repository"
)

// EmployeeService handles business logic for employees
type EmployeeService struct {
	repo repository.EmployeeRepository
}

// NewEmployeeService creates a new employee service
func NewEmployeeService(repo repository.EmployeeRepository) *EmployeeService {
	return &EmployeeService{repo: repo}
}

// CreateEmployee creates a new employee
func (s *EmployeeService) CreateEmployee(ctx context.Context, input domain.EmployeeCreateInput) (*domain.Employee, error) {
	// Check if email already exists
	existing, err := s.repo.GetByEmail(ctx, input.Email)
	if err != nil && err != domain.ErrEmployeeNotFound {
		return nil, err
	}
	if existing != nil {
		return nil, domain.ErrEmployeeAlreadyExists
	}

	employee := &domain.Employee{
		Name:            input.Name,
		Email:           input.Email,
		Role:            input.Role,
		RoleDescription: input.RoleDescription,
		MonthlyHours:    input.MonthlyHours,
	}

	if err := employee.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, employee); err != nil {
		return nil, err
	}

	return employee, nil
}

// GetEmployee retrieves an employee by ID
func (s *EmployeeService) GetEmployee(ctx context.Context, id string) (*domain.Employee, error) {
	return s.repo.GetByID(ctx, id)
}

// GetAllEmployees retrieves all employees
func (s *EmployeeService) GetAllEmployees(ctx context.Context) ([]domain.Employee, error) {
	return s.repo.GetAll(ctx)
}

// GetActiveEmployees retrieves all active employees
func (s *EmployeeService) GetActiveEmployees(ctx context.Context) ([]domain.Employee, error) {
	return s.repo.GetActive(ctx)
}

// UpdateEmployee updates an existing employee
func (s *EmployeeService) UpdateEmployee(ctx context.Context, employee *domain.Employee) error {
	if err := employee.Validate(); err != nil {
		return err
	}

	return s.repo.Update(ctx, employee)
}

// DeleteEmployee soft deletes an employee
func (s *EmployeeService) DeleteEmployee(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
