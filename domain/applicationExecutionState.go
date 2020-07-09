package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ApplicationExecutionState is used to store a snapshot of an application execution state in each price change
type ApplicationExecutionState struct {
	ID          primitive.ObjectID `bson:"_id" json:"_id"`
	ExecutionID primitive.ObjectID `bson:"executionId" json:"executionId"`
	Date        time.Time          `json:"date"`
	State       interface{}        `json:"executionState"`
}

// ApplicationExecutionStateRepository stores and gets application executions states
type ApplicationExecutionStateRepository interface {
	Create(date time.Time, executionID primitive.ObjectID, state interface{}) error
	Aggregate(pipeline mongo.Pipeline) (*[]bson.M, error)
	BulkCreate(documents *[]bson.M) error
	BulkDelete(id string) error
}
