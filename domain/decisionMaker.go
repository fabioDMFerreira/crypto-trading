package domain

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

	Volume           float32 `bson:"volume" json:"volume"`
	VolumeAverage    float64 `bson:"volumeAverage" json:"volumeAverage"`
	VolumeUpperLimit float64 `bson:"volumeUpperLimit" json:"volumeUpperLimit"`
	VolumeLowerLimit float64 `bson:"volumeLowerLimit" json:"volumeLowerLimit"`

	Open  float32 `bson:"open" json:"open"`
	Close float32 `bson:"close" json:"close"`
	High  float32 `bson:"high" json:"high"`
	Low   float32 `bson:"low" json:"low"`

	AccountAmount float64 `bson:"accountAmount,truncate" json:"accountAmount"`
}

// DecisionMaker makes decisions to buy or sell assets
type DecisionMaker interface {
	ShouldBuy() (bool, float32, error)
	ShouldSell() (bool, float32, error)
}

// Strategy uses data to calculate if it is a good time to do an action (buy/sell) and the sureness of doing it
type Strategy interface {
	Execute() (bool, float32, error)
}
