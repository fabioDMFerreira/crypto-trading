package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// EventsLog interacts with events
type EventsLog interface {
	FindAllToNotify() (*[]EventLog, error)
	Create(logType, message string) error
	MarkNotified(ids []primitive.ObjectID) error
}

// Trader buys and sells assets
type Trader interface {
	Buy(amount, price float32, buyTime time.Time) error
	Sell(asset *Asset, price float32) error
}

// Broker add order to buy and sell assets in real brokers
type Broker interface {
	AddBuyOrder(amount, price float32) error
	AddSellOrder(amount, price float32) error
}

// AccountServiceReader reads information about one account
type AccountServiceReader interface {
	GetAmount() (float32, error)
}

// AccountService interacts with one account
type AccountService interface {
	AccountServiceReader
	Withdraw(amount float32) error
	Deposit(amount float32) error
}

// NotificationsService interacts with notifications
type NotificationsService interface {
	FindLastEventLogsNotificationDate() (time.Time, error)
	CreateEmailNotification(subject, message, notificationType string) error
}

// DecisionMaker makes decisions to buy or sell assets
type DecisionMaker interface {
	NewValue(price float32)
	ShouldBuy(price float32, buyTime time.Time) (bool, error)
	ShouldSell(asset *Asset, price float32, buyTime time.Time) (bool, error)
	HowMuchAmountShouldBuy(price float32) (float32, error)
}

// Statistics receives points and do statitics calculations
type Statistics interface {
	AddPoint(p float64)
	GetStandardDeviation() float64
	GetAverage() float64
}
