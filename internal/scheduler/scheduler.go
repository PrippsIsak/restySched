package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/isak/restySched/internal/service"
)

type Scheduler struct {
	scheduleService *service.ScheduleService
	scheduler       gocron.Scheduler
}

// NewScheduler creates a new biweekly scheduler
func NewScheduler(scheduleService *service.ScheduleService) (*Scheduler, error) {
	s, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}

	return &Scheduler{
		scheduleService: scheduleService,
		scheduler:       s,
	}, nil
}

// Start begins the automated schedule generation
func (s *Scheduler) Start() error {
	// Schedule to run every 2 weeks (14 days)
	_, err := s.scheduler.NewJob(
		gocron.DurationJob(14*24*time.Hour),
		gocron.NewTask(s.generateAndSendSchedule),
		gocron.WithName("biweekly-schedule-generation"),
	)
	if err != nil {
		return err
	}

	log.Println("Scheduler started - will generate schedules every 2 weeks")
	s.scheduler.Start()
	return nil
}

// Stop stops the scheduler
func (s *Scheduler) Stop() error {
	return s.scheduler.Shutdown()
}

// RunNow triggers immediate schedule generation (useful for testing)
func (s *Scheduler) RunNow() error {
	s.generateAndSendSchedule()
	return nil
}

func (s *Scheduler) generateAndSendSchedule() {
	ctx := context.Background()

	log.Println("Starting biweekly schedule generation...")

	schedule, err := s.scheduleService.GenerateBiweeklySchedule(ctx)
	if err != nil {
		log.Printf("ERROR: Failed to generate schedule: %v", err)
		return
	}

	log.Printf("Schedule generated successfully: %s", schedule.ID)

	// Send to n8n
	if err := s.scheduleService.SendScheduleToN8N(ctx, schedule.ID); err != nil {
		log.Printf("ERROR: Failed to send schedule to n8n: %v", err)
		return
	}

	log.Printf("Schedule sent to n8n successfully: %s", schedule.ID)
}
