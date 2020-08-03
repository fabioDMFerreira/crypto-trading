package accounts

import (
	"time"

	"github.com/fabiodmferreira/crypto-trading/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccountService struct {
	ID               primitive.ObjectID
	repository       *Repository
	assetsRepository domain.AssetsRepository
}

func NewAccountService(ID primitive.ObjectID, repository *Repository, assetsRepository domain.AssetsRepository) *AccountService {
	return &AccountService{ID, repository, assetsRepository}
}

func (a *AccountService) Withdraw(amount float32) error {
	return a.repository.Withdraw(a.ID, amount)
}

func (a *AccountService) Deposit(amount float32) error {
	return a.repository.Deposit(a.ID, amount)
}

func (a *AccountService) GetAmount() (float32, error) {
	account, err := a.repository.FindById(a.ID)

	if err != nil {
		return 0, err
	}

	return account.Amount, nil
}

func (a *AccountService) FindPendingAssets() (*[]domain.Asset, error) {
	return a.assetsRepository.FindPendingAssets(a.ID)
}

func (a *AccountService) FindAllAssets() (*[]domain.Asset, error) {
	return a.assetsRepository.FindAll(a.ID)
}

func (a *AccountService) CreateAsset(amount, price float32, time time.Time) (*domain.Asset, error) {
	asset := &domain.Asset{ID: primitive.NewObjectID(), Amount: amount, BuyPrice: price, BuyTime: time, AccountID: a.ID}

	err := a.assetsRepository.Create(asset)

	return asset, err
}

func (a *AccountService) SellAsset(assetID primitive.ObjectID, price float32, time time.Time) error {
	return a.assetsRepository.Sell(assetID, price, time)
}

func (a *AccountService) GetBalance(startDate, endDate time.Time) (float32, error) {
	return a.assetsRepository.GetBalance(a.ID, startDate, endDate)
}

func (a *AccountService) CheckAssetWithCloserPriceExists(price, limit float32) (bool, error) {
	return a.assetsRepository.CheckAssetWithCloserPriceExists(a.ID, price, limit)
}
