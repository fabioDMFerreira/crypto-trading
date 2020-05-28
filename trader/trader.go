package trader

import (
	"fmt"
	"log"
	"time"

	"github.com/fabiodmferreira/crypto-trading/assets"
	"github.com/fabiodmferreira/crypto-trading/eventlogs"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DBTrader struct {
	assetsRepository    *assets.AssetsRepository
	eventLogsRepository *eventlogs.EventLogsRepository
}

func NewDBTrader(assetsRepository *assets.AssetsRepository, eventLogsRepository *eventlogs.EventLogsRepository) *DBTrader {
	return &DBTrader{
		assetsRepository,
		eventLogsRepository,
	}
}

func (t *DBTrader) Sell(asset *assets.Asset, price float32) {
	err := t.assetsRepository.Sell(asset.ID, price)

	if err != nil {
		log.Fatal(err)
	}

	message := fmt.Sprintf("Asset sold: {Price: %v Amount: %v Value: %v}", price, asset.Amount, price*asset.Amount)
	err = t.eventLogsRepository.Create("sell", message)
	if err != nil {
		log.Fatal(err)
	}
}

func (t *DBTrader) Buy(amount, price float32) {
	asset := &assets.Asset{ID: primitive.NewObjectID(), Amount: amount, BuyPrice: price, BuyTime: time.Now()}
	err := t.assetsRepository.Create(asset)

	if err != nil {
		log.Fatal(err)
	}

	message := fmt.Sprintf("Asset bought: {Price: %v Amount: %v Value: %v}", price, amount, price*amount)
	err = t.eventLogsRepository.Create("buy", message)
	if err != nil {
		log.Fatal(err)
	}
}
