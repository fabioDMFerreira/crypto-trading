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
	V []interface{}
}

var Pairs = map[string]string{
	"BTC": "XBT/EUR",
	"ETH": "ETH/EUR",
	"ADA": "ADA/EUR",
}

// KrakenCollector collects data from kraken exchange
type KrakenCollector struct {
	options              domain.CollectorOptions
	krakenAPI            *krakenapi.KrakenAPI
	lastTickerPrice      float32
	observables          []domain.OnNewAssetPrice
	lastPricePublishDate time.Time
	pair                 string
	wscon                *websocket.Conn
}

// NewKrakenCollector returns an instance of KrakenCollector
func NewKrakenCollector(asset string, options domain.CollectorOptions, krakenAPI *krakenapi.KrakenAPI) *KrakenCollector {

	pair, ok := Pairs[asset]

	if !ok {
		log.Fatalf("%v does not have a valid kraken pair", asset)
	}

	return &KrakenCollector{options, krakenAPI, 0, []domain.OnNewAssetPrice{}, time.Time{}, pair, nil}
}

// Start connects to a kraken websocket that send prices variations
func (kc *KrakenCollector) Start() {
	u := url.URL{
		Scheme: "wss",
		Host:   "ws.kraken.com",
		Path:   "/",
	}

	con, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	kc.wscon = con

	if err != nil {
		log.Fatal(err)
	}

	defer kc.wscon.Close()

	subscribeEventMessage := fmt.Sprintf(`{
		"event": "subscribe",
		"pair": [
			"%v"
		],
		"subscription": {
			"name": "ohlc",
			"interval": 1
		}
	}`, kc.pair)

	err = kc.wscon.WriteMessage(
		websocket.TextMessage,
		[]byte(subscribeEventMessage),
	)
	if err != nil {
		log.Fatal(err)
	}

	var currentOHLC *domain.OHLC

	// receive message
	for {
		_, message, err := kc.wscon.ReadMessage()

		if err != nil {
			break
		}

		var e SocketEvent
		err = json.Unmarshal(message, &e)

		if e.Event != "heartbeat" {
			// https://eagain.net/articles/go-json-array-to-struct/
			var payload []interface{}
			msg := []interface{}{0, &payload, "", ""}
			err = json.Unmarshal(message, &msg)

			if err == nil {
				ohlc := getOHLCFromPayload(payload)

				if currentOHLC != nil && currentOHLC.Open != ohlc.Open {
					startDate, endDate := GetPreviousIntervalDates(time.Now())

					currentOHLC.Time = startDate
					currentOHLC.EndTime = endDate

					kc.PublishAssetPrice(currentOHLC)
				}

				currentOHLC = ohlc

				// askStr := msg[1].(*TickerMessage).A[0].(string)
				// ask, err := strconv.ParseFloat(askStr, 32)

				// askVolumeStr := msg[1].(*TickerMessage).A[2].(string)
				// askVolume, _ := strconv.ParseFloat(askVolumeStr, 32)

				// bidVolumeStr := msg[1].(*TickerMessage).A[2].(string)
				// bidVolume, _ := strconv.ParseFloat(bidVolumeStr, 32)
				// // bidStr := msg[1].(*TickerMessage).B[0].(string)
				// // bid, _ := strconv.ParseFloat(bidStr, 32)

				// _, minutes, seconds := time.Now().Clock()

				// if currentMinute == minutes {
				// 	currentAskVolume += askVolume
				// 	currentBidVolume += bidVolume
				// } else {
				// 	currentAskVolume = askVolume
				// 	currentBidVolume = bidVolume
				// 	currentMinute = minutes
				// }

				// fmt.Printf("%v:%v %v %v %v\n", currentMinute, seconds, currentAskVolume, currentBidVolume, currentAskVolume-currentBidVolume)

				// if err != nil {
				// 	fmt.Printf("error parsing price in message: %v", err)
				// 	return
				// }

				// price := float32(ask)

				// err = kc.HandlePriceChangeMessage(price, time.Now())

				// if err != nil {
				// 	fmt.Printf("error on handling price change message: %v", err)
				// }
			}

		}
	}
}

// Stop closes connection with kraken websocket
func (kc *KrakenCollector) Stop() {
	if kc.wscon != nil {
		kc.wscon.Close()
	}
}

func (kc *KrakenCollector) PublishAssetPrice(ohlc *domain.OHLC) error {
	for _, observable := range kc.observables {
		observable(ohlc)
	}

	return nil
}

// HandlePriceChangeMessage receives message, extracts parameters and call observable functions with the current asset price
// func (kc *KrakenCollector) HandlePriceChangeMessage(price float32, date time.Time) error {
// 	timeSinceLastPricePublished := date.Sub(kc.lastPricePublishDate).Minutes()

// 	changeVariance := kc.lastTickerPrice * kc.options.PriceVariationDetection

// 	if timeSinceLastPricePublished > float64(kc.options.NewPriceTimeRate) || (kc.lastTickerPrice == 0 ||
// 		price > kc.lastTickerPrice+changeVariance ||
// 		price < kc.lastTickerPrice-changeVariance) {
// 		kc.lastTickerPrice = price
// 		for _, observable := range kc.observables {
// 			observable(price, price, time.Now())

// 			kc.lastPricePublishDate = date
// 		}
// 	}

// 	return nil
// }

func GetPreviousIntervalDates(date time.Time) (time.Time, time.Time) {
	date = date.Add(time.Second * time.Duration(date.Second()) * -1)

	return date.Add(time.Minute * -1), date
}

// Regist add function to be executed when ticker price changes
func (kc *KrakenCollector) Regist(observable domain.OnNewAssetPrice) {
	kc.observables = append(kc.observables, observable)
}

func getOHLCFromPayload(msg []interface{}) *domain.OHLC {
	sTime, _ := strconv.ParseFloat(msg[0].(string), 32)
	startTime := time.Unix(int64(sTime), 0)

	etime, _ := strconv.ParseFloat(msg[1].(string), 32)
	endTime := time.Unix(int64(etime), 0)

	open, _ := strconv.ParseFloat(msg[2].(string), 32)
	high, _ := strconv.ParseFloat(msg[3].(string), 32)
	low, _ := strconv.ParseFloat(msg[4].(string), 32)
	close, _ := strconv.ParseFloat(msg[5].(string), 32)
	volume, _ := strconv.ParseFloat(msg[7].(string), 32)

	return &domain.OHLC{
		Time:    startTime,
		EndTime: endTime,
		Open:    float32(open),
		High:    float32(high),
		Low:     float32(low),
		Close:   float32(close),
		Volume:  float32(volume),
	}
}
