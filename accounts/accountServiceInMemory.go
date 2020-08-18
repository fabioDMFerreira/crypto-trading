package accounts

import (
	"errors"
	"time"

	"github.com/fabiodmferreira/crypto-trading/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AccountServiceInMemory emulates an account service by saving data in memory
type AccountServiceInMemory struct {
	Amount           float32
	withdraws        int
	deposits         int
	assetsRepository domain.AssetsRepository
	ID               string
}

// NewAccountServiceInMemory returns an instance of AccountServiceInMemory
func NewAccountServiceInMemory(initialAmount float32, assetsRepository domain.AssetsRepository) *AccountServiceInMemory {
	return &AccountServiceInMemory{initialAmount, 0, 0, assetsRepository, primitive.NewObjectID().Hex()}
}

// Deposit increases account amount
func (a *AccountServiceInMemory) Deposit(amount float32) error {
	a.deposits++
	a.Amount += amount
	return nil
}

// Withdraw decreases account amount
func (a *AccountServiceInMemory) Withdraw(amount float32) error {
	if amount > a.Amount {
		return errors.New("Insufficient Funds")
	}

	a.Amount -= amount
	a.withdraws++

	return nil
}

// GetAmount returns amount value
func (a *AccountServiceInMemory) GetAmount() (float32, error) {
	return a.Amount, nil
}

func (a *AccountServiceInMemory) FindPendingAssets() (*[]domain.Asset, error) {
	return a.assetsRepository.FindPendingAssets(a.ID)
}

func (a *AccountServiceInMemory) FindAllAssets() (*[]domain.Asset, error) {
	return a.assetsRepository.FindAll(a.ID)
}

func (a *AccountServiceInMemory) CreateAsset(amount, price float32, time time.Time) (*domain.Asset, error) {
	asset := &domain.Asset{ID: primitive.NewObjectID(), Amount: amount, BuyPrice: price, BuyTime: time}

	err := a.assetsRepository.Create(asset)

	return asset, err
}

func (a *AccountServiceInMemory) SellAsset(assetID string, price float32, time time.Time) error {
	return a.assetsRepository.Sell(assetID, price, time)
}

func (a *AccountServiceInMemory) GetBalance(startDate, endDate time.Time) (float32, error) {
	return a.assetsRepository.GetBalance(a.ID, startDate, endDate)
}

func (a *AccountServiceInMemory) CheckAssetWithCloserPriceExists(price, limit float32) (bool, error) {
	return a.assetsRepository.CheckAssetWithCloserPriceExists(a.ID, price, limit)
}
