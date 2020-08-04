package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ApplicationState is the state hold by application execution state object
type ApplicationState struct {
	Average             float32 `json:"average"`
	StandardDeviation   float32 `json:"standardDeviation"`
	LowerBollingerBand  float32 `json:"lowerBollingerBand"`
	HigherBollingerBand float32 `json:"higherBollingerBand"`
	CurrentPrice        float32 `json:"currentPrice"`
	CurrentChange       float32 `json:"currentChange"`
}

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
	FindLast(filter interface{}) (*ApplicationExecutionState, error)
}
