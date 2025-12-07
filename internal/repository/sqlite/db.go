package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// InitDB initializes the SQLite database and creates tables
func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	employeesTable := `
	CREATE TABLE IF NOT EXISTS employees (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		role TEXT NOT NULL,
		role_description TEXT NOT NULL,
		monthly_hours INTEGER NOT NULL,
		active BOOLEAN NOT NULL DEFAULT 1,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_employees_email ON employees(email);
	CREATE INDEX IF NOT EXISTS idx_employees_active ON employees(active);
	`

	schedulesTable := `
	CREATE TABLE IF NOT EXISTS schedules (
		id TEXT PRIMARY KEY,
		period_start DATETIME NOT NULL,
		period_end DATETIME NOT NULL,
		employees TEXT NOT NULL,
		status TEXT NOT NULL,
		sent_to_n8n BOOLEAN NOT NULL DEFAULT 0,
		sent_at DATETIME,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_schedules_period ON schedules(period_start, period_end);
	CREATE INDEX IF NOT EXISTS idx_schedules_status ON schedules(status);
	`

	if _, err := db.Exec(employeesTable); err != nil {
		return fmt.Errorf("failed to create employees table: %w", err)
	}

	if _, err := db.Exec(schedulesTable); err != nil {
		return fmt.Errorf("failed to create schedules table: %w", err)
	}

	return nil
}
