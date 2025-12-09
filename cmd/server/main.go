package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/isak/restySched/internal/config"
	"github.com/isak/restySched/internal/handler"
	"github.com/isak/restySched/internal/logger"
	"github.com/isak/restySched/internal/n8n"
	"github.com/isak/restySched/internal/repository/mongodb"
	"github.com/isak/restySched/internal/scheduler"
	"github.com/isak/restySched/internal/service"
	"github.com/rs/zerolog/log"
)

func main() {
	// Initialize logger
	logger.Init(false) // Set to true for debug mode

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Initialize MongoDB database
	db, err := mongodb.InitDB(cfg.MongoURI, cfg.MongoDatabase)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize MongoDB")
	}

	log.Info().Str("database", cfg.MongoDatabase).Msg("Connected to MongoDB")

	// Initialize repositories
	employeeRepo := mongodb.NewEmployeeRepository(db)
	scheduleRepo := mongodb.NewScheduleRepository(db)

	// Initialize n8n client
	n8nClient := n8n.NewClient(cfg.N8NWebhookURL)

	// Initialize services
	employeeService := service.NewEmployeeService(employeeRepo)
	scheduleService := service.NewScheduleService(scheduleRepo, employeeRepo, n8nClient)

	// Initialize handlers
	homeHandler := handler.NewHomeHandler()
	employeeHandler := handler.NewEmployeeHandler(employeeService)
	scheduleHandler := handler.NewScheduleHandler(scheduleService)
	healthHandler := handler.NewHealthHandler(employeeRepo)

	// Setup routes
	mux := http.NewServeMux()

	// Health check routes
	mux.HandleFunc("GET /health", healthHandler.Health)
	mux.HandleFunc("GET /health/ready", healthHandler.Ready)

	// Home
	mux.HandleFunc("GET /", homeHandler.Home)

	// Employee routes
	mux.HandleFunc("GET /employees", employeeHandler.ListEmployees)
	mux.HandleFunc("GET /employees/new", employeeHandler.ShowNewForm)
	mux.HandleFunc("GET /employees/{id}/edit", employeeHandler.ShowEditForm)
	mux.HandleFunc("POST /employees", employeeHandler.CreateEmployee)
	mux.HandleFunc("PUT /employees/{id}", employeeHandler.UpdateEmployee)
	mux.HandleFunc("DELETE /employees/{id}", employeeHandler.DeleteEmployee)

	// Employee availability routes
	mux.HandleFunc("GET /employees/{id}/availability", employeeHandler.ShowAvailabilityManager)
	mux.HandleFunc("POST /employees/{id}/availability", employeeHandler.AddAvailability)
	mux.HandleFunc("DELETE /employees/{id}/availability/{index}", employeeHandler.DeleteAvailability)

	// Schedule routes
	mux.HandleFunc("GET /schedules", scheduleHandler.ListSchedules)
	mux.HandleFunc("POST /schedules/generate", scheduleHandler.GenerateBiweeklySchedule)
	mux.HandleFunc("POST /schedules/{id}/send", scheduleHandler.SendToN8N)
	mux.HandleFunc("DELETE /schedules/{id}", scheduleHandler.DeleteSchedule)

	// Initialize and start scheduler if enabled
	var sched *scheduler.Scheduler
	if cfg.EnableScheduler {
		sched, err = scheduler.NewScheduler(scheduleService)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to create scheduler")
		}

		if err := sched.Start(); err != nil {
			log.Fatal().Err(err).Msg("Failed to start scheduler")
		}
		defer sched.Stop()
		log.Info().Msg("Automated scheduler started")
	}

	// Setup HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Info().
			Str("port", cfg.ServerPort).
			Str("address", "http://localhost:"+cfg.ServerPort).
			Msg("Server starting")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Server failed")
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Server shutting down...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exited gracefully")
}
