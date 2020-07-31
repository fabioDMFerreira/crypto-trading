package mocks

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RepositorySpy struct {
	FindAllCalls    [][]interface{}
	AggregateCalls  [][]interface{}
	FindOneCalls    [][]interface{}
	InsertOneCalls  []interface{}
	UpdateOneCalls  [][]interface{}
	DeleteByIDCalls []string
	BulkUpsertCalls []interface{}
	BulkCreateCalls []interface{}
	BulkDeleteCalls []interface{}
}

func (r *RepositorySpy) FindAll(documents interface{}, query interface{}, opts *options.FindOptions) error {
	r.FindAllCalls = append(r.FindAllCalls, []interface{}{documents, query, opts})
	return nil
}

func (r *RepositorySpy) Aggregate(documents interface{}, pipelineOptions mongo.Pipeline) error {
	r.AggregateCalls = append(r.AggregateCalls, []interface{}{documents, pipelineOptions})
	return nil
}

func (r *RepositorySpy) FindOne(document interface{}, query interface{}, opts *options.FindOneOptions) error {
	r.FindOneCalls = append(r.FindOneCalls, []interface{}{document, query, opts})
	return nil
}

func (r *RepositorySpy) InsertOne(document interface{}) error {
	r.InsertOneCalls = append(r.InsertOneCalls, document)
	return nil
}

func (r *RepositorySpy) UpdateOne(query interface{}, update interface{}) error {
	r.UpdateOneCalls = append(r.UpdateOneCalls, []interface{}{query, update})
	return nil
}

func (r *RepositorySpy) DeleteByID(id string) error {
	r.DeleteByIDCalls = append(r.DeleteByIDCalls, id)
	return nil
}

func (r *RepositorySpy) BulkUpsert(documents []bson.M) error {
	r.BulkUpsertCalls = append(r.BulkUpsertCalls, documents)
	return nil
}

func (r *RepositorySpy) BulkCreate(documents *[]bson.M) error {
	r.BulkCreateCalls = append(r.BulkCreateCalls, documents)
	return nil
}

func (r *RepositorySpy) BulkDelete(filter bson.M) error {
	r.BulkDeleteCalls = append(r.BulkDeleteCalls, filter)
	return nil
}
