package accounts

import (
	"fmt"
	"time"

	"github.com/fabiodmferreira/crypto-trading/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AccountService interacts with accounts and assets repositories
type AccountService struct {
	ID               string
	repository       *Repository
	assetsRepository domain.AssetsRepository
}

// NewAccountService returns an instance of account service
func NewAccountService(ID string, repository *Repository, assetsRepository domain.AssetsRepository) (*AccountService, error) {
	_, err := repository.FindById(ID)

	if err != nil {
		return nil, fmt.Errorf("Not able to get account with id %v due to %v", ID, err)
	}

	return &AccountService{ID, repository, assetsRepository}, nil
}

// Withdraw decrements an amount from an account
func (a *AccountService) Withdraw(amount float32) error {
	return a.repository.Withdraw(a.ID, amount)
}

// Deposit increments an amount to an account
func (a *AccountService) Deposit(amount float32) error {
	return a.repository.Deposit(a.ID, amount)
}

// GetAmount returns the amount hold by the account
func (a *AccountService) GetAmount() (float32, error) {
	account, err := a.repository.FindById(a.ID)

	if err != nil {
		return 0, err
	}

	return account.Amount, nil
}

// FindPendingAssets returns account assets awaiting to be sold
func (a *AccountService) FindPendingAssets() (*[]domain.Asset, error) {
	return a.assetsRepository.FindPendingAssets(a.ID)
}

// FindAllAssets returns all assets hold by the account
func (a *AccountService) FindAllAssets() (*[]domain.Asset, error) {
	return a.assetsRepository.FindAll(a.ID)
}

// CreateAsset creates an asset hold by the account
func (a *AccountService) CreateAsset(amount, price float32, time time.Time) (*domain.Asset, error) {
	accountOID, err := primitive.ObjectIDFromHex(a.ID)

	if err != nil {
		return nil, err
	}

	asset := &domain.Asset{ID: primitive.NewObjectID(), Amount: amount, BuyPrice: price, BuyTime: time, AccountID: accountOID}

	err = a.assetsRepository.Create(asset)

	return asset, err
}

// SellAsset updates asset status to sold
func (a *AccountService) SellAsset(assetID string, price float32, time time.Time) error {
	return a.assetsRepository.Sell(assetID, price, time)
}

// GetBalance returns the balance between two dates
func (a *AccountService) GetBalance(startDate, endDate time.Time) (float32, error) {
	return a.assetsRepository.GetBalance(a.ID, startDate, endDate)
}

// CheckAssetWithCloserPriceExists verifies whether account already has asset with a price close to the one passed by argument
func (a *AccountService) CheckAssetWithCloserPriceExists(price, limit float32) (bool, error) {
	return a.assetsRepository.CheckAssetWithCloserPriceExists(a.ID, price, limit)
}
