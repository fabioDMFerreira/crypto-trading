package decisionmaker

import "github.com/fabiodmferreira/crypto-trading/assets"

// Trader buys and sells
type Trader interface {
	Buy(amount, price float32)
	Sell(asset *assets.Asset, price float32)
}

// DecisionMaker decides to buy or sell
type DecisionMaker struct {
	trader Trader
}

func NewDecisionMaker(trader Trader) *DecisionMaker {
	return &DecisionMaker{trader}
}

// DecideToSell decides to sell
func (bs *DecisionMaker) DecideToSell(ask float32, assets []*assets.Asset, pretendedProfit float32) {
	for _, asset := range assets {
		if asset.BuyPrice+(asset.BuyPrice*pretendedProfit) < ask {
			bs.trader.Sell(asset, ask)
		}
	}
}

// DecideToBuy decides to buy
func (bs *DecisionMaker) DecideToBuy(ask float32, minimumAssetBuyPrice float32, dropToBuy float32, buyAmount float32) {
	if minimumAssetBuyPrice == 0 ||
		minimumAssetBuyPrice-(minimumAssetBuyPrice*dropToBuy) > ask {
		bs.trader.Buy(buyAmount, ask)
	}
}
