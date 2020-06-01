package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Account
type Account struct {
	ID     primitive.ObjectID `bson:"_id" json:"_id"`
	Amount float32            `bson:"amount,truncate" json:"amount"`
	Broker string             `json:"broker"`
}

// Asset
type Asset struct {
	ID        primitive.ObjectID `bson:"_id" json:"_id"`
	Amount    float32            `bson:"amount,truncate" json:"amount"`
	BuyTime   time.Time          `json:"buyTime"`
	SellTime  time.Time          `json:"sellTime"`
	BuyPrice  float32            `bson:"buyPrice,truncate" json:"buyPrice"`
	SellPrice float32            `bson:"sellPrice,truncate" json:"sellPrice"`
	Sold      bool               `json:"sold"`
}

// EventLog
type EventLog struct {
	ID          primitive.ObjectID `bson:"_id" json:"_id"`
	EventName   string             `json:"eventName"`
	Message     string             `json:"message"`
	Notified    bool               `json:"notified"`
	DateCreated time.Time          `json:"dateCreated"`
}

// Notification
type Notification struct {
	ID                  primitive.ObjectID `bson:"_id" json:"_id"`
	Title               string             `json:"title"`
	Message             string             `json:"message"`
	To                  string             `json:"to"`
	DateCreated         time.Time          `json:"dateCreated"`
	NotificationType    string             `json:"notificationType"`
	NotificationChannel string             `json:"notificationChannel"`
	Sent                bool               `json:"sent"`
}

type NotificationsService interface {
	FindLastEventLogsNotificationDate() (time.Time, error)
	CreateEmailNotification(subject, message, notificationType string) error
}

type DecisionMaker interface {
	ShouldBuy(price float32, buyTime time.Time) (bool, error)
	ShouldSell(asset *Asset, price float32, buyTime time.Time) (bool, error)
	HowMuchAmountShouldBuy(price float32) (float32, error)
}
