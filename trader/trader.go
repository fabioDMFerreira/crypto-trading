package trader

import (
	"fmt"
	"log"
	"time"

	"github.com/fabiodmferreira/crypto-trading/assets"
	"github.com/fabiodmferreira/crypto-trading/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Trader struct {
	assetsRepository    domain.AssetsRepository
	account             domain.AccountService
	eventLogsRepository domain.EventsLog
	broker              domain.Broker
}

func NewTrader(assetsRepository domain.AssetsRepository, account domain.AccountService, eventLogsRepository domain.EventsLog, broker domain.Broker) *Trader {
	return &Trader{
		assetsRepository,
		account,
		eventLogsRepository,
		broker,
	}
}

func (t *Trader) Sell(asset *assets.Asset, price float32) {
	err := t.assetsRepository.Sell(asset.ID, price)

	if err != nil {
		log.Fatal(err)
	}

	err = t.broker.AddSellOrder(asset.Amount, price)
	if err != nil {
		log.Fatal(err)
	}

	message := fmt.Sprintf("Asset sold: {Price: %v Amount: %v Value: %v}", price, asset.Amount, price*asset.Amount)
	err = t.eventLogsRepository.Create("sell", message)
	if err != nil {
		log.Fatal(err)
	}

	amountToDeposit := asset.Amount * price
	t.account.Deposit(amountToDeposit)
}

func (t *Trader) Buy(amount, price float32, buyTime time.Time) {
	amountToWithdraw := amount * price
	err := t.account.Withdraw(amountToWithdraw)

	if err != nil {
		log.Fatal(err)
	}

	asset := &assets.Asset{ID: primitive.NewObjectID(), Amount: amount, BuyPrice: price, BuyTime: buyTime}
	err = t.assetsRepository.Create(asset)

	if err != nil {
		log.Fatal(err)
	}

	err = t.broker.AddBuyOrder(asset.Amount, price)
	if err != nil {
		log.Fatal(err)
	}

	message := fmt.Sprintf("Asset bought: {Price: %v Amount: %v Value: %v}", price, amount, price*amount)
	err = t.eventLogsRepository.Create("buy", message)
	if err != nil {
		log.Fatal(err)
	}
}
