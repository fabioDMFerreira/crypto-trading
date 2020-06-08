package accounts

import "errors"

// AccountServiceInMemory emulates an account service by saving data in memory
type AccountServiceInMemory struct {
	Amount    float32
	withdraws int
	deposits  int
}

// NewAccountServiceInMemory returns an instance of AccountServiceInMemory
func NewAccountServiceInMemory(initialAmount float32) *AccountServiceInMemory {
	return &AccountServiceInMemory{initialAmount, 0, 0}
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
