package main

import (
	"fmt"
	"log"
	"os"

	krakenapi "github.com/beldur/kraken-go-api-client"
	"github.com/fabiodmferreira/crypto-trading/accounts"
	"github.com/fabiodmferreira/crypto-trading/app"
	"github.com/fabiodmferreira/crypto-trading/assets"
	"github.com/fabiodmferreira/crypto-trading/broker"
	"github.com/fabiodmferreira/crypto-trading/collectors"
	"github.com/fabiodmferreira/crypto-trading/db"
	"github.com/fabiodmferreira/crypto-trading/decisionmaker"
	"github.com/fabiodmferreira/crypto-trading/domain"
	"github.com/fabiodmferreira/crypto-trading/eventlogs"
	"github.com/fabiodmferreira/crypto-trading/notifications"
	"github.com/fabiodmferreira/crypto-trading/statistics"
	"github.com/fabiodmferreira/crypto-trading/trader"
	"github.com/joho/godotenv"
)

func main() {
	// load environment variables
	err := godotenv.Load()
	if err != nil {
		fmt.Println(".env file does not exist")
	}

	mongoURL := os.Getenv("MONGO_URL")
	mongoDB := os.Getenv("MONGO_DB")
	notificationsReceiver := os.Getenv("NOTIFICATIONS_RECEIVER")
	notificationsSender := os.Getenv("NOTIFICATIONS_SENDER")
	notificationsSenderPassword := os.Getenv("NOTIFICATIONS_SENDER_PASSWORD")
	appEnv := os.Getenv("APP_ENV")

	// initialize third party instances
	krakenKey := os.Getenv("KRAKEN_API_KEY")
	krakenPrivateKey := os.Getenv("KRAKEN_PRIVATE_KEY")
	krakenAPI := krakenapi.New(krakenKey, krakenPrivateKey)

	dbClient, err := db.ConnectDB(mongoURL)

	if err != nil {
		log.Fatal("connecting db", err)
	}

	mongoDatabase := dbClient.Database(mongoDB)
	assetsCollection := mongoDatabase.Collection(db.ASSETS_COLLECTION)
	eventLogsCollection := mongoDatabase.Collection(db.EVENT_LOGS_COLLECTION)
	notificationsCollection := mongoDatabase.Collection(db.NOTIFICATIONS_COLLECTION)
	accountsCollection := mongoDatabase.Collection(db.ACCOUNTS_COLLECTION)

	// instantiate repositories
	assetsRepository := assets.NewAssetsRepository(assetsCollection)
	eventLogsRepository := eventlogs.NewEventLogsRepository(eventLogsCollection)
	notificationsRepository := notifications.NewNotificationsRepository(notificationsCollection)
	accountsRepository := accounts.NewAccountsRepository(accountsCollection)

	// instantiate services
	var brokerService domain.Broker
	if appEnv == "production" {
		brokerService = broker.NewKrakenBroker(krakenAPI)
	} else {
		fmt.Println("Broker mocked!")
		brokerService = broker.NewBrokerMock()
	}

	var accountService *accounts.AccountService
	accountDocument, err := accountsRepository.FindByBroker("kraken")
	if err != nil {
		accountDocument, err = accountsRepository.Create("kraken", 5000)

		if err != nil {
			log.Fatal("creating account", err)
		}
	}
	accountService = accounts.NewAccountService(accountDocument.ID, accountsRepository, assetsRepository)

	dbTrader := trader.NewTrader(assetsRepository, accountService, brokerService)
	notificationsService := notifications.NewNotificationsService(
		notificationsRepository,
		notificationsReceiver,
		notificationsSender,
		notificationsSenderPassword,
	)

	decisionmakerOptions :=
		decisionmaker.Options{
			MaximumBuyAmount:      0.01,
			MinimumProfitPerSold:  0.01,
			MinimumPriceDropToBuy: 0.01,
		}

	statisticsOptions := statistics.Options{NumberOfPointsHold: 38000}
	macdParams := statistics.MACDParams{Fast: 24, Slow: 12, Lag: 9}
	macd := statistics.NewMACDContainer(macdParams)
	statistics := statistics.NewStatistics(statisticsOptions, macd)

	decisionMaker := decisionmaker.NewDecisionMaker(assetsRepository, decisionmakerOptions, statistics)

	krakenCollector := collectors.NewKrakenCollector(krakenAPI, 0.01)
	application := app.NewApp(notificationsService, decisionMaker, eventLogsRepository, assetsRepository, dbTrader, accountService)

	krakenCollector.Start(application.OnTickerChange)
}
