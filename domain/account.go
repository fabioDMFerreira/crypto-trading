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

// AccountsRepository stores and fetches accounts
type AccountsRepository interface {
	FindById(id string) (*Account, error)
	FindByBroker(broker string) (*Account, error)
	Create(broker string, amount float32) (*Account, error)
	Withdraw(id string, amount float32) error
	Deposit(id string, amount float32) error
}

// AccountServiceReader reads information about one account
type AccountServiceReader interface {
	GetAmount() (float32, error)
	FindPendingAssets() (*[]Asset, error)
	FindAllAssets() (*[]Asset, error)
	GetBalance(startDate, endDate time.Time) (float32, error)
	CheckAssetWithCloserPriceExists(price, limit float32) (bool, error)
}

// AccountService interacts with one account
type AccountService interface {
	AccountServiceReader
	Withdraw(amount float32) error
	Deposit(amount float32) error
	CreateAsset(amount, price float32, time time.Time) (*Asset, error)
	SellAsset(assetID string, price float32, time time.Time) error
}
