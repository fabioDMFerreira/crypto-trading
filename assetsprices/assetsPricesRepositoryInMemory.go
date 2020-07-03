package assetsprices

import (
	"time"

	"github.com/fabiodmferreira/crypto-trading/domain"
)

// RepositoryInMemory stores and gets assets prices in memory
type RepositoryInMemory struct {
	assetsPrices []domain.AssetPrice
}

// FindOne returns an asset price
func (r *RepositoryInMemory) FindOne(date time.Time, value float32, asset string) (*domain.AssetPrice, error) {
	var assetPrice domain.AssetPrice

	for _, assetP := range r.assetsPrices {
		if assetP.Date == date && assetP.Value == value && assetP.Asset == asset {
			assetPrice = assetP
		}
	}

	return &assetPrice, nil
}

// Create stores an asset price
func (r *RepositoryInMemory) Create(date time.Time, value float32, asset string) error {
	foundDocument, err := r.FindOne(date, value, asset)

	if err != nil {
		return err
	}

	if foundDocument != nil {
		return nil
	}

	assetPrice := domain.AssetPrice{Date: date, Value: value, Asset: asset}

	r.assetsPrices = append(r.assetsPrices, assetPrice)

	return nil
}
