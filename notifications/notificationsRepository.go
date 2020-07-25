package notifications

import (
	"time"

	"github.com/fabiodmferreira/crypto-trading/db"
	"github.com/fabiodmferreira/crypto-trading/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Notification domain.Notification

type Repository struct {
	collection *mongo.Collection
}

// NewAssetsRepository returns an instance of OrdersRepository
func NewRepository(collection *mongo.Collection) *Repository {
	return &Repository{
		collection,
	}
}

// Create inserts a new notification in collection
func (or *Repository) Create(notification *Notification) error {
	ctx := db.NewMongoQueryContext()
	_, err := or.collection.InsertOne(ctx, notification)
	return err
}

func (or *Repository) FindLastEventLogsNotificationDate() (time.Time, error) {
	ctx := db.NewMongoQueryContext()

	opts := options.FindOne().SetSort(bson.D{{"createdat", -1}})
	var foundDocument Notification
	err := or.collection.FindOne(ctx, bson.D{{"notificationtype", "eventlogs"}}, opts).Decode(&foundDocument)

	if err != nil {
		return time.Now().AddDate(-1, 0, 0), err
	}

	return foundDocument.CreatedAt, nil
}

// Sent updates notification as sent
func (or *Repository) Sent(id primitive.ObjectID) error {
	ctx := db.NewMongoQueryContext()

	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"send", true}}}}
	_, err := or.collection.UpdateOne(ctx, filter, update)

	return err
}
