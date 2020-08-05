package app

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/fabiodmferreira/crypto-trading/domain"
)

// App holds instances of each application dependency and executes program
type App struct {
	decisionMaker       domain.DecisionMaker
	eventLogsRepository domain.EventsLog
	trader              domain.Trader
	accountService      domain.AccountService
	collector           domain.Collector
	Asset               string
}

// NewApp returns an instance of App
func NewApp(
	collector domain.Collector,
	decisionMaker domain.DecisionMaker,
	trader domain.Trader,
	accountService domain.AccountService,
) *App {
	app := &App{
		collector:      collector,
		decisionMaker:  decisionMaker,
		trader:         trader,
		accountService: accountService,
	}
	app.collector.Regist(app.OnTickerChange)
	return app
}

// SetEventsLog sets events logs repository
func (a *App) SetEventsLog(eventsLog domain.EventsLog) {
	a.eventLogsRepository = eventsLog
}

// log writes message to event log dependency
func (a *App) log(subject, message string) {
	if a.eventLogsRepository != nil {
		a.eventLogsRepository.Create(subject, message)
	} else {
		fmt.Printf("%v: %v", subject, message)
	}
}

// Start starts collecting data
func (a *App) Start() {
	a.collector.Start()
}

// RegistOnTickerChange executes function when the collector receives a change
func (a *App) RegistOnTickerChange(observable domain.OnTickerChange) {
	a.collector.Regist(observable)
}

// DecideToBuy do operations to check if an asset should be bought
func (a *App) DecideToBuy(price float32, currentTime time.Time) error {
	ok, err := a.decisionMaker.ShouldBuy(price, currentTime)
	if ok && err == nil {
		amount, err := a.decisionMaker.HowMuchAmountShouldBuy(price)

		if err != nil {
			return err
		}

		accountAmount, err := a.accountService.GetAmount()

		if err != nil {
			return err
		}

		if accountAmount > amount*price {
			err := a.trader.Buy(amount, price, currentTime)

			if err != nil {
				return err
			}

			message := fmt.Sprintf("Asset bought: {Price: %v Amount: %v Value: %v, Asset: %v}", price, amount, amount*price, a.Asset)
			a.log("buy", message)
		} else {
			a.log("Insuffucient Funds", fmt.Sprintf("want to spend %.4f%v*%.2f$=%v, have %.2f in account", amount, a.Asset, price, amount*price, accountAmount))
		}
	}

	return nil
}

// DecideToSell do operations to check if an asset should be sold
func (a *App) DecideToSell(price float32, currentTime time.Time) error {
	assets, err := a.accountService.FindPendingAssets()

	if err != nil {
		return err
	}

	for _, asset := range *assets {
		if ok, err := a.decisionMaker.ShouldSell(&asset, price, currentTime); ok && err == nil {

			if err := a.trader.Sell(&asset, price, currentTime); err != nil {
				return err
			}

			message := fmt.Sprintf("Asset sold: {Price: %v Amount: %v Value: %v, Asset: %v}", price, asset.Amount, price*asset.Amount, a.Asset)
			a.log("sell", message)
		}
	}

	return nil
}

// OnTickerChange do operations based on asset new price
func (a *App) OnTickerChange(ask, bid float32, currentTime time.Time) {

	a.decisionMaker.NewValue(ask, currentTime)
	a.log("Price change", fmt.Sprintf("%v PRICE: %v", a.Asset, ask))

	err := a.DecideToBuy(ask, currentTime)

	if err != nil {
		log.Fatal(err)
	}

	err = a.DecideToSell(ask, currentTime)

	if err != nil {
		log.Fatal(err)
	}
}

// FetchAssets returns all assets
func (a *App) FetchAssets() (*[]domain.Asset, error) {
	return a.accountService.FindAllAssets()
}

// GetAccountAmount returns the account service amount
func (a *App) GetAccountAmount() (float32, error) {
	return a.accountService.GetAmount()
}

// GetState returns application state
func (a *App) GetState() interface{} {
	accountAmount, _ := a.GetAccountAmount()

	decisionMakerState := a.decisionMaker.GetState()
	accountState := struct {
		AccountAmount float32 `json:"accountAmount"`
	}{AccountAmount: accountAmount}

	state, _ := merge(accountState, decisionMakerState)

	return state
}

// merge merges the two JSON-marshalable values x1 and x2,
// preferring x1 over x2 except where x1 and x2 are
// JSON objects, in which case the keys from both objects
// are included and their values merged recursively.
//
// It returns an error if x1 or x2 cannot be JSON-marshaled.
func merge(x1, x2 interface{}) (interface{}, error) {
	data1, err := json.Marshal(x1)
	if err != nil {
		return nil, err
	}
	data2, err := json.Marshal(x2)
	if err != nil {
		return nil, err
	}
	var j1 interface{}
	err = json.Unmarshal(data1, &j1)
	if err != nil {
		return nil, err
	}
	var j2 interface{}
	err = json.Unmarshal(data2, &j2)
	if err != nil {
		return nil, err
	}
	return merge1(j1, j2), nil
}

func merge1(x1, x2 interface{}) interface{} {
	switch x1 := x1.(type) {
	case map[string]interface{}:
		x2, ok := x2.(map[string]interface{})
		if !ok {
			return x1
		}
		for k, v2 := range x2 {
			if v1, ok := x1[k]; ok {
				x1[k] = merge1(v1, v2)
			} else {
				x1[k] = v2
			}
		}
	case nil:
		// merge(nil, map[string]interface{...}) -> map[string]interface{...}
		x2, ok := x2.(map[string]interface{})
		if ok {
			return x2
		}
	}
	return x1
}
