package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AssetsRepositoryReader fetches assets data
type AssetsRepositoryReader interface {
	FindAll() (*[]Asset, error)
	FindCheaperAssetPrice() (float32, error)
	GetBalance(startDate, endDate time.Time) (float32, error)
}

// AssetsRepository stores and fetches assets
type AssetsRepository interface {
	AssetsRepositoryReader
	Sell(id primitive.ObjectID, price float32) error
	Create(asset *Asset) error
}

// AccountsRepository stores and fetches accounts
type AccountsRepository interface {
	FindById(id primitive.ObjectID) (*Account, error)
	FindByBroker(broker string) (*Account, error)
	Create(broker string, amount float32) (*Account, error)
	Withdraw(id primitive.ObjectID, amount float32) error
	Deposit(id primitive.ObjectID, amount float32) error
}
