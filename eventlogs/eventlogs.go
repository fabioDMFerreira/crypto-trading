package eventlogs

import (
	"log"
	"time"

	"github.com/fabiodmferreira/crypto-trading/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// EventLogsRepository fetches and stores events
type EventLogsRepository struct {
	repo domain.Repository
}

// NewEventLogsRepository returns an instance of Repository
func NewEventLogsRepository(repo domain.Repository) *EventLogsRepository {
	return &EventLogsRepository{
		repo,
	}
}

// Create inserts a new event in collection
func (or *EventLogsRepository) Create(eventName, message string) error {
	event := &domain.EventLog{
		ID:        primitive.NewObjectID(),
		EventName: eventName,
		Message:   message,
		Notified:  false,
		CreatedAt: time.Now(),
	}

	log.Println(message)

	return or.repo.InsertOne(event)
}

// FindAllToNotify returns every event log that needs to be notified
func (or *EventLogsRepository) FindAllToNotify() (*[]domain.EventLog, error) {
	var results []domain.EventLog

	err := or.repo.FindAll(&results, bson.M{"notified": false}, nil)

	return &results, err
}

// MarkNotified marks every eventLog as notified
func (or *EventLogsRepository) MarkNotified(ids []primitive.ObjectID) error {
	filter := bson.M{"_id": bson.M{"$in": ids}}
	update := bson.M{"$set": bson.M{"notified": true}}

	return or.repo.BulkUpdate(filter, update)
}
