package service

import (
	"testing"
	"time"

	"github.com/isak/restySched/internal/domain"
)

func TestShiftGenerator_GenerateShifts(t *testing.T) {
	generator := NewShiftGenerator()

	// Create test employees
	employees := []domain.Employee{
		{
			ID:           "emp1",
			Name:         "John Doe",
			MonthlyHours: 160, // 20 days * 8 hours
		},
		{
			ID:           "emp2",
			Name:         "Jane Smith",
			MonthlyHours: 120, // 15 days * 8 hours
		},
		{
			ID:           "emp3",
			Name:         "Bob Johnson",
			MonthlyHours: 80, // 10 days * 8 hours
		},
	}

	// Generate shifts for 2 weeks (10 workdays)
	start := time.Date(2025, 1, 6, 0, 0, 0, 0, time.UTC)  // Monday
	end := time.Date(2025, 1, 17, 0, 0, 0, 0, time.UTC)   // Friday (2 weeks later)

	assignments := generator.GenerateShifts(employees, start, end)

	// Verify assignments were created
	if len(assignments) == 0 {
		t.Fatal("No assignments were generated")
	}

	// Count assignments per employee
	empAssignments := make(map[string]int)
	empHours := make(map[string]float64)

	for _, assignment := range assignments {
		empAssignments[assignment.EmployeeID]++
		empHours[assignment.EmployeeID] += assignment.Hours

		// Verify assignment has required fields
		if assignment.EmployeeID == "" {
			t.Error("Assignment missing EmployeeID")
		}
		if assignment.EmployeeName == "" {
			t.Error("Assignment missing EmployeeName")
		}
		if assignment.ShiftType == "" {
			t.Error("Assignment missing ShiftType")
		}
		if assignment.Hours <= 0 {
			t.Error("Assignment has invalid hours")
		}
		if assignment.StartTime == "" || assignment.EndTime == "" {
			t.Error("Assignment missing time information")
		}
	}

	t.Logf("Generated %d total assignments", len(assignments))
	for empID, count := range empAssignments {
		hours := empHours[empID]
		t.Logf("Employee %s: %d shifts, %.1f hours", empID, count, hours)
	}

	// Each employee should have at least some assignments
	for _, emp := range employees {
		if empAssignments[emp.ID] == 0 {
			t.Errorf("Employee %s has no assignments", emp.ID)
		}
	}
}

func TestShiftGenerator_CountWorkdays(t *testing.T) {
	generator := NewShiftGenerator()

	tests := []struct {
		name     string
		start    time.Time
		end      time.Time
		expected int
	}{
		{
			name:     "one week monday to friday",
			start:    time.Date(2025, 1, 6, 0, 0, 0, 0, time.UTC),  // Monday
			end:      time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC), // Friday
			expected: 5,
		},
		{
			name:     "two weeks",
			start:    time.Date(2025, 1, 6, 0, 0, 0, 0, time.UTC),  // Monday
			end:      time.Date(2025, 1, 17, 0, 0, 0, 0, time.UTC), // Friday
			expected: 10,
		},
		{
			name:     "including weekend",
			start:    time.Date(2025, 1, 6, 0, 0, 0, 0, time.UTC),  // Monday
			end:      time.Date(2025, 1, 12, 0, 0, 0, 0, time.UTC), // Sunday
			expected: 5, // Should only count Mon-Fri
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generator.countWorkdays(tt.start, tt.end)
			if result != tt.expected {
				t.Errorf("countWorkdays() = %d, want %d", result, tt.expected)
			}
		})
	}
}

func TestShiftGenerator_GetEmployeeStats(t *testing.T) {
	generator := NewShiftGenerator()

	assignments := []domain.ShiftAssignment{
		{
			EmployeeID:   "emp1",
			EmployeeName: "John Doe",
			Date:         time.Now(),
			ShiftType:    domain.ShiftTypeFullDay,
			Hours:        8.0,
		},
		{
			EmployeeID:   "emp1",
			EmployeeName: "John Doe",
			Date:         time.Now().AddDate(0, 0, 1),
			ShiftType:    domain.ShiftTypeFullDay,
			Hours:        8.0,
		},
		{
			EmployeeID:   "emp1",
			EmployeeName: "John Doe",
			Date:         time.Now().AddDate(0, 0, 2),
			ShiftType:    domain.ShiftTypeMorning,
			Hours:        4.0,
		},
	}

	stats := generator.GetEmployeeStats("emp1", assignments)

	if stats.TotalShifts != 3 {
		t.Errorf("TotalShifts = %d, want 3", stats.TotalShifts)
	}

	if stats.TotalHours != 20.0 {
		t.Errorf("TotalHours = %.1f, want 20.0", stats.TotalHours)
	}

	if stats.ShiftTypes[domain.ShiftTypeFullDay] != 2 {
		t.Errorf("Full day shifts = %d, want 2", stats.ShiftTypes[domain.ShiftTypeFullDay])
	}

	if stats.ShiftTypes[domain.ShiftTypeMorning] != 1 {
		t.Errorf("Morning shifts = %d, want 1", stats.ShiftTypes[domain.ShiftTypeMorning])
	}
}

func TestGetShiftDefinition(t *testing.T) {
	tests := []struct {
		shiftType string
		wantNil   bool
	}{
		{domain.ShiftTypeMorning, false},
		{domain.ShiftTypeAfternoon, false},
		{domain.ShiftTypeEvening, false},
		{domain.ShiftTypeFullDay, false},
		{domain.ShiftTypeNight, false},
		{"invalid_type", true},
	}

	for _, tt := range tests {
		t.Run(tt.shiftType, func(t *testing.T) {
			result := domain.GetShiftDefinition(tt.shiftType)
			if tt.wantNil {
				if result != nil {
					t.Errorf("GetShiftDefinition(%s) should return nil", tt.shiftType)
				}
			} else {
				if result == nil {
					t.Errorf("GetShiftDefinition(%s) should not return nil", tt.shiftType)
				}
				if result.Type != tt.shiftType {
					t.Errorf("ShiftType = %s, want %s", result.Type, tt.shiftType)
				}
				if result.Hours <= 0 {
					t.Error("Hours should be > 0")
				}
			}
		})
	}
}
