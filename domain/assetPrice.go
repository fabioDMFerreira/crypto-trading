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
	Date  time.Time          `bson:"date" json:"date"`
	Value float32            `bson:"value" json:"value"`
	Asset string             `bson:"asset" json:"asset"`
}

// AssetPriceRepository stores and gets assets prices
type AssetPriceRepository interface {
	Create(date time.Time, value float32, asset string) error
	FindAll(filter interface{}) (*[]AssetPrice, error)
	Aggregate(pipeline mongo.Pipeline) (*[]bson.M, error)
	GetLastAssetsPrices(asset string, limit int) (*[]AssetPrice, error)
	BulkCreate(documents *[]bson.M) error
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

// AssetsPricesService provides assets prices related methods
type AssetsPricesService interface {
	GetRemotePrices(startDate, endDate time.Time, asset string) (*CoindeskResponse, error)
	GetLastAssetsPrices(asset string, limit int) (*[]AssetPrice, error)
	Create(date time.Time, value float32, asset string) error
	FetchAndStoreAssetPrices(asset string, endDate time.Time) error
}

// CoindeskResponse is the body of Coindesk HTTP Response
type CoindeskResponse struct {
	Iso      string      `json:"iso"`
	Name     string      `json:"name"`
	Slug     string      `json:"slug"`
	Interval string      `json:"interval"`
	Entries  [][]float64 `json:"entries"`
}

// CoindeskHTTPResponse is the response of Coindesk HTTP
type CoindeskHTTPResponse struct {
	StatusCode int              `json:"statusCode"`
	Message    string           `json:"message"`
	Data       CoindeskResponse `json:"data"`
}
