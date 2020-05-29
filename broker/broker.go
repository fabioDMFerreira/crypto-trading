package broker

import (
	"fmt"

	krakenapi "github.com/beldur/kraken-go-api-client"
)

type KrakenBroker struct {
	api    *krakenapi.KrakenAPI
	ticker string
}

func NewKrakenBroker(api *krakenapi.KrakenAPI) *KrakenBroker {
	return &KrakenBroker{api, krakenapi.XXBTZEUR}
}

func (kb *KrakenBroker) addOrder(amount, price float32, orderType string) error {

	_, err := kb.api.AddOrder(kb.ticker, orderType, "limit", fmt.Sprintf("%f", amount), map[string]string{"price": fmt.Sprintf("%.1f", price)})
	return err
}

func (kb *KrakenBroker) AddBuyOrder(amount, price float32) error {
	return kb.addOrder(amount, price, "buy")
}

func (kb *KrakenBroker) AddSellOrder(amount, price float32) error {
	return kb.addOrder(amount, price, "sell")
}

type BrokerMock struct {
}

func NewBrokerMock() *BrokerMock {
	return &BrokerMock{}
}

func (bm *BrokerMock) AddBuyOrder(amount, price float32) error {
	return nil
}

func (bm *BrokerMock) AddSellOrder(amount, price float32) error {
	return nil
}
