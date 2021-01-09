package broker

import (
	"fmt"
	"strings"

	krakenapi "github.com/beldur/kraken-go-api-client"
	"github.com/fabiodmferreira/crypto-trading/utils"
)

// KrakenBroker connects to kraken to sell or buy assets
type KrakenBroker struct {
	api    *krakenapi.KrakenAPI
	ticker string
}

// NewKrakenBroker returns an instance of a kraken broker
func NewKrakenBroker(api *krakenapi.KrakenAPI) *KrakenBroker {
	return &KrakenBroker{api, krakenapi.XXBTZEUR}
}

// SetTicker changes the ticker used to buy or sell assets
func (kb *KrakenBroker) SetTicker(ticker string) {
	switch strings.ToUpper(ticker) {
	case "BTC":
		kb.ticker = krakenapi.XXBTZEUR
	case "ETH":
		kb.ticker = krakenapi.XETHZEUR
	case "ADA":
		kb.ticker = krakenapi.ADAEUR
	case "DOT":
		kb.ticker = "DOTEUR"
	case "ATOM":
		kb.ticker = "ATOMEUR"
	default:
		panic(fmt.Sprintf("invalid ticker set in broker: %s", ticker))
	}
}

// AddBuyOrder request kraken to place a buy order with details passed by arguments
func (kb *KrakenBroker) AddBuyOrder(amount, price float32) error {
	amount = utils.RoundFloorTwoDecimals(amount)
	return kb.addOrder(amount, price, "buy")
}

// AddSellOrder request kraken to place a sell order with details passed by arguments
func (kb *KrakenBroker) AddSellOrder(amount, price float32) error {
	amount = utils.RoundFloorTwoDecimals(amount)
	return kb.addOrder(amount, price, "sell")
}

// addOrder is used by other methods to create orders in kraken
func (kb *KrakenBroker) addOrder(amount, price float32, orderType string) error {

	_, err := kb.api.AddOrder(kb.ticker, orderType, "limit", fmt.Sprintf("%f", amount), map[string]string{"price": fmt.Sprintf("%.1f", price)})
	return err
}

// BrokerMock is a broker stub to test it locally
type BrokerMock struct {
}

// NewBrokerMock returns an instance of a broker mock
func NewBrokerMock() *BrokerMock {
	return &BrokerMock{}
}

// SetTicker stub
func (bm *BrokerMock) SetTicker(ticker string) {
	fmt.Printf("Set ticker %s\n", ticker)
}

// AddBuyOrder stub
func (bm *BrokerMock) AddBuyOrder(amount, price float32) error {
	fmt.Printf("Add buy order (amount:%f,price:%f)\n", amount, price)
	return nil
}

// AddSellOrder stub
func (bm *BrokerMock) AddSellOrder(amount, price float32) error {
	fmt.Printf("Add sell order (amount:%f,price:%f)\n", amount, price)
	return nil
}
