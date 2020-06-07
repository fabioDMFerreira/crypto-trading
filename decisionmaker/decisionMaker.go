package decisionmaker

import (
	"time"

	"github.com/fabiodmferreira/crypto-trading/domain"
)

// Options are used to change DecisionMaker behaviour
type Options struct {
	MaximumBuyAmount      float32
	MinimumProfitPerSold  float32
	MinimumPriceDropToBuy float32
}

// DecisionMaker decides to buy or sell
type DecisionMaker struct {
	assetsRepository domain.AssetsRepositoryReader
	options          Options
	statistics       domain.Statistics
}

// NewDecisionMaker returns a new instance of DecisionMaker
func NewDecisionMaker(assetsRepository domain.AssetsRepositoryReader, options Options, statistics domain.Statistics) *DecisionMaker {
	return &DecisionMaker{assetsRepository, options, statistics}
}

// NewValue adds a new price to recalculate statistics
func (dm *DecisionMaker) NewValue(price float32) {
	dm.statistics.AddPoint(float64(price))
}

// ShouldBuy returns true or false if it is a good time to buy
func (dm *DecisionMaker) ShouldBuy(price float32, buyTime time.Time) (bool, error) {
	cheaperAssetPrice, err := dm.assetsRepository.FindCheaperAssetPrice()

	if err != nil {
		return false, err
	}

	if cheaperAssetPrice > 0 && cheaperAssetPrice-(cheaperAssetPrice*dm.options.MinimumPriceDropToBuy) < price {
		return false, nil
	}

	if float32(dm.statistics.GetAverage()) < price {
		return false, nil
	}

	return true, nil
}

// ShouldSell returns true or false if it is a good time to sell
func (dm *DecisionMaker) ShouldSell(asset *domain.Asset, price float32, byTime time.Time) (bool, error) {
	if asset.BuyPrice+(asset.BuyPrice*dm.options.MinimumProfitPerSold) > price {
		return false, nil
	}

	if float32(dm.statistics.GetAverage()) > price {
		return false, nil
	}

	return true, nil
}

// HowMuchAmountShouldBuy returns the amount of asset that should be bought
func (dm *DecisionMaker) HowMuchAmountShouldBuy(price float32) (float32, error) {

	standardDeviation := dm.statistics.GetStandardDeviation()
	average := dm.statistics.GetAverage()

	if float32(average-2*standardDeviation) < price {
		return dm.options.MaximumBuyAmount, nil
	} else if float32(average-1*standardDeviation) < price {
		return 0.8 * dm.options.MaximumBuyAmount, nil
	} else {
		return 0.1 * dm.options.MaximumBuyAmount, nil
	}
}
