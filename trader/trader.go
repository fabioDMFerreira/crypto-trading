package trader

import (
	"time"

	"github.com/fabiodmferreira/crypto-trading/domain"
)

// Trader execute operations to buy and sell assets
type Trader struct {
	broker domain.Broker
}

// NewTrader returns a Trader instance
func NewTrader(broker domain.Broker) *Trader {
	return &Trader{
		broker,
	}
}

// Sell requests broker to sell an asset
func (t *Trader) Sell(asset *domain.Asset, price float32, sellTime time.Time) error {
	return t.broker.AddSellOrder(asset.Amount, price)
}

// Buy requests broker to buy an asset
func (t *Trader) Buy(amount, price float32, buyTime time.Time) error {
	return t.broker.AddBuyOrder(amount, price)
}
