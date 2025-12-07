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

type employeeRepository struct {
	collection *mongo.Collection
}

// NewEmployeeRepository creates a new MongoDB employee repository
func NewEmployeeRepository(db *mongo.Database) repository.EmployeeRepository {
	return &employeeRepository{
		collection: db.Collection("employees"),
	}
}

func (r *employeeRepository) Create(ctx context.Context, employee *domain.Employee) error {
	if employee.ID == "" {
		employee.ID = uuid.New().String()
	}

	now := time.Now()
	employee.CreatedAt = now
	employee.UpdatedAt = now
	employee.Active = true

	_, err := r.collection.InsertOne(ctx, employee)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return domain.ErrEmployeeAlreadyExists
		}
		return err
	}

	return nil
}

func (r *employeeRepository) GetByID(ctx context.Context, id string) (*domain.Employee, error) {
	var employee domain.Employee

	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&employee)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrEmployeeNotFound
		}
		return nil, err
	}

	return &employee, nil
}

func (r *employeeRepository) GetAll(ctx context.Context) ([]domain.Employee, error) {
	opts := options.Find().SetSort(bson.D{{Key: "name", Value: 1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var employees []domain.Employee
	if err := cursor.All(ctx, &employees); err != nil {
		return nil, err
	}

	if employees == nil {
		employees = []domain.Employee{}
	}

	return employees, nil
}

func (r *employeeRepository) GetActive(ctx context.Context) ([]domain.Employee, error) {
	opts := options.Find().SetSort(bson.D{{Key: "name", Value: 1}})

	cursor, err := r.collection.Find(ctx, bson.M{"active": true}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var employees []domain.Employee
	if err := cursor.All(ctx, &employees); err != nil {
		return nil, err
	}

	if employees == nil {
		employees = []domain.Employee{}
	}

	return employees, nil
}

func (r *employeeRepository) Update(ctx context.Context, employee *domain.Employee) error {
	employee.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"name":             employee.Name,
			"email":            employee.Email,
			"role":             employee.Role,
			"role_description": employee.RoleDescription,
			"monthly_hours":    employee.MonthlyHours,
			"active":           employee.Active,
			"updated_at":       employee.UpdatedAt,
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"id": employee.ID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return domain.ErrEmployeeNotFound
	}

	return nil
}

func (r *employeeRepository) Delete(ctx context.Context, id string) error {
	update := bson.M{
		"$set": bson.M{
			"active":     false,
			"updated_at": time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"id": id}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return domain.ErrEmployeeNotFound
	}

	return nil
}

func (r *employeeRepository) GetByEmail(ctx context.Context, email string) (*domain.Employee, error) {
	var employee domain.Employee

	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&employee)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrEmployeeNotFound
		}
		return nil, err
	}

	return &employee, nil
}
