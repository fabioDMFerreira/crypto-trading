package app

import (
	"fmt"
	"log"
	"time"

	"github.com/fabiodmferreira/crypto-trading/domain"
	"github.com/fabiodmferreira/crypto-trading/notifications"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// App holds instances of each application dependency and executes program
type App struct {
	notificationsService domain.NotificationsService
	decisionMaker        domain.DecisionMaker
	eventLogsRepository  domain.EventsLog
	assetsRepository     domain.AssetsRepositoryReader
	trader               domain.Trader
	accountService       domain.AccountService
	collector            domain.Collector
}

// NewApp returns an instance of App
func NewApp(
	notificationsService domain.NotificationsService,
	decisionMaker domain.DecisionMaker,
	log domain.EventsLog,
	assetsRepository domain.AssetsRepositoryReader,
	trader domain.Trader,
	accountService domain.AccountService,
	collector domain.Collector,
) *App {
	app := &App{notificationsService, decisionMaker, log, assetsRepository, trader, accountService, collector}
	app.collector.Regist(app.OnTickerChange)
	return app
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

			message := fmt.Sprintf("Asset bought: {Price: %v Amount: %v Value: %v}", price, amount, amount*price)
			err = a.eventLogsRepository.Create("buy", message)
			if err != nil {
				return err
			}
		} else {
			a.eventLogsRepository.Create("Insuffucient Funds", fmt.Sprintf("want to spend %.4fBTC*%.2f$=%v, have %.2f in account", amount, price, amount*price, accountAmount))
		}
	}

	return nil
}

// DecideToSell do operations to check if an asset should be sold
func (a *App) DecideToSell(price float32, currentTime time.Time) error {
	assets, err := a.assetsRepository.FindPendingAssets()

	if err != nil {
		return err
	}

	for _, asset := range *assets {
		if ok, err := a.decisionMaker.ShouldSell(&asset, price, currentTime); ok && err == nil {

			if err != nil {
				return err
			}

			err := a.trader.Sell(&asset, price, currentTime)

			if err != nil {
				return err
			}

			message := fmt.Sprintf("Asset sold: {Price: %v Amount: %v Value: %v}", price, asset.Amount, price*asset.Amount)
			err = a.eventLogsRepository.Create("sell", message)
			if err != nil {
				return err
			}

		}
	}

	return nil
}

// OnTickerChange do operations based on asset new price
func (a *App) OnTickerChange(ask, bid float32, currentTime time.Time) {

	a.decisionMaker.NewValue(ask, currentTime)
	a.eventLogsRepository.Create("btc price change", fmt.Sprintf("BTC PRICE: %v", ask))

	err := a.DecideToBuy(ask, currentTime)

	if err != nil {
		log.Fatal(err)
	}

	err = a.DecideToSell(ask, currentTime)

	if err != nil {
		log.Fatal(err)
	}

	err = a.CheckEventLogs()

	if err != nil {
		log.Fatal(err)
	}
}

// CheckEventLogs verifies wheter there are log events to notify the user
func (a *App) CheckEventLogs() error {
	lastNotificationTime, err := a.notificationsService.FindLastEventLogsNotificationDate()

	if err != nil || time.Now().Sub(lastNotificationTime).Hours() > 12 {
		eventLogs, err := a.eventLogsRepository.FindAllToNotify()

		if err != nil {
			return err
		}

		pendingAssets, err := a.assetsRepository.FindPendingAssets()
		if err != nil {
			return err
		}

		accountAmount, err := a.accountService.GetAmount()

		if err != nil {
			return err
		}

		startDate, endDate := lastNotificationTime, time.Now()
		balance, err := a.assetsRepository.GetBalance(startDate, endDate)

		if err != nil {
			return err
		}

		subject := "Crypto-Trading: Report"
		var eventLogsIds []primitive.ObjectID

		for _, event := range *eventLogs {
			eventLogsIds = append(eventLogsIds, event.ID)
		}

		message, err := notifications.GenerateEventlogReportEmail(accountAmount, len(*pendingAssets), balance, startDate, endDate, eventLogs, pendingAssets)

		if err != nil {
			return err
		}

		err = a.notificationsService.CreateEmailNotification(subject, message.String(), "eventlogs")

		if err != nil {
			log.Fatal(err)
		}

		err = a.eventLogsRepository.MarkNotified(eventLogsIds)
		if err != nil {
			return err
		}
	}

	return nil
}

// FetchAssets returns all assets
func (a *App) FetchAssets() (*[]domain.Asset, error) {
	return a.assetsRepository.FindAll()
}

// GetAccountAmount returns the account service amount
func (a *App) GetAccountAmount() (float32, error) {
	return a.accountService.GetAmount()
}

// GetState returns application state
func (a *App) GetState() interface{} {
	return a.decisionMaker.GetState()
}
