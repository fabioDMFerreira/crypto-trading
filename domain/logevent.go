package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// EventLog is a regist of an application event
type EventLog struct {
	ID        primitive.ObjectID `bson:"_id" json:"_id"`
	EventName string             `json:"eventName"`
	Message   string             `json:"message"`
	Notified  bool               `json:"notified"`
	CreatedAt time.Time          `json:"createdAt"`
}

// EventsLog interacts with events
type EventsLog interface {
	FindAllToNotify() (*[]EventLog, error)
	Create(logType, message string) error
	MarkNotified(ids []primitive.ObjectID) error
}
