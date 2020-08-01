package domain

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Repository is a generic interface to be used by other repositories
type Repository interface {
	FindAll(documents interface{}, query interface{}, opts *options.FindOptions) error
	Aggregate(documents interface{}, pipelineOptions mongo.Pipeline) error
	FindOne(document interface{}, query interface{}, opts *options.FindOneOptions) error
	InsertOne(document interface{}) error
	UpdateOne(query interface{}, update interface{}) error
	DeleteByID(id string) error
	BulkUpsert(documents []bson.M) error
	BulkCreate(documents *[]bson.M) error
	BulkDelete(filter bson.M) error
	BulkUpdate(filter bson.M, update bson.M) error
}
