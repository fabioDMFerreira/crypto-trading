package app

import (
	"fmt"
	"log"
	"time"

	"github.com/fabiodmferreira/crypto-trading/assets"
	"github.com/fabiodmferreira/crypto-trading/domain"
)

const ()

type NotificationsService interface {
	CheckEventLogs() error
}

type DecisionMaker interface {
	ShouldBuy(price float32, buyTime time.Time) (bool, error)
	ShouldSell(asset *assets.Asset, price float32, buyTime time.Time) (bool, error)
	HowMuchAmountShouldBuy(price float32) (float32, error)
}

type App struct {
	lastTickerPrice         float32
	notificationsService    NotificationsService
	decisionMaker           DecisionMaker
	eventLogsRepository     domain.EventsLog
	PriceVariationDetection float32
	assetsRepository        domain.AssetsRepositoryReader
	trader                  domain.Trader
	accountService          domain.AccountService
}

func NewApp(
	notificationsService NotificationsService,
	decisionMaker DecisionMaker,
	log domain.EventsLog,
	priceVariationDetection float32,
	assetsRepository domain.AssetsRepositoryReader,
	trader domain.Trader,
	accountService domain.AccountService,
) *App {
	return &App{0, notificationsService, decisionMaker, log, priceVariationDetection, assetsRepository, trader, accountService}
}

func (a *App) OnTickerChange(ask, bid float32, buyTime time.Time) {
	if a.lastTickerPrice == 0 ||
		ask > a.lastTickerPrice+(a.lastTickerPrice*a.PriceVariationDetection) ||
		ask < a.lastTickerPrice-(a.lastTickerPrice*a.PriceVariationDetection) {
		a.lastTickerPrice = ask
		a.eventLogsRepository.Create("btc price change", fmt.Sprintf("BTC PRICE: %v", a.lastTickerPrice))

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
	}

	err := a.notificationsService.CheckEventLogs()

	if err != nil {
		log.Fatal(err)
	}
}
