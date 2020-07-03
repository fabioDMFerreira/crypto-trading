package domain

// Broker add order to buy and sell assets in real brokers
type Broker interface {
	AddBuyOrder(amount, price float32) error
	AddSellOrder(amount, price float32) error
}
