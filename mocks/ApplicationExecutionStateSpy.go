package mocks

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ApplicationExecutionStatesRepositorySpy struct {
	CreateCalls     [][]interface{}
	AggregateCalls  []interface{}
	BulkCreateCalls []interface{}
	BulkDeleteCalls []string
}

func (a *ApplicationExecutionStatesRepositorySpy) Create(date time.Time, executionID primitive.ObjectID, state interface{}) error {
	a.CreateCalls = append(a.CreateCalls, []interface{}{date, executionID, state})
	return nil
}

func (a *ApplicationExecutionStatesRepositorySpy) Aggregate(pipeline mongo.Pipeline) (*[]bson.M, error) {
	a.AggregateCalls = append(a.AggregateCalls, pipeline)
	return &[]bson.M{}, nil
}

func (a *ApplicationExecutionStatesRepositorySpy) BulkCreate(documents *[]bson.M) error {
	a.BulkCreateCalls = append(a.BulkCreateCalls, documents)
	return nil
}

func (a *ApplicationExecutionStatesRepositorySpy) BulkDelete(id string) error {
	a.BulkDeleteCalls = append(a.BulkDeleteCalls, id)
	return nil
}
