package applicationExecutionStates

import (
	"time"

	"github.com/fabiodmferreira/crypto-trading/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Repository stores and gets assets prices
type RepositoryInMemory struct {
}

// NewRepository returns an instance of ApplicationExecutionStatesRepository
func NewRepositoryInMemory() *Repository {

	return &Repository{}
}

// Aggregate returns assets prices aggregated
func (r *RepositoryInMemory) Aggregate(pipeline mongo.Pipeline) (*[]bson.M, error) {
	var results []bson.M

	return &results, nil
}

// FindAll returns assets prices
func (r *RepositoryInMemory) FindAll(filter interface{}) (*[]domain.ApplicationExecutionState, error) {
	var results []domain.ApplicationExecutionState

	return &results, nil
}

// Create stores an application execution state
func (r *RepositoryInMemory) Create(date time.Time, executionID primitive.ObjectID, state interface{}) error {
	return nil
}

// BulkCreate creates multiple documents
func (r *RepositoryInMemory) BulkCreate(documents *[]bson.M) error {
	return nil
}

// BulkDelete deletes multiple documents
func (r *RepositoryInMemory) BulkDelete(id string) error {
	return nil
}
