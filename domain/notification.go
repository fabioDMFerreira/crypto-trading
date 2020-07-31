package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Notification is a message sent to a specific user
type Notification struct {
	ID                  primitive.ObjectID `bson:"_id" json:"_id"`
	Title               string             `json:"title"`
	Message             string             `json:"message"`
	To                  string             `json:"to"`
	NotificationType    string             `json:"notificationType"`
	NotificationChannel string             `json:"notificationChannel"`
	Sent                bool               `json:"sent"`
	CreatedAt           time.Time          `json:"createdAt"`
}

// NotificationOptions has notifications service options
type NotificationOptions struct {
	Receiver       string
	Sender         string
	SenderPassword string
}

// NotificationsService interacts with notifications
type NotificationsService interface {
	FindLastEventLogsNotificationDate() (time.Time, error)
	CreateEmailNotification(subject, message, notificationType string) error
	ShouldSendNotification() bool
}

// NotificationsRepository stores and gets notifications
type NotificationsRepository interface {
	Create(notification *Notification) error
	FindLastEventLogsNotificationDate() (time.Time, error)
	Sent(id primitive.ObjectID) error
}
