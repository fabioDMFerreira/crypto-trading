package domain

import "time"

// DecisionMakerOptions are used to change DecisionMaker behaviour
type DecisionMakerOptions struct {
	MaximumBuyAmount      float32 `bson:"maximumBuyAmount,truncate" json:"maximumBuyAmount"`
	MaximumFIATBuyAmount  float32 `bson:"maximumFIATBuyAmount,truncate" json:"maximumFIATBuyAmount"`
	MinimumProfitPerSold  float32 `bson:"minimumProfitPerSold,truncate" json:"minimumProfitPerSold"`
	MinimumPriceDropToBuy float32 `bson:"minimumPriceDropToBuy,truncate" json:"minimumPriceDropToBuy"`
	GrowthDecreaseLimit   float32 `bson:"growthDecreaseLimit,truncate" json:"growthDecreaseLimit"`
	GrowthIncreaseLimit   float32 `bson:"growthIncreaseLimit,truncate" json:"growthIncreaseLimit"`
}

// DecisionMakerState represents the decision maker state
type DecisionMakerState struct {
	Average             float64 `bson:"average,truncate" json:"average"`
	StandardDeviation   float64 `bson:"standardDeviation,truncate" json:"standardDeviation"`
	LowerBollingerBand  float64 `bson:"lowerBollingerBand,truncate" json:"lowerBollingerBand"`
	HigherBollingerBand float64 `bson:"higherBollingerBand,truncate" json:"higherBollingerBand"`
	CurrentPrice        float32 `bson:"currentPrice,truncate" json:"currentPrice"`
	CurrentChange       float32 `bson:"currentChange,truncate" json:"currentChange"`
}

// DecisionMaker makes decisions to buy or sell assets
type DecisionMaker interface {
	NewValue(price float32, date time.Time)
	ShouldBuy(price float32, buyTime time.Time) (bool, error)
	ShouldSell(asset *Asset, price float32, buyTime time.Time) (bool, error)
	HowMuchAmountShouldBuy(price float32) (float32, error)
	GetState() DecisionMakerState
}
