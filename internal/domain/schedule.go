package domain

import "time"

// Schedule represents a generated schedule for a period
type Schedule struct {
	ID          string       `json:"id" bson:"id"`
	PeriodStart time.Time    `json:"period_start" bson:"period_start"`
	PeriodEnd   time.Time    `json:"period_end" bson:"period_end"`
	Employees   []Employee   `json:"employees" bson:"employees"`
	Status      string       `json:"status" bson:"status"` // draft, sent, completed
	SentToN8N   bool         `json:"sent_to_n8n" bson:"sent_to_n8n"`
	SentAt      *time.Time   `json:"sent_at,omitempty" bson:"sent_at,omitempty"`
	CreatedAt   time.Time    `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at" bson:"updated_at"`
}

// ScheduleStatus constants
const (
	ScheduleStatusDraft     = "draft"
	ScheduleStatusSent      = "sent"
	ScheduleStatusCompleted = "completed"
)

// N8NSchedulePayload represents the data sent to n8n webhook
type N8NSchedulePayload struct {
	ScheduleID  string              `json:"schedule_id"`
	PeriodStart string              `json:"period_start"`
	PeriodEnd   string              `json:"period_end"`
	Employees   []N8NEmployeeData   `json:"employees"`
	GeneratedAt string              `json:"generated_at"`
}

// N8NEmployeeData represents employee data for n8n
type N8NEmployeeData struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	Role            string `json:"role"`
	RoleDescription string `json:"role_description"`
	MonthlyHours    int    `json:"monthly_hours"`
}
