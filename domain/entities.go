package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Account has details about an exchange account
type Account struct {
	ID     primitive.ObjectID `bson:"_id" json:"_id"`
	Amount float32            `bson:"amount,truncate" json:"amount"`
	Broker string             `json:"broker"`
}

// Asset is a financial instrument
type Asset struct {
	ID        primitive.ObjectID `bson:"_id" json:"_id"`
	Amount    float32            `bson:"amount,truncate" json:"amount"`
	BuyTime   time.Time          `json:"buyTime"`
	SellTime  time.Time          `json:"sellTime"`
	BuyPrice  float32            `bson:"buyPrice,truncate" json:"buyPrice"`
	SellPrice float32            `bson:"sellPrice,truncate" json:"sellPrice"`
	Sold      bool               `json:"sold"`
}

// EventLog is a regist of an application event
type EventLog struct {
	ID          primitive.ObjectID `bson:"_id" json:"_id"`
	EventName   string             `json:"eventName"`
	Message     string             `json:"message"`
	Notified    bool               `json:"notified"`
	DateCreated time.Time          `json:"dateCreated"`
}

// Notification is a message sent to a specific user
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
