package mongodb

import (
	"context"
	"time"

	"github.com/isak/restySched/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CompanyConfigRepository struct {
	collection *mongo.Collection
}

func NewCompanyConfigRepository(db *mongo.Database) *CompanyConfigRepository {
	return &CompanyConfigRepository{
		collection: db.Collection("company_config"),
	}
}

// Get retrieves the company configuration (there should only be one)
func (r *CompanyConfigRepository) Get(ctx context.Context) (*domain.CompanyConfig, error) {
	var config domain.CompanyConfig
	err := r.collection.FindOne(ctx, bson.M{}).Decode(&config)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrCompanyConfigNotFound
		}
		return nil, err
	}
	return &config, nil
}

// Create creates a new company configuration
func (r *CompanyConfigRepository) Create(ctx context.Context, config *domain.CompanyConfig) error {
	// Check if config already exists
	existing, err := r.Get(ctx)
	if err != nil && err != domain.ErrCompanyConfigNotFound {
		return err
	}
	if existing != nil {
		return domain.ErrCompanyConfigAlreadyExists
	}

	config.CreatedAt = time.Now()
	config.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, config)
	if err != nil {
		return err
	}

	config.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

// Update updates the company configuration (upsert if doesn't exist)
func (r *CompanyConfigRepository) Update(ctx context.Context, config *domain.CompanyConfig) error {
	// Get existing config to preserve ID
	existing, err := r.Get(ctx)
	if err != nil && err != domain.ErrCompanyConfigNotFound {
		return err
	}

	config.UpdatedAt = time.Now()

	// If no existing config, create new one
	if existing == nil {
		config.CreatedAt = time.Now()
		return r.Create(ctx, config)
	}

	// Update existing config
	config.ID = existing.ID
	config.CreatedAt = existing.CreatedAt

	objectID, err := primitive.ObjectIDFromHex(config.ID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": config}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return domain.ErrCompanyConfigNotFound
	}

	return nil
}

// GetOrCreate retrieves the config or creates a default one if it doesn't exist
func (r *CompanyConfigRepository) GetOrCreate(ctx context.Context) (*domain.CompanyConfig, error) {
	config, err := r.Get(ctx)
	if err == nil {
		return config, nil
	}

	if err != domain.ErrCompanyConfigNotFound {
		return nil, err
	}

	// Create default configuration
	defaultConfig := &domain.CompanyConfig{
		CompanyName: "Your Company Name",
		WorkingHours: domain.WorkingHours{
			WorkingDays: []int{1, 2, 3, 4, 5}, // Monday to Friday
			OpenTime:    "09:00",
			CloseTime:   "17:00",
			Timezone:    "Europe/Oslo",
		},
		ShiftRequirements: []domain.ShiftRequirement{
			{
				ShiftType:    domain.ShiftTypeMorning,
				MinEmployees: 1,
				MaxEmployees: 2,
				Description:  "Morning shift coverage",
			},
			{
				ShiftType:    domain.ShiftTypeAfternoon,
				MinEmployees: 1,
				MaxEmployees: 2,
				Description:  "Afternoon shift coverage",
			},
		},
		SchedulingPolicies: domain.SchedulingPolicies{
			MaxConsecutiveDays:     5,
			MinRestHours:           12,
			AllowOvertime:          true,
			MaxOvertimeHours:       20,
			WeekendConsentRequired: true,
			FairDistribution:       true,
		},
		AIContext: "Please ensure fair distribution of shifts and respect employee availability preferences.",
	}

	if err := r.Create(ctx, defaultConfig); err != nil {
		return nil, err
	}

	return defaultConfig, nil
}
