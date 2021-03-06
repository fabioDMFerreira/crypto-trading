package assetsprices

import (
	"time"

	"github.com/fabiodmferreira/crypto-trading/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// RepositoryInMemory stores and gets assets prices in memory
type RepositoryInMemory struct {
	assetsPrices []domain.AssetPrice
}

// NewRepositoryInMemory returns assets prices repository mock
func NewRepositoryInMemory() *RepositoryInMemory {
	return &RepositoryInMemory{}
}

// FindAll returns assets prices
func (r *RepositoryInMemory) FindAll(filter interface{}) (*[]domain.AssetPrice, error) {
	return &r.assetsPrices, nil
}

// FindOne returns an asset price
func (r *RepositoryInMemory) FindOne(date time.Time, value float32, asset string) (*domain.AssetPrice, error) {
	var assetPrice domain.AssetPrice

	for _, assetP := range r.assetsPrices {
		if assetP.Date == date && assetP.Close == value && assetP.Asset == asset {
			assetPrice = assetP
		}
	}

	return &assetPrice, nil
}

// Create stores an asset price
func (r *RepositoryInMemory) Create(ohlc *domain.OHLC, asset string) error {
	foundDocument, err := r.FindOne(ohlc.Time, ohlc.Close, asset)

	if err != nil {
		return err
	}

	if foundDocument != nil {
		return nil
	}

	assetPrice := domain.AssetPrice{Date: ohlc.Time, Close: ohlc.Close, Asset: asset}

	r.assetsPrices = append(r.assetsPrices, assetPrice)

	return nil
}

// Aggregate returns assets prices aggregated
func (r *RepositoryInMemory) Aggregate(pipeline mongo.Pipeline) (*[]bson.M, error) {
	var results []bson.M

	return &results, nil
}

// GetLastAssetsPrices stub
func (r *RepositoryInMemory) GetLastAssetsPrices(asset string, limit int) (*[]domain.AssetPrice, error) {
	return nil, nil
}

// BulkCreate stub
func (r *RepositoryInMemory) BulkCreate(documents *[]bson.M) error {
	return nil
}
