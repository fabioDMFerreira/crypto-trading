package domain

import (
	"time"

	"github.com/fabiodmferreira/crypto-trading/assets"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventsLog interface {
	Create(logType, message string) error
}

// Trader buys and sells
type Trader interface {
	Buy(amount, price float32, buyTime time.Time)
	Sell(asset *assets.Asset, price float32)
}

type Account interface {
	Withdraw(amount float32) error
	Deposit(amount float32) error
}

type AssetsRepository interface {
	FindAll() (*[]assets.Asset, error)
	FindCheaperAssetPrice() (float32, error)
	Sell(id primitive.ObjectID, price float32) error
	Create(asset *assets.Asset) error
}
