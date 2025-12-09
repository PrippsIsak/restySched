package domain

import (
	"testing"
)

func TestEmployee_Validate(t *testing.T) {
	tests := []struct {
		name     string
		employee Employee
		wantErr  error
	}{
		{
			name: "valid employee",
			employee: Employee{
				Name:         "John Doe",
				Email:        "john@example.com",
				Role:         "Developer",
				MonthlyHours: 160,
			},
			wantErr: nil,
		},
		{
			name: "empty name",
			employee: Employee{
				Name:         "",
				Email:        "john@example.com",
				Role:         "Developer",
				MonthlyHours: 160,
			},
			wantErr: ErrInvalidEmployeeName,
		},
		{
			name: "whitespace only name",
			employee: Employee{
				Name:         "   ",
				Email:        "john@example.com",
				Role:         "Developer",
				MonthlyHours: 160,
			},
			wantErr: ErrInvalidEmployeeName,
		},
		{
			name: "name too long",
			employee: Employee{
				Name:         string(make([]byte, 101)), // 101 characters
				Email:        "john@example.com",
				Role:         "Developer",
				MonthlyHours: 160,
			},
			wantErr: ErrInvalidEmployeeName,
		},
		{
			name: "empty email",
			employee: Employee{
				Name:         "John Doe",
				Email:        "",
				Role:         "Developer",
				MonthlyHours: 160,
			},
			wantErr: ErrInvalidEmployeeEmail,
		},
		{
			name: "invalid email format",
			employee: Employee{
				Name:         "John Doe",
				Email:        "not-an-email",
				Role:         "Developer",
				MonthlyHours: 160,
			},
			wantErr: ErrInvalidEmployeeEmail,
		},
		{
			name: "invalid email missing @",
			employee: Employee{
				Name:         "John Doe",
				Email:        "johnexample.com",
				Role:         "Developer",
				MonthlyHours: 160,
			},
			wantErr: ErrInvalidEmployeeEmail,
		},
		{
			name: "invalid email missing domain",
			employee: Employee{
				Name:         "John Doe",
				Email:        "john@",
				Role:         "Developer",
				MonthlyHours: 160,
			},
			wantErr: ErrInvalidEmployeeEmail,
		},
		{
			name: "empty role",
			employee: Employee{
				Name:         "John Doe",
				Email:        "john@example.com",
				Role:         "",
				MonthlyHours: 160,
			},
			wantErr: ErrInvalidEmployeeRole,
		},
		{
			name: "zero monthly hours",
			employee: Employee{
				Name:         "John Doe",
				Email:        "john@example.com",
				Role:         "Developer",
				MonthlyHours: 0,
			},
			wantErr: ErrInvalidMonthlyHours,
		},
		{
			name: "negative monthly hours",
			employee: Employee{
				Name:         "John Doe",
				Email:        "john@example.com",
				Role:         "Developer",
				MonthlyHours: -10,
			},
			wantErr: ErrInvalidMonthlyHours,
		},
		{
			name: "monthly hours too high",
			employee: Employee{
				Name:         "John Doe",
				Email:        "john@example.com",
				Role:         "Developer",
				MonthlyHours: 800,
			},
			wantErr: ErrInvalidMonthlyHours,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.employee.Validate()
			if err != tt.wantErr {
				t.Errorf("Employee.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSanitizeEmployeeInput(t *testing.T) {
	tests := []struct {
		name     string
		input    EmployeeCreateInput
		expected EmployeeCreateInput
	}{
		{
			name: "trim whitespace",
			input: EmployeeCreateInput{
				Name:            "  John Doe  ",
				Email:           "  JOHN@EXAMPLE.COM  ",
				Role:            "  Developer  ",
				RoleDescription: "  Writes code  ",
				MonthlyHours:    160,
			},
			expected: EmployeeCreateInput{
				Name:            "John Doe",
				Email:           "john@example.com",
				Role:            "Developer",
				RoleDescription: "Writes code",
				MonthlyHours:    160,
			},
		},
		{
			name: "lowercase email",
			input: EmployeeCreateInput{
				Name:         "John Doe",
				Email:        "John.Doe@EXAMPLE.COM",
				Role:         "Developer",
				MonthlyHours: 160,
			},
			expected: EmployeeCreateInput{
				Name:         "John Doe",
				Email:        "john.doe@example.com",
				Role:         "Developer",
				MonthlyHours: 160,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := tt.input
			SanitizeEmployeeInput(&input)

			if input.Name != tt.expected.Name {
				t.Errorf("Name = %v, want %v", input.Name, tt.expected.Name)
			}
			if input.Email != tt.expected.Email {
				t.Errorf("Email = %v, want %v", input.Email, tt.expected.Email)
			}
			if input.Role != tt.expected.Role {
				t.Errorf("Role = %v, want %v", input.Role, tt.expected.Role)
			}
			if input.RoleDescription != tt.expected.RoleDescription {
				t.Errorf("RoleDescription = %v, want %v", input.RoleDescription, tt.expected.RoleDescription)
			}
		})
	}
}
