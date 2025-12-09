package service

import (
	"math"
	"time"

	"github.com/isak/restySched/internal/domain"
	"github.com/rs/zerolog/log"
)

// ShiftGenerator handles the logic for generating shift assignments
type ShiftGenerator struct{}

// NewShiftGenerator creates a new shift generator
func NewShiftGenerator() *ShiftGenerator {
	return &ShiftGenerator{}
}

// GenerateShifts creates shift assignments for employees over the schedule period
func (g *ShiftGenerator) GenerateShifts(employees []domain.Employee, periodStart, periodEnd time.Time) []domain.ShiftAssignment {
	var assignments []domain.ShiftAssignment

	// Calculate total days in period (excluding weekends for now)
	totalDays := g.countWorkdays(periodStart, periodEnd)
	if totalDays == 0 {
		log.Warn().Msg("No workdays in schedule period")
		return assignments
	}

	log.Debug().
		Int("total_days", totalDays).
		Int("employees", len(employees)).
		Msg("Generating shifts")

	// Calculate how many hours each employee should work during this period
	employeeTargets := g.calculateEmployeeTargets(employees, periodStart, periodEnd)

	// Track hours assigned to each employee
	assignedHours := make(map[string]float64)
	for _, emp := range employees {
		assignedHours[emp.ID] = 0
	}

	// Generate shifts day by day
	currentDate := periodStart
	for currentDate.Before(periodEnd) || currentDate.Equal(periodEnd) {
		// Skip weekends (Saturday = 6, Sunday = 0)
		if currentDate.Weekday() == time.Saturday || currentDate.Weekday() == time.Sunday {
			currentDate = currentDate.AddDate(0, 0, 1)
			continue
		}

		// Assign shifts for this day
		dayShifts := g.assignDayShifts(employees, currentDate, employeeTargets, assignedHours)
		assignments = append(assignments, dayShifts...)

		currentDate = currentDate.AddDate(0, 0, 1)
	}

	log.Info().
		Int("total_assignments", len(assignments)).
		Msg("Shift generation complete")

	return assignments
}

// countWorkdays counts the number of weekdays in the period
func (g *ShiftGenerator) countWorkdays(start, end time.Time) int {
	count := 0
	current := start
	for current.Before(end) || current.Equal(end) {
		if current.Weekday() != time.Saturday && current.Weekday() != time.Sunday {
			count++
		}
		current = current.AddDate(0, 0, 1)
	}
	return count
}

// calculateEmployeeTargets calculates target hours for each employee in the period
func (g *ShiftGenerator) calculateEmployeeTargets(employees []domain.Employee, start, end time.Time) map[string]float64 {
	targets := make(map[string]float64)

	// Calculate how many months are in the period
	days := end.Sub(start).Hours() / 24
	monthFraction := days / 30.0 // Approximate

	for _, emp := range employees {
		// Target hours = (monthly hours * fraction of month in this period)
		targets[emp.ID] = float64(emp.MonthlyHours) * monthFraction
	}

	return targets
}

// assignDayShifts assigns shifts to employees for a single day
func (g *ShiftGenerator) assignDayShifts(
	employees []domain.Employee,
	date time.Time,
	targets map[string]float64,
	assignedHours map[string]float64,
) []domain.ShiftAssignment {
	var assignments []domain.ShiftAssignment

	// Strategy: Assign full-day shifts to employees who need the most hours
	// This is a simple round-robin approach - can be enhanced with more sophisticated algorithms

	// Find employees who still need hours
	type employeeNeed struct {
		employee      domain.Employee
		hoursNeeded   float64
		percentNeeded float64
	}

	var needsList []employeeNeed
	shiftType := domain.ShiftTypeFullDay // We'll use full-day shifts

	for _, emp := range employees {
		// Check if employee is available on this date
		if !emp.IsAvailableOn(date, shiftType) {
			log.Debug().
				Str("employee", emp.Name).
				Time("date", date).
				Msg("Employee unavailable, skipping")
			continue
		}

		target := targets[emp.ID]
		assigned := assignedHours[emp.ID]
		needed := target - assigned

		if needed > 0 {
			percentNeeded := (needed / target) * 100

			// Add preference bonus to prioritize preferred shifts
			preferenceBonus := float64(emp.GetPreference(date, shiftType)) * 10.0

			needsList = append(needsList, employeeNeed{
				employee:      emp,
				hoursNeeded:   needed,
				percentNeeded: percentNeeded + preferenceBonus,
			})
		}
	}

	// Sort by percent needed + preference (highest first)
	for i := 0; i < len(needsList); i++ {
		for j := i + 1; j < len(needsList); j++ {
			if needsList[j].percentNeeded > needsList[i].percentNeeded {
				needsList[i], needsList[j] = needsList[j], needsList[i]
			}
		}
	}

	// Assign shifts based on need and availability
	// For simplicity, we'll assign 2-3 people per day with full-day shifts
	shiftsToAssign := int(math.Min(float64(len(needsList)), 3))

	for i := 0; i < shiftsToAssign; i++ {
		emp := needsList[i].employee

		// Assign a full-day shift (8 hours)
		shiftDef := domain.GetShiftDefinition(shiftType)
		if shiftDef == nil {
			continue
		}

		assignment := domain.ShiftAssignment{
			EmployeeID:   emp.ID,
			EmployeeName: emp.Name,
			Date:         date,
			ShiftType:    shiftDef.Type,
			StartTime:    shiftDef.StartTime,
			EndTime:      shiftDef.EndTime,
			Hours:        shiftDef.Hours,
		}

		assignments = append(assignments, assignment)
		assignedHours[emp.ID] += shiftDef.Hours
	}

	return assignments
}

// GetEmployeeStats returns statistics about shift assignments for an employee
func (g *ShiftGenerator) GetEmployeeStats(employeeID string, assignments []domain.ShiftAssignment) EmployeeShiftStats {
	stats := EmployeeShiftStats{
		EmployeeID:  employeeID,
		TotalHours:  0,
		TotalShifts: 0,
		ShiftTypes:  make(map[string]int),
	}

	for _, assignment := range assignments {
		if assignment.EmployeeID == employeeID {
			stats.TotalHours += assignment.Hours
			stats.TotalShifts++
			stats.ShiftTypes[assignment.ShiftType]++
		}
	}

	return stats
}

// EmployeeShiftStats holds statistics about an employee's shifts
type EmployeeShiftStats struct {
	EmployeeID  string
	TotalHours  float64
	TotalShifts int
	ShiftTypes  map[string]int
}
