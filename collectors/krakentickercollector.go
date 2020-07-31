package collectors

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"

	krakenapi "github.com/beldur/kraken-go-api-client"
	"github.com/fabiodmferreira/crypto-trading/domain"
	"github.com/gorilla/websocket"
)

// SocketEvent is a type used to decode kraken websocket messages
type SocketEvent struct {
	Event string
}

// TickerMessage is a type used to decode messages of ticker price change events
type TickerMessage struct {
	A []interface{}
	B []interface{}
}

// KrakenCollector collects data from kraken exchange
type KrakenCollector struct {
	options              domain.CollectorOptions
	krakenAPI            *krakenapi.KrakenAPI
	lastTickerPrice      float32
	observables          []domain.OnTickerChange
	lastPricePublishDate time.Time
}

// NewKrakenCollector returns an instance of KrakenCollector
func NewKrakenCollector(options domain.CollectorOptions, krakenAPI *krakenapi.KrakenAPI) *KrakenCollector {
	return &KrakenCollector{options, krakenAPI, 0, []domain.OnTickerChange{}, time.Time{}}
}

// Start connects to a kraken websocket that send prices variations
func (kc *KrakenCollector) Start() {
	u := url.URL{
		Scheme: "wss",
		Host:   "ws.kraken.com",
		Path:   "/",
	}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal(err)
	}

	defer c.Close()

	err = c.WriteMessage(
		websocket.TextMessage,
		[]byte(`{
			"event": "subscribe",
			"pair": [
				"XBT/EUR"
			],
			"subscription": {
				"name": "ticker"
			}
		}`),
	)
	if err != nil {
		log.Fatal(err)
	}

	// receive message
	for {
		_, message, err := c.ReadMessage()

		if err != nil {
			log.Fatal(err)
		}

		var e SocketEvent
		err = json.Unmarshal(message, &e)

		if e.Event != "heartbeat" {
			// https://eagain.net/articles/go-json-array-to-struct/
			msg := []interface{}{0, &TickerMessage{}, "", ""}
			err = json.Unmarshal(message, &msg)

			if err == nil {
				askStr := msg[1].(*TickerMessage).A[0].(string)
				// bidStr := msg[1].(*TickerMessage).B[0].(string)
				ask, err := strconv.ParseFloat(askStr, 32)
				// bid, err2 := strconv.ParseFloat(bidStr, 32)

				if err != nil {
					fmt.Printf("error parsing price in message: %v", err)
					return
				}

				price := float32(ask)

				err = kc.HandlePriceChangeMessage(price, time.Now())

				if err != nil {
					fmt.Printf("error on handling price change message: %v", err)
				}
			}

		}
	}
}

// HandlePriceChangeMessage receives message, extracts parameters and call observable functions with the current asset price
func (kc *KrakenCollector) HandlePriceChangeMessage(price float32, date time.Time) error {
	timeSinceLastPricePublished := date.Sub(kc.lastPricePublishDate).Minutes()

	if timeSinceLastPricePublished < float64(kc.options.NewPriceTimeRate) {
		return nil
	}

	changeVariance := kc.lastTickerPrice * kc.options.PriceVariationDetection

	if kc.lastTickerPrice == 0 ||
		price > kc.lastTickerPrice+changeVariance ||
		price < kc.lastTickerPrice-changeVariance {
		kc.lastTickerPrice = price
		for _, observable := range kc.observables {
			observable(price, price, time.Now())

			kc.lastPricePublishDate = date
		}
	}

	return nil
}

// Regist add function to be executed when ticker price changes
func (kc *KrakenCollector) Regist(observable domain.OnTickerChange) {
	kc.observables = append(kc.observables, observable)
}
