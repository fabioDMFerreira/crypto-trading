package app

import (
	"fmt"
	"log"
	"time"

	"github.com/fabiodmferreira/crypto-trading/domain"
)

const ()

type NotificationsService interface {
	CheckEventLogs() error
}

type DecisionMaker interface {
	MakeDecisions(price float32, buyTime time.Time)
}

type App struct {
	lastTickerPrice         float32
	notificationsService    NotificationsService
	decisionMaker           DecisionMaker
	eventLogsRepository     domain.EventsLog
	PriceVariationDetection float32
}

func NewApp(notificationsService NotificationsService, decisionMaker DecisionMaker, log domain.EventsLog, priceVariationDetection float32) *App {
	return &App{0, notificationsService, decisionMaker, log, priceVariationDetection}
}

func (a *App) OnTickerChange(ask, bid float32, buyTime time.Time) {
	if a.lastTickerPrice == 0 ||
		ask > a.lastTickerPrice+(a.lastTickerPrice*a.PriceVariationDetection) ||
		ask < a.lastTickerPrice-(a.lastTickerPrice*a.PriceVariationDetection) {
		a.lastTickerPrice = ask
		a.eventLogsRepository.Create("btc price change", fmt.Sprintf("BTC PRICE: %v", a.lastTickerPrice))

		a.decisionMaker.MakeDecisions(ask, buyTime)
	}

	err := a.notificationsService.CheckEventLogs()

	if err != nil {
		log.Fatal(err)
	}
}
