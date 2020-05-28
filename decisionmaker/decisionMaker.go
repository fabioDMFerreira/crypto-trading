package decisionmaker

import (
	"log"

	"github.com/fabiodmferreira/crypto-trading/assets"
)

// Trader buys and sells
type Trader interface {
	Buy(amount, price float32)
	Sell(asset *assets.Asset, price float32)
}

type Account interface {
	Withdraw(amount float32) error
	Deposit(amount float32) error
}

// DecisionMaker decides to buy or sell
type DecisionMaker struct {
	trader  Trader
	account Account
}

func NewDecisionMaker(trader Trader, account Account) *DecisionMaker {
	return &DecisionMaker{trader, account}
}

// DecideToSell decides to sell
func (dm *DecisionMaker) DecideToSell(ask float32, assets []*assets.Asset, pretendedProfit float32) {
	for _, asset := range assets {
		if asset.BuyPrice+(asset.BuyPrice*pretendedProfit) < ask {

			dm.trader.Sell(asset, ask)

			amountToDeposit := asset.Amount * ask
			dm.account.Deposit(amountToDeposit)
		}
	}
}

// DecideToBuy decides to buy
func (dm *DecisionMaker) DecideToBuy(ask float32, minimumAssetBuyPrice float32, dropToBuy float32, buyAmount float32) {
	if minimumAssetBuyPrice == 0 ||
		minimumAssetBuyPrice-(minimumAssetBuyPrice*dropToBuy) > ask {
		amountToWithdraw := buyAmount * ask
		err := dm.account.Withdraw(amountToWithdraw)

		if err == nil {
			dm.trader.Buy(buyAmount, ask)
		} else {
			log.Fatal(err)
		}
	}
}
