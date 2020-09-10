package domain

import (
	"encoding/csv"
	"time"
)

// CollectorOptions are used to change collectors behaviour
type CollectorOptions struct {
	PriceVariationDetection float32 `bson:"priceVariationDetection,truncate" json:"priceVariationDetection"`
	DataSource              *csv.Reader
	NewPriceTimeRate        int `bson:"newPriceTimeRate,truncate" json:"newPriceTimeRate"`
}

// Collector notifies when price asset changes
type Collector interface {
	Start()
	Stop()
	Regist(observable OnNewAssetPrice)
	SetIndicators(indicators *[]Indicator)
}

// OHLC is a type with interval asset prices
type OHLC struct {
	Time    time.Time `json:"time"`
	EndTime time.Time `json:"etime"`
	Open    float32   `json:"open"`
	Close   float32   `json:"close"`
	High    float32   `json:"high"`
	Low     float32   `json:"low"`
	Volume  float32   `json:"volume"`
}

// OnNewAssetPrice
type OnNewAssetPrice = func(ohlc *OHLC)
