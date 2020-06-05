package decisionmaker

import (
	"time"

	"github.com/fabiodmferreira/crypto-trading/domain"
)

type DecisionMakerOptions struct {
	MaximumBuyAmount       float32
	PretendedProfitPerSold float32
	PriceDropToBuy         float32
}

// DecisionMaker decides to buy or sell
type DecisionMaker struct {
	assetsRepository domain.AssetsRepositoryReader
	options          DecisionMakerOptions
}

func NewDecisionMaker(assetsRepository domain.AssetsRepositoryReader, options DecisionMakerOptions) *DecisionMaker {
	return &DecisionMaker{assetsRepository, options}
}

// New Value is a mock
func (db *DecisionMaker) NewValue(price float32) {}

func (dm *DecisionMaker) ShouldBuy(price float32, buyTime time.Time) (bool, error) {
	cheaperAssetPrice, err := dm.assetsRepository.FindCheaperAssetPrice()

	if err != nil {
		return false, err
	}

	return cheaperAssetPrice == 0 || cheaperAssetPrice-(cheaperAssetPrice*dm.options.PriceDropToBuy) > price, nil
}

func (dm *DecisionMaker) ShouldSell(asset *domain.Asset, price float32, byTime time.Time) (bool, error) {
	return asset.BuyPrice+(asset.BuyPrice*dm.options.PretendedProfitPerSold) < price, nil
}

func (dm *DecisionMaker) HowMuchAmountShouldBuy(price float32) (float32, error) {
	return dm.options.MaximumBuyAmount, nil
}
