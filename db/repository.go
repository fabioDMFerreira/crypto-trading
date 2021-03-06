package db

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
	ctx, cancel := NewMongoQueryContext()
	cur, err := r.collection.Find(ctx, filter, opts)

	if err != nil {
		cancel()
		return err
	}

	defer cur.Close(ctx)

	if err = cur.All(ctx, documents); err != nil {
		cancel()
		return err
	}

	return nil
}

// Aggregate returns rows aggregated
func (r *Repository) Aggregate(documents interface{}, pipelineOptions mongo.Pipeline) error {
	ctx, cancel := NewMongoQueryContext()
	cur, err := r.collection.Aggregate(ctx, pipelineOptions)

	if err != nil {
		cancel()
		return err
	}

	defer cur.Close(ctx)

	if err = cur.All(ctx, documents); err != nil {
		cancel()
		return err
	}

	return nil
}

// FindOne returns one row that match criteria
func (r *Repository) FindOne(document interface{}, filter interface{}, opts *options.FindOneOptions) error {
	ctx, cancel := NewMongoQueryContext()

	err := r.collection.FindOne(ctx, filter, opts).Decode(document)

	if err != nil {
		cancel()
	}

	return err
}

// InsertOne creates one document
func (r *Repository) InsertOne(document interface{}) error {
	ctx, cancel := NewMongoQueryContext()

	_, err := r.collection.InsertOne(ctx, document)

	if err != nil {
		cancel()
	}

	return err
}

// UpdateOne updates one document found with match criteria
func (r *Repository) UpdateOne(filter interface{}, update interface{}) error {
	ctx, cancel := NewMongoQueryContext()

	_, err := r.collection.UpdateOne(ctx, filter, update)

	if err != nil {
		cancel()
	}

	return err
}

// DeleteByID delete one document by ID
func (r *Repository) DeleteByID(id string) error {
	ctx, cancel := NewMongoQueryContext()

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf(fmt.Sprint("primitive.ObjectIDFromHex ERROR:", err))
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": idPrimitive})

	if err != nil {
		cancel()
	}

	return err
}

// BulkUpsert creates multiple documents if there are no documents with the same data
func (r *Repository) BulkUpsert(documents []bson.M) error {
	ctx, cancel := NewMongoQueryContext()

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

	if err != nil {
		cancel()
	}

	return err
}

// BulkCreate creates multiple documents
func (r *Repository) BulkCreate(documents *[]bson.M) error {
	ctx, cancel := NewMongoQueryContext()

	bulkOptions := options.InsertManyOptions{}

	ds := []interface{}{}

	for _, d := range *documents {
		ds = append(ds, d)
	}

	_, err := r.collection.InsertMany(ctx, ds, &bulkOptions)

	if err != nil {
		cancel()
	}

	return err
}

// BulkDelete deletes multiple documents
func (r *Repository) BulkDelete(filter bson.M) error {
	ctx, cancel := NewMongoQueryContext()
	defer cancel()

	_, err := r.collection.DeleteMany(ctx, filter)

	return err
}

// BulkUpdate updates multiple documents
func (r *Repository) BulkUpdate(filter bson.M, update bson.M) error {
	ctx, cancel := NewMongoQueryContext()
	defer cancel()

	_, err := r.collection.UpdateMany(ctx, filter, update)

	return err
}
