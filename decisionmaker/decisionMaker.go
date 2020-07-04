package decisionmaker

import (
	"fmt"
	"time"

	"github.com/fabiodmferreira/crypto-trading/domain"
)

// DecisionMaker decides to buy or sell
type DecisionMaker struct {
	assetsRepository domain.AssetsRepositoryReader
	options          domain.DecisionMakerOptions
	statistics       domain.Statistics

	lastSell      int
	sellCumulator int
	sellsLimit    int
}

// NewDecisionMaker returns a new instance of DecisionMaker
func NewDecisionMaker(assetsRepository domain.AssetsRepositoryReader, options domain.DecisionMakerOptions, statistics domain.Statistics) *DecisionMaker {
	return &DecisionMaker{assetsRepository, options, statistics, -1, 0, 500}
}

// NewValue adds a new price to recalculate statistics
func (dm *DecisionMaker) NewValue(price float32) {
	dm.statistics.AddPoint(float64(price))
}

// ShouldBuy returns true or false if it is a good time to buy
func (dm *DecisionMaker) ShouldBuy(price float32, buyTime time.Time) (bool, error) {
	assetWithCloserPrice, err := dm.assetsRepository.CheckAssetWithCloserPriceExists(price, 0.02)

	if err != nil {
		return false, err
	}

	if assetWithCloserPrice {
		return false, nil
	}

	if float32(dm.statistics.GetAverage()-dm.statistics.GetStandardDeviation()) < price {
		return false, nil
	}

	if dm.lastSell < 0 {
		dm.lastSell = 0
	} else if dm.lastSell < dm.sellsLimit {
		dm.sellCumulator++
		dm.lastSell++
		return false, nil
	} else {
		dm.lastSell = 0
	}

	return true, nil
}

// ShouldSell returns true or false if it is a good time to sell
func (dm *DecisionMaker) ShouldSell(asset *domain.Asset, price float32, byTime time.Time) (bool, error) {
	if asset.BuyPrice+(asset.BuyPrice*dm.options.MinimumProfitPerSold) > price {
		return false, nil
	}

	if float32(dm.statistics.GetAverage()+dm.statistics.GetStandardDeviation()) > price {
		return false, nil
	}

	return true, nil
}

// HowMuchAmountShouldBuy returns the amount of asset that should be bought
func (dm *DecisionMaker) HowMuchAmountShouldBuy(price float32) (float32, error) {

	standardDeviation := dm.statistics.GetStandardDeviation()
	average := dm.statistics.GetAverage()

	if dm.sellCumulator > 0 {
		factor := dm.sellCumulator / dm.sellsLimit
		value := 2 * dm.options.MaximumBuyAmount
		fmt.Printf("%v %v %v %v\n", dm.sellCumulator, dm.sellsLimit, factor, value)
		dm.sellCumulator = 0
		return value, nil
	}

	if float32(average-2*standardDeviation) < price {
		return dm.options.MaximumBuyAmount, nil
	} else if float32(average-1*standardDeviation) < price {
		return 0.8 * dm.options.MaximumBuyAmount, nil
	} else {
		return 0.5 * dm.options.MaximumBuyAmount, nil
	}
}
