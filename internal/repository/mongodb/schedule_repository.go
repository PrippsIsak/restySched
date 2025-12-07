package mongodb

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/isak/restySched/internal/domain"
	"github.com/isak/restySched/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type scheduleRepository struct {
	collection *mongo.Collection
}

// NewScheduleRepository creates a new MongoDB schedule repository
func NewScheduleRepository(db *mongo.Database) repository.ScheduleRepository {
	return &scheduleRepository{
		collection: db.Collection("schedules"),
	}
}

func (r *scheduleRepository) Create(ctx context.Context, schedule *domain.Schedule) error {
	if schedule.ID == "" {
		schedule.ID = uuid.New().String()
	}

	now := time.Now()
	schedule.CreatedAt = now
	schedule.UpdatedAt = now

	_, err := r.collection.InsertOne(ctx, schedule)
	return err
}

func (r *scheduleRepository) GetByID(ctx context.Context, id string) (*domain.Schedule, error) {
	var schedule domain.Schedule

	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&schedule)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrScheduleNotFound
		}
		return nil, err
	}

	return &schedule, nil
}

func (r *scheduleRepository) GetAll(ctx context.Context) ([]domain.Schedule, error) {
	opts := options.Find().SetSort(bson.D{{Key: "period_start", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var schedules []domain.Schedule
	if err := cursor.All(ctx, &schedules); err != nil {
		return nil, err
	}

	if schedules == nil {
		schedules = []domain.Schedule{}
	}

	return schedules, nil
}

func (r *scheduleRepository) GetByPeriod(ctx context.Context, start, end time.Time) ([]domain.Schedule, error) {
	filter := bson.M{
		"period_start": bson.M{"$gte": start},
		"period_end":   bson.M{"$lte": end},
	}

	opts := options.Find().SetSort(bson.D{{Key: "period_start", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var schedules []domain.Schedule
	if err := cursor.All(ctx, &schedules); err != nil {
		return nil, err
	}

	if schedules == nil {
		schedules = []domain.Schedule{}
	}

	return schedules, nil
}

func (r *scheduleRepository) Update(ctx context.Context, schedule *domain.Schedule) error {
	schedule.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"period_start": schedule.PeriodStart,
			"period_end":   schedule.PeriodEnd,
			"employees":    schedule.Employees,
			"status":       schedule.Status,
			"sent_to_n8n":  schedule.SentToN8N,
			"sent_at":      schedule.SentAt,
			"updated_at":   schedule.UpdatedAt,
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"id": schedule.ID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return domain.ErrScheduleNotFound
	}

	return nil
}

func (r *scheduleRepository) Delete(ctx context.Context, id string) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return domain.ErrScheduleNotFound
	}

	return nil
}

func (r *scheduleRepository) MarkAsSent(ctx context.Context, id string) error {
	now := time.Now()

	update := bson.M{
		"$set": bson.M{
			"sent_to_n8n": true,
			"sent_at":     now,
			"status":      domain.ScheduleStatusSent,
			"updated_at":  now,
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"id": id}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return domain.ErrScheduleNotFound
	}

	return nil
}
