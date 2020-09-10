package app

import (
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
	collectors          *[]domain.Collector
	Asset               string
}

// NewApp returns an instance of App
func NewApp(
	collectors *[]domain.Collector,
	decisionMaker domain.DecisionMaker,
	trader domain.Trader,
	accountService domain.AccountService,
) *App {
	app := &App{
		collectors:     collectors,
		decisionMaker:  decisionMaker,
		trader:         trader,
		accountService: accountService,
	}

	app.RegistOnNewAssetPrice(app.OnNewAssetPrice)

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
	}
}

// Start starts collecting data
func (a *App) Start() {
	for _, collector := range *a.collectors {
		collector.Start()
	}
}

// Stop stops collecting data
func (a *App) Stop() {
	for _, collector := range *a.collectors {
		collector.Stop()
	}
}

// RegistOnNewAssetPrice executes function when the collector receives a change
func (a *App) RegistOnNewAssetPrice(observable domain.OnNewAssetPrice) {
	for _, collector := range *a.collectors {
		collector.Regist(observable)
	}
}

// DecideToBuy do operations to check if an asset should be bought
func (a *App) DecideToBuy(price float32, currentTime time.Time) error {
	ok, amount, err := a.decisionMaker.ShouldBuy()
	if ok && err == nil {
		accountAmount, err := a.accountService.GetAmount()

		if err != nil {
			return err
		}

		if accountAmount > amount*price {
			err := a.trader.Buy(amount, price, currentTime)

			if err != nil {
				return err
			}

			err = a.accountService.Withdraw(amount * price)

			if err != nil {
				return err
			}

			_, err = a.accountService.CreateAsset(amount, price, currentTime)

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

	ok, _, err := a.decisionMaker.ShouldSell()

	if ok {
		for _, asset := range *assets {
			if asset.BuyPrice+(asset.BuyPrice*0.01) < price {
				if err := a.trader.Sell(&asset, price, currentTime); err != nil {
					return err
				}

				err := a.accountService.SellAsset(asset.ID.Hex(), price, currentTime)

				if err != nil {
					return err
				}

				err = a.accountService.Deposit(asset.Amount * price)

				if err != nil {
					return err
				}

				message := fmt.Sprintf("Asset sold: {Price: %v Amount: %v Value: %v, Asset: %v}", price, asset.Amount, price*asset.Amount, a.Asset)
				a.log("sell", message)
			}
		}
	}

	return nil
}

// OnNewAssetPrice do operations based on asset new price
func (a *App) OnNewAssetPrice(ohlc *domain.OHLC) {
	a.log("Price change", fmt.Sprintf("%v PRICE: %v", a.Asset, ohlc.Close))

	err := a.DecideToBuy(ohlc.Close, ohlc.Time)

	if err != nil {
		log.Fatal(err)
	}

	err = a.DecideToSell(ohlc.Close, ohlc.Time)

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

	accountState := struct {
		AccountAmount float32 `json:"accountAmount"`
	}{AccountAmount: accountAmount}

	return accountState
}
