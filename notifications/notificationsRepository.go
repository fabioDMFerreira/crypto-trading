package notifications

import (
	"time"

	"github.com/fabiodmferreira/crypto-trading/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Notification domain.Notification

type Repository struct {
	repo domain.Repository
}

// NewAssetsRepository returns an instance of OrdersRepository
func NewRepository(repo domain.Repository) *Repository {
	return &Repository{
		repo,
	}
}

// Create inserts a new notification in collection
func (r *Repository) Create(notification *domain.Notification) error {
	return r.repo.InsertOne(notification)
}

func (r *Repository) FindLastEventLogsNotificationDate() (time.Time, error) {
	var document Notification

	opts := options.FindOne().SetSort(bson.M{"createdAt": -1})

	err := r.repo.FindOne(&document, bson.M{"notificationType": "eventlogs"}, opts)

	if err != nil {
		return time.Now().AddDate(-1, 0, 0), err
	}

	return document.CreatedAt, nil
}

// Sent updates notification as sent
func (r *Repository) Sent(id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"sent": true}}

	return r.repo.UpdateOne(filter, update)
}

// BulkDelete removes multiple notifications
func (r *Repository) BulkDelete(filter bson.M) error {
	return r.repo.BulkDelete(filter)
}

// BulkDeleteByApplicationID deletes all notifications associated with an application id
func (r *Repository) BulkDeleteByApplicationID(id string) error {

	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}

	filter := bson.M{"applicationID": oid}

	return r.BulkDelete(filter)
}
