package domain

import "time"

// Schedule represents a generated schedule for a period
type Schedule struct {
	ID          string              `json:"id" bson:"id"`
	PeriodStart time.Time           `json:"period_start" bson:"period_start"`
	PeriodEnd   time.Time           `json:"period_end" bson:"period_end"`
	Employees   []Employee          `json:"employees" bson:"employees"`
	Assignments []ShiftAssignment   `json:"assignments" bson:"assignments"`
	Status      string              `json:"status" bson:"status"` // draft, sent, completed
	SentToN8N   bool                `json:"sent_to_n8n" bson:"sent_to_n8n"`
	SentAt      *time.Time          `json:"sent_at,omitempty" bson:"sent_at,omitempty"`
	CreatedAt   time.Time           `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at" bson:"updated_at"`
}

// ShiftAssignment represents an employee's shift on a specific day
type ShiftAssignment struct {
	EmployeeID   string    `json:"employee_id" bson:"employee_id"`
	EmployeeName string    `json:"employee_name" bson:"employee_name"`
	Date         time.Time `json:"date" bson:"date"`
	ShiftType    string    `json:"shift_type" bson:"shift_type"` // morning, afternoon, evening, full_day
	StartTime    string    `json:"start_time" bson:"start_time"` // e.g., "09:00"
	EndTime      string    `json:"end_time" bson:"end_time"`     // e.g., "17:00"
	Hours        float64   `json:"hours" bson:"hours"`           // Duration in hours
}

// ScheduleStatus constants
const (
	ScheduleStatusDraft     = "draft"
	ScheduleStatusSent      = "sent"
	ScheduleStatusCompleted = "completed"
)

// ShiftType constants
const (
	ShiftTypeMorning    = "morning"      // 09:00 - 13:00
	ShiftTypeAfternoon  = "afternoon"    // 13:00 - 17:00
	ShiftTypeEvening    = "evening"      // 17:00 - 21:00
	ShiftTypeFullDay    = "full_day"     // 09:00 - 17:00
	ShiftTypeNight      = "night"        // 21:00 - 05:00
)

// ShiftDefinition defines the time range for each shift type
type ShiftDefinition struct {
	Type      string
	StartTime string
	EndTime   string
	Hours     float64
}

// GetShiftDefinitions returns all available shift types
func GetShiftDefinitions() []ShiftDefinition {
	return []ShiftDefinition{
		{Type: ShiftTypeMorning, StartTime: "09:00", EndTime: "13:00", Hours: 4.0},
		{Type: ShiftTypeAfternoon, StartTime: "13:00", EndTime: "17:00", Hours: 4.0},
		{Type: ShiftTypeEvening, StartTime: "17:00", EndTime: "21:00", Hours: 4.0},
		{Type: ShiftTypeFullDay, StartTime: "09:00", EndTime: "17:00", Hours: 8.0},
		{Type: ShiftTypeNight, StartTime: "21:00", EndTime: "05:00", Hours: 8.0},
	}
}

// GetShiftDefinition returns the shift definition for a given type
func GetShiftDefinition(shiftType string) *ShiftDefinition {
	for _, def := range GetShiftDefinitions() {
		if def.Type == shiftType {
			return &def
		}
	}
	return nil
}

// N8NSchedulePayload represents the data sent to n8n webhook
type N8NSchedulePayload struct {
	ScheduleID  string              `json:"schedule_id"`
	PeriodStart string              `json:"period_start"`
	PeriodEnd   string              `json:"period_end"`
	Employees   []N8NEmployeeData   `json:"employees"`
	Assignments []ShiftAssignment   `json:"assignments"`
	TotalShifts int                 `json:"total_shifts"`
	TotalHours  float64             `json:"total_hours"`
	GeneratedAt string              `json:"generated_at"`
}

// N8NEmployeeData represents employee data for n8n
type N8NEmployeeData struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	Email           string  `json:"email"`
	Role            string  `json:"role"`
	RoleDescription string  `json:"role_description"`
	MonthlyHours    int     `json:"monthly_hours"`
	AssignedHours   float64 `json:"assigned_hours"`
	AssignedShifts  int     `json:"assigned_shifts"`
}
