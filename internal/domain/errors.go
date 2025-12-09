package domain

import "errors"

// Domain errors
var (
	// Employee errors
	ErrEmployeeNotFound      = errors.New("employee not found")
	ErrInvalidEmployeeName   = errors.New("employee name is required and must be less than 100 characters")
	ErrInvalidEmployeeEmail  = errors.New("valid employee email is required (max 255 characters)")
	ErrInvalidEmployeeRole   = errors.New("employee role is required and must be less than 100 characters")
	ErrInvalidMonthlyHours   = errors.New("monthly hours must be between 1 and 744")
	ErrEmployeeAlreadyExists = errors.New("an employee with this email already exists")

	// Schedule errors
	ErrScheduleNotFound      = errors.New("schedule not found")
	ErrInvalidSchedulePeriod = errors.New("schedule period end must be after period start")
	ErrScheduleAlreadySent   = errors.New("schedule has already been sent to n8n")

	// General errors
	ErrInternalServer = errors.New("internal server error")
)
