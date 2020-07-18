package decisionmaker

import (
	"time"

	"github.com/fabiodmferreira/crypto-trading/domain"
)

// DecisionMaker decides to buy or sell
type DecisionMaker struct {
	assetsRepository    domain.AssetsRepositoryReader
	options             domain.DecisionMakerOptions
	statistics          domain.Statistics
	growthStatistics    domain.Statistics
	assetsPricesService domain.AssetsPricesService

	currentPrice       float32
	lastPrice          float32
	currentChange      float32
	lastPointAddedDate time.Time
}

// NewDecisionMaker returns a new instance of DecisionMaker
func NewDecisionMaker(assetsRepository domain.AssetsRepositoryReader, options domain.DecisionMakerOptions, statistics domain.Statistics, growthStatistics domain.Statistics, assetsPricesService domain.AssetsPricesService) *DecisionMaker {
	return &DecisionMaker{assetsRepository, options, statistics, growthStatistics, assetsPricesService, 0, 0, 0, time.Time{}}
}

// NewValue adds a new price to recalculate statistics
func (dm *DecisionMaker) NewValue(price float32, date time.Time) {
	timeSinceLastPointAdded := date.Sub(dm.lastPointAddedDate).Minutes()

	change := price - dm.lastPrice

	dm.currentChange = change

	if timeSinceLastPointAdded < float64(dm.options.MinutesToCollectNewPoint) {
		return
	}

	dm.statistics.AddPoint(float64(price))
	dm.assetsPricesService.Create(date, price, "BTC")

	if dm.lastPrice > 0 {
		dm.growthStatistics.AddPoint(float64(price - dm.lastPrice))
	}

	dm.lastPrice = dm.currentPrice
	dm.currentPrice = price
	dm.lastPointAddedDate = date
}

// ShouldBuy returns true or false if it is a good time to buy
func (dm *DecisionMaker) ShouldBuy(price float32, buyTime time.Time) (bool, error) {
	if !dm.statistics.HasRequiredNumberOfPoints() {
		return false, nil
	}

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

	if dm.currentChange < dm.options.GrowthDecreaseLimit {
		return false, nil
	}

	return true, nil
}

// ShouldSell returns true or false if it is a good time to sell
func (dm *DecisionMaker) ShouldSell(asset *domain.Asset, price float32, byTime time.Time) (bool, error) {
	if !dm.statistics.HasRequiredNumberOfPoints() {
		return false, nil
	}

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

	maximumBuyAmount := dm.getMaximumBuyAmountBasedOnOptions(price)

	if float32(average-2*standardDeviation) < price {
		return maximumBuyAmount, nil
	} else if float32(average-1*standardDeviation) < price {
		return 0.8 * maximumBuyAmount, nil
	} else {
		return 0.5 * maximumBuyAmount, nil
	}
}

// getMaximumBuyAmountBasedOnOptions returns the asset amount value based on options
func (dm *DecisionMaker) getMaximumBuyAmountBasedOnOptions(price float32) float32 {
	if dm.options.MaximumFIATBuyAmount > 0 {
		return dm.options.MaximumFIATBuyAmount / price
	}

	return dm.options.MaximumBuyAmount
}

// GetState returns decision maker state
func (dm *DecisionMaker) GetState() domain.DecisionMakerState {
	average := dm.statistics.GetAverage()
	std := dm.statistics.GetStandardDeviation()

	return domain.DecisionMakerState{
		Average:             average,
		StandardDeviation:   std,
		LowerBollingerBand:  average - std,
		HigherBollingerBand: average + std,
		CurrentChange:       dm.currentChange,
	}
}
