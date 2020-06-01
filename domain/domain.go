package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventsLog interface {
	FindAllToNotify() (*[]EventLog, error)
	Create(logType, message string) error
	MarkNotified(ids []primitive.ObjectID) error
}

// Trader buys and sells
type Trader interface {
	Buy(amount, price float32, buyTime time.Time) error
	Sell(asset *Asset, price float32) error
}

type AssetsRepositoryReader interface {
	FindAll() (*[]Asset, error)
	FindCheaperAssetPrice() (float32, error)
	GetBalance(startDate, endDate time.Time) (float32, error)
}

type AssetsRepository interface {
	AssetsRepositoryReader
	Sell(id primitive.ObjectID, price float32) error
	Create(asset *Asset) error
}

type Broker interface {
	AddBuyOrder(amount, price float32) error
	AddSellOrder(amount, price float32) error
}

type AccountServiceReader interface {
	GetAmount() (float32, error)
}

type AccountService interface {
	AccountServiceReader
	Withdraw(amount float32) error
	Deposit(amount float32) error
}
