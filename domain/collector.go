package domain

import (
	"encoding/csv"
	"time"
)

// CollectorOptions are used to change collectors behaviour
type CollectorOptions struct {
	PriceVariationDetection float32 `json:"priceVariationDetection"`
	DataSource              *csv.Reader
}

// Collector notifies when price asset changes
type Collector interface {
	Start()
	Regist(observable OnTickerChange)
}

// OnTickerChange is a function that receives 2 floats (ask and bid) and the timestamp of a price change
type OnTickerChange = func(float32, float32, time.Time)
