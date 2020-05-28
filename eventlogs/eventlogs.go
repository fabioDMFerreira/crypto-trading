package eventlogs

import (
	"log"
	"time"

	"github.com/fabiodmferreira/crypto-trading/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// EventLog
type EventLog struct {
	ID          primitive.ObjectID `bson:"_id" json:"_id"`
	EventName   string             `json:"eventName"`
	Message     string             `json:"message"`
	Notified    bool               `json:"notified"`
	DateCreated time.Time          `json:"dateCreated"`
}

type EventLogsRepository struct {
	collection *mongo.Collection
}

// NewAssetsRepository returns an instance of OrdersRepository
func NewEventLogsRepository(collection *mongo.Collection) *EventLogsRepository {
	return &EventLogsRepository{
		collection,
	}
}

// Create inserts a new asset in collection
func (or *EventLogsRepository) Create(eventName, message string) error {
	event := &EventLog{
		ID:          primitive.NewObjectID(),
		EventName:   eventName,
		Message:     message,
		Notified:    false,
		DateCreated: time.Now(),
	}

	log.Println(message)

	ctx := db.NewMongoQueryContext()
	_, err := or.collection.InsertOne(ctx, event)
	return err
}

// FindAllToNotify returns every event log that needs to be notified
func (or *EventLogsRepository) FindAllToNotify() ([]*EventLog, error) {
	ctx := db.NewMongoQueryContext()
	cur, err := or.collection.Find(ctx, bson.D{{"notified", false}})

	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)
	var results []*EventLog
	if err = cur.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

// MarkNotified marks every eventLog as notified
func (or *EventLogsRepository) MarkNotified(ids []primitive.ObjectID) error {
	ctx := db.NewMongoQueryContext()

	filter := bson.D{{"_id", bson.M{"$in": ids}}}
	update := bson.D{{"$set", bson.D{{"notified", true}}}}
	_, err := or.collection.UpdateMany(ctx, filter, update)

	return err
}
