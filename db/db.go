package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	ASSETS_COLLECTION                       = "assets"
	EVENT_LOGS_COLLECTION                   = "eventlogs"
	NOTIFICATIONS_COLLECTION                = "notifications"
	ACCOUNTS_COLLECTION                     = "accounts"
	BENCHMARKS_COLLECTION                   = "benchmarks"
	ASSETS_PRICES_COLLECTION                = "assetsprices"
	APPLICATION_EXECUTION_STATES_COLLECTION = "applicationExecutionStates"
)

func NewMongoQueryContext() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	return ctx
}

func ConnectDB(mongoUrl string) (*mongo.Client, error) {

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUrl))

	if err != nil {
		return nil, err
	}

	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	err = client.Ping(ctx, readpref.Primary())

	if err != nil {
		return nil, err
	}

	return client, nil
}

// Repository is a generic mongo service to used by other repositories
type Repository struct {
	collection *mongo.Collection
}

// NewRepository returns an instance of Repository
func NewRepository(c *mongo.Collection) *Repository {
	return &Repository{c}
}

// FindAll returns rows that match criteria
func (r *Repository) FindAll(documents interface{}, filter interface{}, opts *options.FindOptions) error {
	ctx := NewMongoQueryContext()
	cur, err := r.collection.Find(ctx, filter, opts)

	if err != nil {
		return err
	}

	defer cur.Close(ctx)

	if err = cur.All(ctx, documents); err != nil {
		return err
	}

	return nil
}

// Aggregate returns rows aggregated
func (r *Repository) Aggregate(documents interface{}, pipelineOptions mongo.Pipeline) error {
	ctx := NewMongoQueryContext()
	cur, err := r.collection.Aggregate(ctx, pipelineOptions)

	if err != nil {
		return err
	}

	defer cur.Close(ctx)

	if err = cur.All(ctx, documents); err != nil {
		return err
	}

	return nil
}

// FindOne returns one row that match criteria
func (r *Repository) FindOne(document interface{}, filter interface{}, opts *options.FindOneOptions) error {
	ctx := NewMongoQueryContext()

	err := r.collection.FindOne(ctx, filter, opts).Decode(document)

	return err
}

// InsertOne creates one document
func (r *Repository) InsertOne(document interface{}) error {
	ctx := NewMongoQueryContext()

	_, err := r.collection.InsertOne(ctx, document)

	return err
}

// UpdateOne updates one document found with match criteria
func (r *Repository) UpdateOne(filter interface{}, update interface{}) error {
	ctx := NewMongoQueryContext()

	_, err := r.collection.UpdateOne(ctx, filter, update)

	return err
}

// DeleteByID delete one document by ID
func (r *Repository) DeleteByID(id string) error {
	ctx := NewMongoQueryContext()

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf(fmt.Sprint("primitive.ObjectIDFromHex ERROR:", err))
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": idPrimitive})

	return err
}

// BulkUpsert creates multiple documents if there are no documents with the same data
func (r *Repository) BulkUpsert(documents []bson.M) error {
	ctx := NewMongoQueryContext()

	var operations []mongo.WriteModel

	for _, document := range documents {
		operation := mongo.NewUpdateOneModel()

		operation.SetFilter(document)
		operation.SetUpdate(document)
		operation.SetUpsert(true)

		operations = append(operations, operation)
	}

	bulkOptions := options.BulkWriteOptions{}

	_, err := r.collection.BulkWrite(ctx, operations, &bulkOptions)

	return err
}

// BulkCreate creates multiple documents
func (r *Repository) BulkCreate(documents *[]bson.M) error {
	ctx := NewMongoQueryContext()

	bulkOptions := options.InsertManyOptions{}

	ds := []interface{}{}

	for _, d := range *documents {
		ds = append(ds, d)
	}

	_, err := r.collection.InsertMany(ctx, ds, &bulkOptions)

	return err
}

// BulkDelete deletes multiple documents
func (r *Repository) BulkDelete(filter bson.M) error {
	ctx := NewMongoQueryContext()

	_, err := r.collection.DeleteMany(ctx, filter)

	return err
}
