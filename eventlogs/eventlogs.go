package eventlogs

import (
	"log"
	"time"

	"github.com/fabiodmferreira/crypto-trading/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// EventLogsRepository fetches and stores events
type EventLogsRepository struct {
	repo  domain.Repository
	appID primitive.ObjectID
}

// NewEventLogsRepository returns an instance of Repository
func NewEventLogsRepository(repo domain.Repository, appID primitive.ObjectID) *EventLogsRepository {
	return &EventLogsRepository{
		repo,
		appID,
	}
}

// Create inserts a new event in collection
func (or *EventLogsRepository) Create(eventName, message string) error {
	event := &domain.EventLog{
		ID:            primitive.NewObjectID(),
		EventName:     eventName,
		Message:       message,
		Notified:      false,
		CreatedAt:     time.Now(),
		ApplicationID: or.appID,
	}

	log.Println(message)

	return or.repo.InsertOne(event)
}

// FindAllToNotify returns every event log that needs to be notified
func (or *EventLogsRepository) FindAllToNotify() (*[]domain.EventLog, error) {
	var results []domain.EventLog

	err := or.repo.FindAll(&results, bson.M{"notified": false, "applicationID": or.appID}, nil)

	return &results, err
}

// MarkNotified marks every eventLog as notified
func (or *EventLogsRepository) MarkNotified(ids []primitive.ObjectID) error {
	filter := bson.M{"_id": bson.M{"$in": ids}, "applicationID": or.appID}
	update := bson.M{"$set": bson.M{"notified": true}}

	return or.repo.BulkUpdate(filter, update)
}

// FindAll returns all log events that match the filter
func (e *EventLogsRepository) FindAll(filter interface{}) (*[]domain.EventLog, error) {
	var events []domain.EventLog

	err := e.repo.FindAll(&events, filter, &options.FindOptions{Sort: bson.M{"createdat": -1}})

	return &events, err
}

// BulkDeleteByApplicationID deletes rows related with an application id.
func (e *EventLogsRepository) BulkDeleteByApplicationID(id string) error {

	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}

	filter := bson.M{"applicationID": oid}

	return e.repo.BulkDelete(filter)
}
