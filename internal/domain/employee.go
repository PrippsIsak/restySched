package domain

import "time"

// Employee represents an employee in the system
type Employee struct {
	ID               string    `json:"id" bson:"id"`
	Name             string    `json:"name" bson:"name"`
	Email            string    `json:"email" bson:"email"`
	Role             string    `json:"role" bson:"role"`
	RoleDescription  string    `json:"role_description" bson:"role_description"`
	MonthlyHours     int       `json:"monthly_hours" bson:"monthly_hours"`
	Active           bool      `json:"active" bson:"active"`
	CreatedAt        time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" bson:"updated_at"`
}

// EmployeeCreateInput represents the data needed to create a new employee
type EmployeeCreateInput struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Role            string `json:"role"`
	RoleDescription string `json:"role_description"`
	MonthlyHours    int    `json:"monthly_hours"`
}

// Validate checks if the employee data is valid
func (e *Employee) Validate() error {
	if e.Name == "" {
		return ErrInvalidEmployeeName
	}
	if e.Email == "" {
		return ErrInvalidEmployeeEmail
	}
	if e.Role == "" {
		return ErrInvalidEmployeeRole
	}
	if e.MonthlyHours <= 0 {
		return ErrInvalidMonthlyHours
	}
	return nil
}
