package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/isak/restySched/internal/domain"
)

// MockEmployeeRepository is a mock implementation of EmployeeRepository for testing
type MockEmployeeRepository struct {
	employees map[string]*domain.Employee
	idCounter int
}

func NewMockEmployeeRepository() *MockEmployeeRepository {
	return &MockEmployeeRepository{
		employees: make(map[string]*domain.Employee),
		idCounter: 0,
	}
}

func (m *MockEmployeeRepository) Create(ctx context.Context, employee *domain.Employee) error {
	if employee.ID == "" {
		m.idCounter++
		employee.ID = fmt.Sprintf("mock-id-%d", m.idCounter)
	}
	employee.Active = true
	m.employees[employee.ID] = employee
	return nil
}

func (m *MockEmployeeRepository) GetByID(ctx context.Context, id string) (*domain.Employee, error) {
	emp, ok := m.employees[id]
	if !ok {
		return nil, domain.ErrEmployeeNotFound
	}
	return emp, nil
}

func (m *MockEmployeeRepository) GetAll(ctx context.Context) ([]domain.Employee, error) {
	var result []domain.Employee
	for _, emp := range m.employees {
		result = append(result, *emp)
	}
	return result, nil
}

func (m *MockEmployeeRepository) GetActive(ctx context.Context) ([]domain.Employee, error) {
	var result []domain.Employee
	for _, emp := range m.employees {
		if emp.Active {
			result = append(result, *emp)
		}
	}
	return result, nil
}

func (m *MockEmployeeRepository) Update(ctx context.Context, employee *domain.Employee) error {
	if _, ok := m.employees[employee.ID]; !ok {
		return domain.ErrEmployeeNotFound
	}
	m.employees[employee.ID] = employee
	return nil
}

func (m *MockEmployeeRepository) Delete(ctx context.Context, id string) error {
	emp, ok := m.employees[id]
	if !ok {
		return domain.ErrEmployeeNotFound
	}
	emp.Active = false
	return nil
}

func (m *MockEmployeeRepository) GetByEmail(ctx context.Context, email string) (*domain.Employee, error) {
	for _, emp := range m.employees {
		if emp.Email == email {
			return emp, nil
		}
	}
	return nil, domain.ErrEmployeeNotFound
}

func TestCreateEmployee(t *testing.T) {
	repo := NewMockEmployeeRepository()
	service := NewEmployeeService(repo)

	input := domain.EmployeeCreateInput{
		Name:            "John Doe",
		Email:           "john@example.com",
		Role:            "Developer",
		RoleDescription: "Full-stack developer working on web applications",
		MonthlyHours:    160,
	}

	employee, err := service.CreateEmployee(context.Background(), input)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if employee.Name != input.Name {
		t.Errorf("Expected name %s, got %s", input.Name, employee.Name)
	}

	if employee.Email != input.Email {
		t.Errorf("Expected email %s, got %s", input.Email, employee.Email)
	}

	if !employee.Active {
		t.Error("Expected employee to be active")
	}
}

func TestCreateEmployeeDuplicate(t *testing.T) {
	repo := NewMockEmployeeRepository()
	service := NewEmployeeService(repo)

	input := domain.EmployeeCreateInput{
		Name:            "Jane Doe",
		Email:           "jane@example.com",
		Role:            "Designer",
		RoleDescription: "UI/UX Designer",
		MonthlyHours:    160,
	}

	// Create first employee
	_, err := service.CreateEmployee(context.Background(), input)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Try to create duplicate
	_, err = service.CreateEmployee(context.Background(), input)
	if err != domain.ErrEmployeeAlreadyExists {
		t.Errorf("Expected ErrEmployeeAlreadyExists, got %v", err)
	}
}

func TestGetActiveEmployees(t *testing.T) {
	repo := NewMockEmployeeRepository()
	service := NewEmployeeService(repo)

	// Create active employee
	input1 := domain.EmployeeCreateInput{
		Name:            "Active Employee",
		Email:           "active@example.com",
		Role:            "Developer",
		RoleDescription: "Active developer",
		MonthlyHours:    160,
	}
	emp1, _ := service.CreateEmployee(context.Background(), input1)

	// Create and delete another employee
	input2 := domain.EmployeeCreateInput{
		Name:            "Inactive Employee",
		Email:           "inactive@example.com",
		Role:            "Developer",
		RoleDescription: "Inactive developer",
		MonthlyHours:    160,
	}
	emp2, _ := service.CreateEmployee(context.Background(), input2)
	service.DeleteEmployee(context.Background(), emp2.ID)

	// Get active employees
	active, err := service.GetActiveEmployees(context.Background())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(active) != 1 {
		t.Errorf("Expected 1 active employee, got %d", len(active))
	}

	if active[0].ID != emp1.ID {
		t.Errorf("Expected active employee ID %s, got %s", emp1.ID, active[0].ID)
	}
}
