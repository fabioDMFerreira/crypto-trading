package domain

import "time"

// DecisionMakerOptions are used to change DecisionMaker behaviour
type DecisionMakerOptions struct {
	MaximumBuyAmount      float32 `json:"maximumBuyAmount"`
	MinimumProfitPerSold  float32 `json:"minimumProfitPerSold"`
	MinimumPriceDropToBuy float32 `json:"minimumPriceDropToBuy"`
}

// DecisionMaker makes decisions to buy or sell assets
type DecisionMaker interface {
	NewValue(price float32)
	ShouldBuy(price float32, buyTime time.Time) (bool, error)
	ShouldSell(asset *Asset, price float32, buyTime time.Time) (bool, error)
	HowMuchAmountShouldBuy(price float32) (float32, error)
}
