package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// AssetPrice represents the price of an asset on a moment
type AssetPrice struct {
	ID    primitive.ObjectID `bson:"_id" json:"_id"`
	Date  time.Time          `json:"date"`
	Value float32            `json:"value"`
	Asset string             `json:"asset"`
}

// AssetPriceRepository stores and gets assets prices
type AssetPriceRepository interface {
	Create(date time.Time, value float32, asset string) error
	FindAll(filter interface{}) (*[]AssetPrice, error)
	Aggregate(pipeline mongo.Pipeline) (*[]bson.M, error)
}

// AssetPriceGroupByDate is a group id struct
type AssetPriceGroupByDate struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`
}

// AssetPriceAggregatedByDate is the output of the aggregate query that groups prices per date
type AssetPriceAggregatedByDate struct {
	ID    AssetPriceGroupByDate `json:"_id"`
	Price float32               `json:"price"`
}