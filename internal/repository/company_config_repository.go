package repository

import (
	"context"

	"github.com/isak/restySched/internal/domain"
)

// CompanyConfigRepository defines the interface for company configuration operations
type CompanyConfigRepository interface {
	// Get retrieves the company configuration
	Get(ctx context.Context) (*domain.CompanyConfig, error)

	// Create creates a new company configuration
	Create(ctx context.Context, config *domain.CompanyConfig) error

	// Update updates the company configuration
	Update(ctx context.Context, config *domain.CompanyConfig) error

	// GetOrCreate retrieves the config or creates a default one
	GetOrCreate(ctx context.Context) (*domain.CompanyConfig, error)
}
