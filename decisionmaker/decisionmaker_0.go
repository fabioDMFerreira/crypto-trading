package decisionmaker

import (
	"time"

	"github.com/fabiodmferreira/crypto-trading/domain"
)

type DecisionMaker0Options struct {
	MaximumBuyAmount       float32
	PretendedProfitPerSold float32
	PriceDropToBuy         float32
}

// DecisionMaker0 decides to buy or sell
type DecisionMaker0 struct {
	assetsRepository domain.AssetsRepositoryReader
	options          DecisionMaker0Options
}

func NewDecisionMaker0(assetsRepository domain.AssetsRepositoryReader, options DecisionMaker0Options) *DecisionMaker0 {
	return &DecisionMaker0{assetsRepository, options}
}

// New Value is a mock
func (db *DecisionMaker0) NewValue(price float32) {}

func (dm *DecisionMaker0) ShouldBuy(price float32, buyTime time.Time) (bool, error) {
	cheaperAssetPrice, err := dm.assetsRepository.FindCheaperAssetPrice()

	if err != nil {
		return false, err
	}

	return cheaperAssetPrice == 0 || cheaperAssetPrice-(cheaperAssetPrice*dm.options.PriceDropToBuy) > price, nil
}

func (dm *DecisionMaker0) ShouldSell(asset *domain.Asset, price float32, byTime time.Time) (bool, error) {
	return asset.BuyPrice+(asset.BuyPrice*dm.options.PretendedProfitPerSold) < price, nil
}

func (dm *DecisionMaker0) HowMuchAmountShouldBuy(price float32) (float32, error) {
	return dm.options.MaximumBuyAmount, nil
}
