package decisionmaker

import (
	"github.com/fabiodmferreira/crypto-trading/domain"
	"github.com/fabiodmferreira/crypto-trading/indicators"
)

type SellStrategy struct {
	priceStats  *indicators.PriceIndicator
	volumeStats *indicators.VolumeIndicator
	account     domain.AccountService
	options     domain.DecisionMakerOptions
}

func NewSellStrategy(
	priceStats *indicators.PriceIndicator,
	volumeStats *indicators.VolumeIndicator,
	accountService domain.AccountService,
	options domain.DecisionMakerOptions,
) *SellStrategy {
	return &SellStrategy{
		priceStats:  priceStats,
		volumeStats: volumeStats,
		account:     accountService,
		options:     options,
	}
}

func (s *SellStrategy) Execute() (bool, float32, error) {
	priceStatsState := s.priceStats.GetState().(*indicators.MetricStatisticsIndicatorState)
	volumeStatsState := s.volumeStats.GetState().(*indicators.MetricStatisticsIndicatorState)

	assetWithCloserPrice, err := s.account.CheckAssetWithCloserPriceExists(float32(priceStatsState.Value), s.options.MinimumPriceDropToBuy)

	if err != nil || assetWithCloserPrice {
		return false, 0, err
	}

	if priceStatsState.HasRequiredPoints &&
		volumeStatsState.HasRequiredPoints &&
		priceStatsState.Average+priceStatsState.StandardDeviation < priceStatsState.Value &&
		volumeStatsState.Value < volumeStatsState.Average {
		return true, 100, nil
	}

	return false, 0, nil
}
