package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InitDB initializes the MongoDB connection and creates indexes
func InitDB(mongoURI, dbName string) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create client options
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping the database to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	db := client.Database(dbName)

	// Create indexes
	if err := createIndexes(ctx, db); err != nil {
		return nil, fmt.Errorf("failed to create indexes: %w", err)
	}

	return db, nil
}

func createIndexes(ctx context.Context, db *mongo.Database) error {
	// Employees collection indexes
	employeesCollection := db.Collection("employees")

	// Email unique index
	_, err := employeesCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    map[string]interface{}{"email": 1},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("failed to create email index: %w", err)
	}

	// Active index
	_, err = employeesCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: map[string]interface{}{"active": 1},
	})
	if err != nil {
		return fmt.Errorf("failed to create active index: %w", err)
	}

	// Schedules collection indexes
	schedulesCollection := db.Collection("schedules")

	// Period index
	_, err = schedulesCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: map[string]interface{}{
			"period_start": 1,
			"period_end":   1,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create period index: %w", err)
	}

	// Status index
	_, err = schedulesCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: map[string]interface{}{"status": 1},
	})
	if err != nil {
		return fmt.Errorf("failed to create status index: %w", err)
	}

	return nil
}
