package decisionmaker

import (
	"time"

	"github.com/fabiodmferreira/crypto-trading/domain"
	"github.com/fabiodmferreira/crypto-trading/statistics"
)

type DecisionMakerOptions1 struct {
	MaximumBuyAmount      float32
	MinimumProfitPerSold  float32
	MinimumPriceDropToBuy float32
}

// DecisionMaker decides to buy or sell
type DecisionMaker1 struct {
	assetsRepository domain.AssetsRepositoryReader
	options          DecisionMakerOptions1
	statistics       *statistics.Statistics
}

func NewDecisionMaker1(assetsRepository domain.AssetsRepositoryReader, options DecisionMakerOptions1, statistics *statistics.Statistics) *DecisionMaker1 {
	return &DecisionMaker1{assetsRepository, options, statistics}
}

func (dm *DecisionMaker1) NewValue(price float32) {
	dm.statistics.AddPoint(float64(price))
}

func (dm *DecisionMaker1) ShouldBuy(price float32, buyTime time.Time) (bool, error) {
	cheaperAssetPrice, err := dm.assetsRepository.FindCheaperAssetPrice()

	if err != nil {
		return false, err
	}

	if cheaperAssetPrice > 0 && cheaperAssetPrice-(cheaperAssetPrice*dm.options.MinimumPriceDropToBuy) < price {
		return false, nil
	}

	// if float32(dm.statistics.Average) < price {
	// 	return false, nil
	// }

	if float32(dm.statistics.MACD.GetLastHistogramPoint()) >= -5 {
		return false, nil
	}

	return true, nil
}

func (dm *DecisionMaker1) ShouldSell(asset *domain.Asset, price float32, byTime time.Time) (bool, error) {
	if asset.BuyPrice+(asset.BuyPrice*dm.options.MinimumProfitPerSold) > price {
		return false, nil
	}

	// if float32(dm.statistics.Average) > price {
	// 	return false, nil
	// }

	// if float32(dm.statistics.MACD.GetLastHistogramPoint()) <= 5 {
	// 	return false, nil
	// }

	return true, nil
}

func (dm *DecisionMaker1) HowMuchAmountShouldBuy(price float32) (float32, error) {
	return dm.options.MaximumBuyAmount, nil
}
