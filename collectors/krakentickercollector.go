package collectors

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"

	krakenapi "github.com/beldur/kraken-go-api-client"
	"github.com/gorilla/websocket"
)

type OnTickerChange = func(float32, float32, time.Time)

type SocketEvent struct {
	Event string
}

type TickerMessage struct {
	A []interface{}
	B []interface{}
}

type KrakenCollector struct {
	krakenAPI               *krakenapi.KrakenAPI
	priceVariationDetection float32
	lastTickerPrice         float32
}

func NewKrakenCollector(krakenAPI *krakenapi.KrakenAPI, priceVariationDetection float32) *KrakenCollector {
	return &KrakenCollector{krakenAPI, priceVariationDetection, 0}
}

func (kc *KrakenCollector) Start(onChange OnTickerChange) {
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
	// send message
	// event := Event{event: "heartbeat"}
	// message, err := json.Marshal(event)
	// if message != []byte(`{"event":"heartbeat"}`) {
	// 	log.Fatal()
	// }
	// fmt.Printf("%v\n", string(message))
	// if err != nil {
	// 	log.Fatal(err)
	// }

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
				ask, err1 := strconv.ParseFloat(askStr, 32)
				// bid, err2 := strconv.ParseFloat(bidStr, 32)

				if err1 == nil {
					price := float32(ask)

					changeVariance := kc.lastTickerPrice * kc.priceVariationDetection

					if kc.lastTickerPrice == 0 ||
						price > kc.lastTickerPrice+changeVariance ||
						price < kc.lastTickerPrice-changeVariance {
						kc.lastTickerPrice = price
						onChange(price, price, time.Now())
					}
				} else {
					fmt.Printf("%v", err1)
				}

			}
			// else {
			// 	fmt.Println("%v", string(message))
			// }
			// msgTicker := msg[1]
			// fmt.Printf("%v", msgTicker)
			// fmt.Println(string(message))
		}
	}
}
