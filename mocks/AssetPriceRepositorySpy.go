package mocks

import (
	"time"

	"github.com/fabiodmferreira/crypto-trading/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AssetPriceRepositorySpy struct {
	AssetsPrices             *[]domain.AssetPrice
	CreateCalls              [][]interface{}
	FindAllCalls             []interface{}
	AggregateCalls           []interface{}
	GetLastAssetsPricesCalls [][]interface{}
	BulkCreateCalls          []interface{}
}

func (a *AssetPriceRepositorySpy) Create(date time.Time, value float32, asset string) error {
	a.CreateCalls = append(a.CreateCalls, []interface{}{date, value, asset})
	return nil
}

func (a *AssetPriceRepositorySpy) FindAll(filter interface{}) (*[]domain.AssetPrice, error) {
	a.FindAllCalls = append(a.FindAllCalls, filter)
	return &[]domain.AssetPrice{}, nil
}

func (a *AssetPriceRepositorySpy) Aggregate(pipeline mongo.Pipeline) (*[]bson.M, error) {
	a.AggregateCalls = append(a.AggregateCalls, pipeline)
	return &[]bson.M{}, nil
}

func (a *AssetPriceRepositorySpy) GetLastAssetsPrices(asset string, limit int) (*[]domain.AssetPrice, error) {
	a.GetLastAssetsPricesCalls = append(a.GetLastAssetsPricesCalls, []interface{}{asset, limit})
	return a.AssetsPrices, nil
}

func (a *AssetPriceRepositorySpy) BulkCreate(documents *[]bson.M) error {
	a.BulkCreateCalls = append(a.BulkCreateCalls, documents)
	return nil
}
