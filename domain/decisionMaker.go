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
	Price                  float32 `bson:"price,truncate" json:"price"`
	PriceAverage           float64 `bson:"priceAverage,truncate" json:"priceAverage"`
	PriceStandardDeviation float64 `bson:"priceStandardDeviation,truncate" json:"priceStandardDeviation"`
	PriceUpperLimit        float64 `bson:"priceUpperLimit,truncate" json:"priceUpperLimit"`
	PriceLowerLimit        float64 `bson:"priceLowerLimit,truncate" json:"priceLowerLimit"`

	Change                  float32 `bson:"change,truncate" json:"change"`
	ChangeAverage           float64 `bson:"changeAverage,truncate" json:"changeAverage"`
	ChangeStandardDeviation float64 `bson:"changeStandardDeviation,truncate" json:"changeStandardDeviation"`
	ChangeUpperLimit        float64 `bson:"changeUpperLimit,truncate" json:"changeUpperLimit"`
	ChangeLowerLimit        float64 `bson:"changeLowerLimit,truncate" json:"changeLowerLimit"`

	Acceleration                  float32 `bson:"acceleration,truncate" json:"acceleration"`
	AccelerationAverage           float64 `bson:"accelerationAverage,truncate" json:"accelerationAverage"`
	AccelerationStandardDeviation float64 `bson:"accelerationStandardDeviation,truncate" json:"accelerationStandardDeviation"`
	AccelerationUpperLimit        float64 `bson:"accelerationUpperLimit,truncate" json:"accelerationUpperLimit"`
	AccelerationLowerLimit        float64 `bson:"accelerationLowerLimit,truncate" json:"accelerationLowerLimit"`

	AccountAmount float64 `bson:"accountAmount,truncate" json:"accountAmount"`
}

// DecisionMaker makes decisions to buy or sell assets
type DecisionMaker interface {
	NewValue(price float32, date time.Time)
	ShouldBuy(price float32, buyTime time.Time) (bool, error)
	ShouldSell(asset *Asset, price float32, buyTime time.Time) (bool, error)
	HowMuchAmountShouldBuy(price float32) (float32, error)
	GetState() DecisionMakerState
}
