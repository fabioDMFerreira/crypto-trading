package applicationExecutionStates

import (
	"time"

	"github.com/fabiodmferreira/crypto-trading/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Repository stores and gets assets prices
type Repository struct {
	repo domain.Repository
}

// NewRepository returns an instance of ApplicationExecutionStatesRepository
func NewRepository(repo domain.Repository) *Repository {

	return &Repository{
		repo,
	}
}

// Aggregate returns assets prices aggregated
func (r *Repository) Aggregate(pipeline mongo.Pipeline) (*[]bson.M, error) {
	var results []bson.M
	err := r.repo.Aggregate(&results, pipeline)

	if err != nil {
		return nil, err
	}

	return &results, nil
}

// FindAll returns assets prices
func (r *Repository) FindAll(filter interface{}) (*[]domain.ApplicationExecutionState, error) {
	var results []domain.ApplicationExecutionState
	err := r.repo.FindAll(&results, filter, nil)

	if err != nil {
		return nil, err
	}

	return &results, nil
}

// Create stores an application execution state
func (r *Repository) Create(date time.Time, executionID primitive.ObjectID, state interface{}) error {
	ApplicationExecutionState := domain.ApplicationExecutionState{ID: primitive.NewObjectID(), Date: date, ExecutionID: executionID, State: state}

	return r.repo.InsertOne(ApplicationExecutionState)
}

// BulkCreate creates multiple documents
func (r *Repository) BulkCreate(documents *[]bson.M) error {
	return r.repo.BulkCreate(documents)
}

// BulkDelete deletes multiple documents
func (r *Repository) BulkDelete(id string) error {

	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}

	filter := bson.M{"executionId": oid}

	return r.repo.BulkDelete(filter)
}
