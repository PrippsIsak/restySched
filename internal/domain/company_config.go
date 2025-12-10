package domain

import (
	"errors"
	"time"
)

var (
	ErrInvalidCompanyName        = errors.New("company name is required")
	ErrInvalidWorkingHours       = errors.New("invalid working hours configuration")
	ErrInvalidShiftRequirements  = errors.New("invalid shift requirements")
	ErrCompanyConfigNotFound     = errors.New("company configuration not found")
	ErrCompanyConfigAlreadyExists = errors.New("company configuration already exists")
)

// CompanyConfig represents the company's scheduling configuration
type CompanyConfig struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	CompanyName string    `json:"company_name" bson:"company_name"`

	// Business hours and days
	WorkingHours WorkingHours `json:"working_hours" bson:"working_hours"`

	// Shift requirements
	ShiftRequirements []ShiftRequirement `json:"shift_requirements" bson:"shift_requirements"`

	// Scheduling policies
	SchedulingPolicies SchedulingPolicies `json:"scheduling_policies" bson:"scheduling_policies"`

	// AI Context - additional instructions for n8n AI agent
	AIContext string `json:"ai_context" bson:"ai_context"`

	// Metadata
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// WorkingHours defines when the company operates
type WorkingHours struct {
	// Days of week the company operates (0 = Sunday, 6 = Saturday)
	WorkingDays []int `json:"working_days" bson:"working_days"`

	// Standard operating hours
	OpenTime  string `json:"open_time" bson:"open_time"`   // e.g., "09:00"
	CloseTime string `json:"close_time" bson:"close_time"` // e.g., "17:00"

	// Timezone
	Timezone string `json:"timezone" bson:"timezone"` // e.g., "Europe/Oslo"
}

// ShiftRequirement defines what shifts are needed and how many employees
type ShiftRequirement struct {
	ShiftType         string `json:"shift_type" bson:"shift_type"` // morning, afternoon, evening, full_day, night
	MinEmployees      int    `json:"min_employees" bson:"min_employees"`
	MaxEmployees      int    `json:"max_employees" bson:"max_employees"`
	RequiredSkills    []string `json:"required_skills,omitempty" bson:"required_skills,omitempty"`
	Description       string `json:"description" bson:"description"`
}

// SchedulingPolicies defines rules and constraints for scheduling
type SchedulingPolicies struct {
	// Maximum consecutive work days
	MaxConsecutiveDays int `json:"max_consecutive_days" bson:"max_consecutive_days"`

	// Minimum rest hours between shifts
	MinRestHours int `json:"min_rest_hours" bson:"min_rest_hours"`

	// Allow overtime
	AllowOvertime bool `json:"allow_overtime" bson:"allow_overtime"`

	// Maximum overtime hours per month
	MaxOvertimeHours int `json:"max_overtime_hours" bson:"max_overtime_hours"`

	// Require employee consent for weekend shifts
	WeekendConsentRequired bool `json:"weekend_consent_required" bson:"weekend_consent_required"`

	// Fair distribution of shifts
	FairDistribution bool `json:"fair_distribution" bson:"fair_distribution"`
}

// Validate checks if the company configuration is valid
func (c *CompanyConfig) Validate() error {
	if c.CompanyName == "" {
		return ErrInvalidCompanyName
	}

	// Validate working hours
	if c.WorkingHours.OpenTime == "" || c.WorkingHours.CloseTime == "" {
		return ErrInvalidWorkingHours
	}

	if len(c.WorkingHours.WorkingDays) == 0 {
		return ErrInvalidWorkingHours
	}

	// Validate shift requirements
	if len(c.ShiftRequirements) == 0 {
		return ErrInvalidShiftRequirements
	}

	for _, req := range c.ShiftRequirements {
		if req.MinEmployees < 0 || req.MaxEmployees < req.MinEmployees {
			return ErrInvalidShiftRequirements
		}
	}

	return nil
}

// GetContextForAI returns a formatted context string for the AI agent
func (c *CompanyConfig) GetContextForAI() string {
	context := "Company: " + c.CompanyName + "\n\n"

	context += "Working Hours:\n"
	context += "- Days: "
	dayNames := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	for i, day := range c.WorkingHours.WorkingDays {
		if i > 0 {
			context += ", "
		}
		context += dayNames[day]
	}
	context += "\n"
	context += "- Hours: " + c.WorkingHours.OpenTime + " - " + c.WorkingHours.CloseTime + "\n"
	context += "- Timezone: " + c.WorkingHours.Timezone + "\n\n"

	context += "Shift Requirements:\n"
	for _, req := range c.ShiftRequirements {
		context += "- " + req.ShiftType + ": "
		context += req.Description
		if req.MinEmployees > 0 {
			context += " (Min: " + string(rune(req.MinEmployees+'0')) + " employees)"
		}
		context += "\n"
	}
	context += "\n"

	context += "Scheduling Policies:\n"
	context += "- Maximum consecutive work days: " + string(rune(c.SchedulingPolicies.MaxConsecutiveDays+'0')) + "\n"
	context += "- Minimum rest hours between shifts: " + string(rune(c.SchedulingPolicies.MinRestHours+'0')) + "h\n"
	if c.SchedulingPolicies.FairDistribution {
		context += "- Fair distribution of shifts across all employees\n"
	}
	if c.SchedulingPolicies.WeekendConsentRequired {
		context += "- Weekend shifts require employee consent\n"
	}
	context += "\n"

	if c.AIContext != "" {
		context += "Additional Instructions:\n" + c.AIContext + "\n"
	}

	return context
}
