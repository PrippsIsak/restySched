package domain

import (
	"regexp"
	"strings"
	"time"
)

// Employee represents an employee in the system
type Employee struct {
	ID               string         `json:"id" bson:"id"`
	Name             string         `json:"name" bson:"name"`
	Email            string         `json:"email" bson:"email"`
	Role             string         `json:"role" bson:"role"`
	RoleDescription  string         `json:"role_description" bson:"role_description"`
	MonthlyHours     int            `json:"monthly_hours" bson:"monthly_hours"`
	Active           bool           `json:"active" bson:"active"`
	Availability     []Availability `json:"availability,omitempty" bson:"availability,omitempty"`
	CreatedAt        time.Time      `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at" bson:"updated_at"`
}

// Availability represents an employee's availability for a date range
type Availability struct {
	StartDate   time.Time `json:"start_date" bson:"start_date"`
	EndDate     time.Time `json:"end_date" bson:"end_date"`
	Type        string    `json:"type" bson:"type"` // available, unavailable, preferred
	Reason      string    `json:"reason,omitempty" bson:"reason,omitempty"`
	ShiftTypes  []string  `json:"shift_types,omitempty" bson:"shift_types,omitempty"` // If empty, applies to all shift types
}

// AvailabilityType constants
const (
	AvailabilityTypeAvailable   = "available"
	AvailabilityTypeUnavailable = "unavailable"
	AvailabilityTypePreferred   = "preferred"
)

// EmployeeCreateInput represents the data needed to create a new employee
type EmployeeCreateInput struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Role            string `json:"role"`
	RoleDescription string `json:"role_description"`
	MonthlyHours    int    `json:"monthly_hours"`
}

// Email validation regex pattern
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// Validate checks if the employee data is valid
func (e *Employee) Validate() error {
	// Validate name
	if strings.TrimSpace(e.Name) == "" {
		return ErrInvalidEmployeeName
	}
	if len(e.Name) > 100 {
		return ErrInvalidEmployeeName
	}

	// Validate email
	if strings.TrimSpace(e.Email) == "" {
		return ErrInvalidEmployeeEmail
	}
	if !emailRegex.MatchString(e.Email) {
		return ErrInvalidEmployeeEmail
	}
	if len(e.Email) > 255 {
		return ErrInvalidEmployeeEmail
	}

	// Validate role
	if strings.TrimSpace(e.Role) == "" {
		return ErrInvalidEmployeeRole
	}
	if len(e.Role) > 100 {
		return ErrInvalidEmployeeRole
	}

	// Validate monthly hours
	if e.MonthlyHours <= 0 {
		return ErrInvalidMonthlyHours
	}
	if e.MonthlyHours > 744 { // Max hours in a month (31 days * 24 hours)
		return ErrInvalidMonthlyHours
	}

	return nil
}

// SanitizeEmployeeInput sanitizes and trims input data
func SanitizeEmployeeInput(input *EmployeeCreateInput) {
	input.Name = strings.TrimSpace(input.Name)
	input.Email = strings.TrimSpace(strings.ToLower(input.Email))
	input.Role = strings.TrimSpace(input.Role)
	input.RoleDescription = strings.TrimSpace(input.RoleDescription)
}

// IsAvailableOn checks if the employee is available on a specific date
func (e *Employee) IsAvailableOn(date time.Time, shiftType string) bool {
	// If no availability constraints, employee is available
	if len(e.Availability) == 0 {
		return true
	}

	// Check each availability entry
	for _, avail := range e.Availability {
		// Check if date falls within this availability range
		if (date.After(avail.StartDate) || date.Equal(avail.StartDate)) &&
			(date.Before(avail.EndDate) || date.Equal(avail.EndDate)) {

			// If shift types are specified, check if current shift type matches
			if len(avail.ShiftTypes) > 0 {
				shiftMatch := false
				for _, st := range avail.ShiftTypes {
					if st == shiftType {
						shiftMatch = true
						break
					}
				}
				if !shiftMatch {
					continue // This availability doesn't apply to this shift type
				}
			}

			// If this is an unavailable period, employee is not available
			if avail.Type == AvailabilityTypeUnavailable {
				return false
			}
		}
	}

	// Default to available if no unavailable periods match
	return true
}

// GetPreference returns the preference level for a date (0 = no preference, 1 = preferred)
func (e *Employee) GetPreference(date time.Time, shiftType string) int {
	for _, avail := range e.Availability {
		if (date.After(avail.StartDate) || date.Equal(avail.StartDate)) &&
			(date.Before(avail.EndDate) || date.Equal(avail.EndDate)) {

			// Check shift types if specified
			if len(avail.ShiftTypes) > 0 {
				shiftMatch := false
				for _, st := range avail.ShiftTypes {
					if st == shiftType {
						shiftMatch = true
						break
					}
				}
				if !shiftMatch {
					continue
				}
			}

			if avail.Type == AvailabilityTypePreferred {
				return 1
			}
		}
	}
	return 0
}
