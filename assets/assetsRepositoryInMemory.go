package assets

import (
	"time"

	"github.com/fabiodmferreira/crypto-trading/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AssetsRepositoryInMemory stores assets in memory
type AssetsRepositoryInMemory struct {
	Assets []domain.Asset
}

// FindAll returns all assets not sold stored
func (ar *AssetsRepositoryInMemory) FindAll() (*[]domain.Asset, error) {

	pendingAssets := []domain.Asset{}

	for _, asset := range ar.Assets {
		if !asset.Sold {
			pendingAssets = append(pendingAssets, asset)
		}
	}
	return &pendingAssets, nil
}

// FindCheaperAssetPrice returns the lowest price of non sold assets
func (ar *AssetsRepositoryInMemory) FindCheaperAssetPrice() (float32, error) {
	var minimumPrice float32

	for _, asset := range ar.Assets {
		if asset.Sold == false && minimumPrice > asset.BuyPrice {
			minimumPrice = asset.BuyPrice
		}
	}

	return minimumPrice, nil
}

// GetBalance mocks the returning of balance between two dates
func (ar *AssetsRepositoryInMemory) GetBalance(startDate, endDate time.Time) (float32, error) {
	return 0, nil
}

// Create creates an asset and stores it in a data structure
func (ar *AssetsRepositoryInMemory) Create(asset *domain.Asset) error {
	ar.Assets = append(ar.Assets, *asset)
	return nil
}

// Sell updates asset state to sold and other related attributes
func (ar *AssetsRepositoryInMemory) Sell(id primitive.ObjectID, price float32) error {

	for index, asset := range ar.Assets {
		if asset.ID == id {
			ar.Assets[index].SellPrice = price
			ar.Assets[index].Sold = true
			ar.Assets[index].SellTime = time.Now()
			break
		}
	}

	return nil
}
