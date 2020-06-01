package trader

import (
	"time"

	"github.com/fabiodmferreira/crypto-trading/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Trader struct {
	assetsRepository domain.AssetsRepository
	accountService   domain.AccountService
	broker           domain.Broker
}

func NewTrader(assetsRepository domain.AssetsRepository, accountService domain.AccountService, broker domain.Broker) *Trader {
	return &Trader{
		assetsRepository,
		accountService,
		broker,
	}
}

func (t *Trader) Sell(asset *domain.Asset, price float32) error {
	err := t.assetsRepository.Sell(asset.ID, price)

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

func (t *Trader) Buy(amount, price float32, buyTime time.Time) error {
	amountToWithdraw := amount * price
	err := t.accountService.Withdraw(amountToWithdraw)

	if err != nil {
		return err
	}

	asset := &domain.Asset{ID: primitive.NewObjectID(), Amount: amount, BuyPrice: price, BuyTime: buyTime}
	err = t.assetsRepository.Create(asset)

	if err != nil {
		return err
	}

	err = t.broker.AddBuyOrder(asset.Amount, price)
	if err != nil {
		return err
	}

	return nil
}
