package accounts

import (
	"github.com/fabiodmferreira/crypto-trading/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccountService struct {
	ID               primitive.ObjectID
	repository       *Repository
	assetsRepository domain.AssetsRepositoryReader
}

func NewAccountService(ID primitive.ObjectID, repository *Repository, assetsRepository domain.AssetsRepositoryReader) *AccountService {
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
	return a.assetsRepository.FindPendingAssets()
}

func (a *AccountService) FindAllAssets() (*[]domain.Asset, error) {
	return a.assetsRepository.FindAll()
}
