package domain

import "errors"

// Domain errors
var (
	// Employee errors
	ErrEmployeeNotFound      = errors.New("employee not found")
	ErrInvalidEmployeeName   = errors.New("invalid employee name")
	ErrInvalidEmployeeEmail  = errors.New("invalid employee email")
	ErrInvalidEmployeeRole   = errors.New("invalid employee role")
	ErrInvalidMonthlyHours   = errors.New("monthly hours must be greater than 0")
	ErrEmployeeAlreadyExists = errors.New("employee already exists")

	// Schedule errors
	ErrScheduleNotFound     = errors.New("schedule not found")
	ErrInvalidSchedulePeriod = errors.New("invalid schedule period")
	ErrScheduleAlreadySent  = errors.New("schedule already sent to n8n")

	// General errors
	ErrInternalServer = errors.New("internal server error")
)
