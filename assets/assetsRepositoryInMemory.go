package assets

import (
	"time"

	"github.com/fabiodmferreira/crypto-trading/domain"
)

// AssetsRepositoryInMemory stores assets in memory
type AssetsRepositoryInMemory struct {
	Assets []domain.Asset
}

// NewAssetsRepositoryInMemory returns an instance of NewAssetsRepositoryInMemory
func NewAssetsRepositoryInMemory() *AssetsRepositoryInMemory {
	return &AssetsRepositoryInMemory{}
}

// FindPendingAssets returns all assets sold stored
func (ar *AssetsRepositoryInMemory) FindPendingAssets(accountID string) (*[]domain.Asset, error) {

	pendingAssets := []domain.Asset{}

	for _, asset := range ar.Assets {
		if !asset.Sold {
			pendingAssets = append(pendingAssets, asset)
		}
	}
	return &pendingAssets, nil
}

// FindAll returns all assets stored
func (ar *AssetsRepositoryInMemory) FindAll(accountID string) (*[]domain.Asset, error) {
	return &ar.Assets, nil
}

// FindCheaperAssetPrice returns the lowest price of non sold assets
func (ar *AssetsRepositoryInMemory) FindCheaperAssetPrice(accountID string) (float32, error) {
	var minimumPrice float32

	for _, asset := range ar.Assets {
		if asset.Sold == false && minimumPrice > asset.BuyPrice {
			minimumPrice = asset.BuyPrice
		}
	}

	return minimumPrice, nil
}

// GetBalance mocks the returning of balance between two dates
func (ar *AssetsRepositoryInMemory) GetBalance(accountID string, startDate, endDate time.Time) (float32, error) {
	return 0, nil
}

// Create creates an asset and stores it in a data structure
func (ar *AssetsRepositoryInMemory) Create(asset *domain.Asset) error {
	ar.Assets = append(ar.Assets, *asset)
	return nil
}

// Sell updates asset state to sold and other related attributes
func (ar *AssetsRepositoryInMemory) Sell(id string, price float32, sellTime time.Time) error {

	for index, asset := range ar.Assets {
		if asset.ID.Hex() == id {
			ar.Assets[index].SellPrice = price
			ar.Assets[index].Sold = true
			ar.Assets[index].SellTime = sellTime
			break
		}
	}

	return nil
}

// CheckAssetWithCloserPriceExists checks whether exist an asset that has the same price within limits defined
func (ar *AssetsRepositoryInMemory) CheckAssetWithCloserPriceExists(accountID string, price, limit float32) (bool, error) {
	lowerLimit := price - (price * limit)
	upperLimit := price + (price * limit)

	for _, asset := range ar.Assets {
		if !asset.Sold && asset.BuyPrice > lowerLimit && asset.BuyPrice < upperLimit {
			return true, nil
		}
	}

	return false, nil
}
