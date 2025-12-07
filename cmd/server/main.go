package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/isak/restySched/internal/config"
	"github.com/isak/restySched/internal/handler"
	"github.com/isak/restySched/internal/n8n"
	"github.com/isak/restySched/internal/repository/mongodb"
	"github.com/isak/restySched/internal/scheduler"
	"github.com/isak/restySched/internal/service"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize MongoDB database
	db, err := mongodb.InitDB(cfg.MongoURI, cfg.MongoDatabase)
	if err != nil {
		log.Fatalf("Failed to initialize MongoDB: %v", err)
	}

	log.Printf("Connected to MongoDB database: %s", cfg.MongoDatabase)

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

	// Setup routes
	mux := http.NewServeMux()

	// Home
	mux.HandleFunc("GET /", homeHandler.Home)

	// Employee routes
	mux.HandleFunc("GET /employees", employeeHandler.ListEmployees)
	mux.HandleFunc("GET /employees/new", employeeHandler.ShowNewForm)
	mux.HandleFunc("GET /employees/{id}/edit", employeeHandler.ShowEditForm)
	mux.HandleFunc("POST /employees", employeeHandler.CreateEmployee)
	mux.HandleFunc("PUT /employees/{id}", employeeHandler.UpdateEmployee)
	mux.HandleFunc("DELETE /employees/{id}", employeeHandler.DeleteEmployee)

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
			log.Fatalf("Failed to create scheduler: %v", err)
		}

		if err := sched.Start(); err != nil {
			log.Fatalf("Failed to start scheduler: %v", err)
		}
		defer sched.Stop()
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
		log.Printf("Server starting on port %s", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server shutting down...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
