package domain

import "time"

// DecisionMakerOptions are used to change DecisionMaker behaviour
type DecisionMakerOptions struct {
	MaximumBuyAmount      float32 `json:"maximumBuyAmount"`
	MaximumFIATBuyAmount  float32 `json:"maximumFIATBuyAmount"`
	MinimumProfitPerSold  float32 `json:"minimumProfitPerSold"`
	MinimumPriceDropToBuy float32 `json:"minimumPriceDropToBuy"`
	GrowthDecreaseLimit   float32 `json:"growthDecreaseLimit"`
	GrowthIncreaseLimit   float32 `json:"growthIncreaseLimit"`
}

// DecisionMakerState represents the decision maker state
type DecisionMakerState struct {
	Average             float64 `json:"average"`
	StandardDeviation   float64 `json:"standardDeviation"`
	LowerBollingerBand  float64 `json:"lowerBollingerBand"`
	HigherBollingerBand float64 `json:"higherBollingerBand"`
	CurrentChange       float32 `json:"currentChange"`
}

// DecisionMaker makes decisions to buy or sell assets
type DecisionMaker interface {
	NewValue(price float32, date time.Time)
	ShouldBuy(price float32, buyTime time.Time) (bool, error)
	ShouldSell(asset *Asset, price float32, buyTime time.Time) (bool, error)
	HowMuchAmountShouldBuy(price float32) (float32, error)
	GetState() DecisionMakerState
}
