package domain

import "time"

// Trader buys and sells assets
type Trader interface {
	Buy(amount, price float32, buyTime time.Time) error
	Sell(asset *Asset, price float32, sellTime time.Time) error
}
