package app

import (
	"fmt"
	"log"
	"time"

	"github.com/fabiodmferreira/crypto-trading/domain"
	"github.com/fabiodmferreira/crypto-trading/notifications"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type App struct {
	notificationsService domain.NotificationsService
	decisionMaker        domain.DecisionMaker
	eventLogsRepository  domain.EventsLog
	assetsRepository     domain.AssetsRepositoryReader
	trader               domain.Trader
	accountService       domain.AccountService
}

func NewApp(
	notificationsService domain.NotificationsService,
	decisionMaker domain.DecisionMaker,
	log domain.EventsLog,
	assetsRepository domain.AssetsRepositoryReader,
	trader domain.Trader,
	accountService domain.AccountService,
) *App {
	return &App{notificationsService, decisionMaker, log, assetsRepository, trader, accountService}
}

func (a *App) OnTickerChange(ask, bid float32, buyTime time.Time) {

	a.decisionMaker.NewValue(ask)
	a.eventLogsRepository.Create("btc price change", fmt.Sprintf("BTC PRICE: %v", ask))

	ok, err := a.decisionMaker.ShouldBuy(ask, buyTime)
	if ok && err == nil {
		amount, err := a.decisionMaker.HowMuchAmountShouldBuy(ask)

		if err != nil {
			log.Fatal(err)
		}

		accountAmount, err := a.accountService.GetAmount()

		if err != nil {
			log.Fatal(err)
		}

		if accountAmount > amount*ask {
			err := a.trader.Buy(amount, ask, buyTime)

			if err != nil {
				log.Fatal(err)
			}

			message := fmt.Sprintf("Asset bought: {Price: %v Amount: %v Value: %v}", ask, amount, amount*ask)
			err = a.eventLogsRepository.Create("buy", message)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			a.eventLogsRepository.Create("Insuffucient Funds", fmt.Sprintf("want to spend %.4fBTC*%.2f$=%v, have %.2f in account", amount, ask, amount*ask, accountAmount))
		}
	}

	assets, err := a.assetsRepository.FindAll()

	if err != nil {
		log.Fatal(err)
	}

	for _, asset := range *assets {
		if ok, err := a.decisionMaker.ShouldSell(&asset, ask, buyTime); ok && err == nil {

			if err != nil {
				log.Fatal(err)
			}

			err := a.trader.Sell(&asset, ask)

			if err != nil {
				log.Fatal(err)
			}

			message := fmt.Sprintf("Asset sold: {Price: %v Amount: %v Value: %v}", ask, asset.Amount, ask*asset.Amount)
			err = a.eventLogsRepository.Create("sell", message)
			if err != nil {
				log.Fatal(err)
			}

		}
	}

	err = a.CheckEventLogs()

	if err != nil {
		log.Fatal(err)
	}
}

func (a *App) CheckEventLogs() error {
	lastNotificationTime, err := a.notificationsService.FindLastEventLogsNotificationDate()

	if err != nil || time.Now().Sub(lastNotificationTime).Hours() > 12 {
		eventLogs, err := a.eventLogsRepository.FindAllToNotify()

		if err != nil {
			return err
		}

		pendingAssets, err := a.assetsRepository.FindAll()
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
