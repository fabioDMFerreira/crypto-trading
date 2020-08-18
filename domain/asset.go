package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Asset is a financial instrument
type Asset struct {
	ID        primitive.ObjectID `bson:"_id" json:"_id"`
	Amount    float32            `bson:"amount,truncate" json:"amount"`
	BuyTime   time.Time          `json:"buyTime"`
	SellTime  time.Time          `json:"sellTime"`
	BuyPrice  float32            `bson:"buyPrice,truncate" json:"buyPrice"`
	SellPrice float32            `bson:"sellPrice,truncate" json:"sellPrice"`
	Sold      bool               `json:"sold"`
	AccountID primitive.ObjectID `bson:"accountID" json:"accountID"`
}

// AssetsRepositoryReader fetches assets data
type AssetsRepositoryReader interface {
	FindAll(accountID string) (*[]Asset, error)
	FindPendingAssets(accountID string) (*[]Asset, error)
	FindCheaperAssetPrice(accountID string) (float32, error)
	CheckAssetWithCloserPriceExists(accountID string, price float32, limit float32) (bool, error)
	GetBalance(accountID string, startDate, endDate time.Time) (float32, error)
}

// AssetsRepository stores and fetches assets
type AssetsRepository interface {
	AssetsRepositoryReader
	Sell(id string, price float32, sellTime time.Time) error
	Create(asset *Asset) error
}
