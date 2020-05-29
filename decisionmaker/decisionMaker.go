package decisionmaker

import (
	"log"
	"time"

	"github.com/fabiodmferreira/crypto-trading/assets"
	"github.com/fabiodmferreira/crypto-trading/domain"
)

type DecisionMakerOptions struct {
	MaximumBuyAmount       float32
	PretendedProfitPerSold float32
	PriceDropToBuy         float32
}

// DecisionMaker decides to buy or sell
type DecisionMaker struct {
	trader           domain.Trader
	account          domain.Account
	assetsRepository domain.AssetsRepositoryReader
	options          DecisionMakerOptions
}

func NewDecisionMaker(trader domain.Trader, account domain.Account, assetsRepository domain.AssetsRepositoryReader, options DecisionMakerOptions) *DecisionMaker {
	return &DecisionMaker{trader, account, assetsRepository, options}
}

func (dm *DecisionMaker) MakeDecisions(price float32, buyTime time.Time) {
	assets, err := dm.assetsRepository.FindAll()
	if err == nil {
		dm.DecideToSell(price, assets, dm.options.PretendedProfitPerSold)
	}

	cheaperAssetPrice, err := dm.assetsRepository.FindCheaperAssetPrice()
	if err == nil {
		dm.DecideToBuy(price, cheaperAssetPrice, dm.options.PriceDropToBuy, dm.options.MaximumBuyAmount, buyTime)
	} else {
		log.Fatal(err)
	}

}

// DecideToSell decides to sell
func (dm *DecisionMaker) DecideToSell(ask float32, assets *[]assets.Asset, pretendedProfit float32) {
	for _, asset := range *assets {
		if asset.BuyPrice+(asset.BuyPrice*pretendedProfit) < ask {

			dm.trader.Sell(&asset, ask)

			amountToDeposit := asset.Amount * ask
			dm.account.Deposit(amountToDeposit)
		}
	}
}

// DecideToBuy decides to buy
func (dm *DecisionMaker) DecideToBuy(ask float32, minimumAssetBuyPrice float32, dropToBuy float32, buyAmount float32, buyTime time.Time) {
	if minimumAssetBuyPrice == 0 ||
		minimumAssetBuyPrice-(minimumAssetBuyPrice*dropToBuy) > ask {
		amountToWithdraw := buyAmount * ask
		err := dm.account.Withdraw(amountToWithdraw)

		if err == nil {
			dm.trader.Buy(buyAmount, ask, buyTime)
		} else {
			log.Fatal(err)
		}
	}
}
