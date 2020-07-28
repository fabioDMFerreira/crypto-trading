package assetsprices

import (
	"time"

	"github.com/fabiodmferreira/crypto-trading/domain"
)

// ServiceMock is a mock service
type ServiceMock struct{}

// Create stub
func (a *ServiceMock) Create(date time.Time, value float32, asset string) error {
	return nil
}

// FetchAndStoreAssetPrices stub
func (a *ServiceMock) FetchAndStoreAssetPrices(asset string, endDate time.Time) error {
	return nil
}

// GetRemotePrices stub
func (a *ServiceMock) GetRemotePrices(startDate, endDate time.Time, asset string) (*domain.CoindeskResponse, error) {
	return nil, nil
}

// GetLastAssetsPrices stub
func (a *ServiceMock) GetLastAssetsPrices(asset string, limit int) (*[]AssetPrice, error) {
	return nil, nil
}
