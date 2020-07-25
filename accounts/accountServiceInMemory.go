package accounts

import (
	"errors"

	"github.com/fabiodmferreira/crypto-trading/domain"
)

// AccountServiceInMemory emulates an account service by saving data in memory
type AccountServiceInMemory struct {
	Amount           float32
	withdraws        int
	deposits         int
	assetsRepository domain.AssetsRepository
}

// NewAccountServiceInMemory returns an instance of AccountServiceInMemory
func NewAccountServiceInMemory(initialAmount float32, assetsRepository domain.AssetsRepository) *AccountServiceInMemory {
	return &AccountServiceInMemory{initialAmount, 0, 0, assetsRepository}
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
	return a.assetsRepository.FindPendingAssets()
}

func (a *AccountServiceInMemory) FindAllAssets() (*[]domain.Asset, error) {
	return a.assetsRepository.FindAll()
}
