package trader

import (
	"time"

	"github.com/fabiodmferreira/crypto-trading/domain"
)

// Trader execute operations to buy and sell assets
type Trader struct {
	accountService domain.AccountService
	broker         domain.Broker
}

// NewTrader returns a Trader instance
func NewTrader(accountService domain.AccountService, broker domain.Broker) *Trader {
	return &Trader{
		accountService,
		broker,
	}
}

// Sell updates asset status to sold, requests broker to sell an asset and updates account ammount
func (t *Trader) Sell(asset *domain.Asset, price float32, sellTime time.Time) error {
	err := t.accountService.SellAsset(asset.ID, price, sellTime)

	if err != nil {
		return err
	}

	err = t.broker.AddSellOrder(asset.Amount, price)
	if err != nil {
		return err
	}

	amountToDeposit := asset.Amount * price
	err = t.accountService.Deposit(amountToDeposit)

	if err != nil {
		return err
	}

	return nil
}

// Buy updates account ammount, creates an asset and requests broker to buy an asset
func (t *Trader) Buy(amount, price float32, buyTime time.Time) error {
	amountToWithdraw := amount * price
	err := t.accountService.Withdraw(amountToWithdraw)

	if err != nil {
		return err
	}

	asset, err := t.accountService.CreateAsset(amount, price, buyTime)

	if err != nil {
		return err
	}

	err = t.broker.AddBuyOrder(asset.Amount, price)
	if err != nil {
		return err
	}

	return nil
}
