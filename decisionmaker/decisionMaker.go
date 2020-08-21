package decisionmaker

import (
	"time"

	"github.com/fabiodmferreira/crypto-trading/domain"
)

// DecisionMaker decides to buy or sell
type DecisionMaker struct {
	account                domain.AccountService
	options                domain.DecisionMakerOptions
	statistics             domain.Statistics
	growthStatistics       domain.Statistics
	accelerationStatistics domain.Statistics

	currentPrice        float32
	lastPrice           float32
	currentChange       float32
	lastChange          float32
	currentAcceleration float32
	lastAcceleration    float32
}

// NewDecisionMaker returns a new instance of DecisionMaker
func NewDecisionMaker(
	account domain.AccountService,
	options domain.DecisionMakerOptions,
	statistics domain.Statistics,
	growthStatistics domain.Statistics,
	accelerationStatistics domain.Statistics,
) *DecisionMaker {
	return &DecisionMaker{account, options, statistics, growthStatistics, accelerationStatistics, 0, 0, 0, 0, 0, 0}
}

// NewValue adds a new price to recalculate statistics
func (dm *DecisionMaker) NewValue(price float32, date time.Time) {
	dm.statistics.AddPoint(float64(price))

	if dm.lastPrice > 0 {
		change := price - dm.lastPrice

		dm.lastChange = dm.currentChange
		dm.currentChange = change

		dm.growthStatistics.AddPoint(float64(change))
	}

	if dm.lastChange != 0 {
		acceleration := dm.currentChange - dm.lastChange

		dm.accelerationStatistics.AddPoint(float64(acceleration))

		dm.lastAcceleration = dm.currentAcceleration
		dm.currentAcceleration = acceleration
	}

	dm.lastPrice = dm.currentPrice
	dm.currentPrice = price
}

// ShouldBuy returns true or false if it is a good time to buy
func (dm *DecisionMaker) ShouldBuy(price float32, buyTime time.Time) (bool, error) {
	if !dm.statistics.HasRequiredNumberOfPoints() {
		return false, nil
	}

	assetWithCloserPrice, err := dm.account.CheckAssetWithCloserPriceExists(price, 0.02)

	if err != nil {
		return false, err
	}

	if assetWithCloserPrice ||
		float32(dm.statistics.GetAverage()-dm.statistics.GetStandardDeviation()) < price ||
		dm.currentChange < dm.options.GrowthDecreaseLimit ||
		((float32(dm.accelerationStatistics.GetAverage()) > 0) && (float32(dm.growthStatistics.GetAverage()) < 0)) {
		return false, nil
	}

	return true, nil
}

// ShouldSell returns true or false if it is a good time to sell
func (dm *DecisionMaker) ShouldSell(asset *domain.Asset, price float32, byTime time.Time) (bool, error) {
	if !dm.statistics.HasRequiredNumberOfPoints() ||
		asset.BuyPrice+(asset.BuyPrice*dm.options.MinimumProfitPerSold) > price ||
		float32(dm.statistics.GetAverage()+dm.statistics.GetStandardDeviation()) > price ||
		((float32(dm.accelerationStatistics.GetAverage()) > 0) && (float32(dm.growthStatistics.GetAverage()) > 0)) {
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
	priceAverage := dm.statistics.GetAverage()
	priceStd := dm.statistics.GetStandardDeviation()

	changeAverage := dm.growthStatistics.GetAverage()
	changeStd := dm.growthStatistics.GetStandardDeviation()

	accelerationAverage := dm.accelerationStatistics.GetAverage()
	accelerationStd := dm.accelerationStatistics.GetStandardDeviation()

	return domain.DecisionMakerState{
		Price:                  dm.currentPrice,
		PriceAverage:           priceAverage,
		PriceStandardDeviation: priceStd,
		PriceUpperLimit:        priceAverage + priceStd,
		PriceLowerLimit:        priceAverage - priceStd,

		Change:                  dm.currentChange,
		ChangeAverage:           changeAverage,
		ChangeStandardDeviation: changeStd,
		ChangeUpperLimit:        changeAverage + changeStd,
		ChangeLowerLimit:        changeAverage - changeStd,

		Acceleration:                  dm.currentAcceleration,
		AccelerationAverage:           accelerationAverage,
		AccelerationStandardDeviation: accelerationStd,
		AccelerationUpperLimit:        accelerationAverage + accelerationStd,
		AccelerationLowerLimit:        accelerationAverage - accelerationStd,
	}
}
