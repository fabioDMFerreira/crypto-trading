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
	"github.com/fabiodmferreira/crypto-trading/eventlogs"
	"github.com/fabiodmferreira/crypto-trading/notifications"
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
	krakenBroker := broker.NewKrakenBroker(krakenAPI)
	dbTrader := trader.NewTrader(assetsRepository, eventLogsRepository, krakenBroker)
	notificationsService := notifications.NewNotificationsService(
		notificationsRepository,
		eventLogsRepository,
		notificationsReceiver,
		notificationsSender,
		notificationsSenderPassword,
	)

	var account *accounts.AccountService
	accountDocument, err := accountsRepository.FindByBroker("kraken")
	if err != nil {
		accountDocument, err = accountsRepository.Create("kraken", 5000)

		if err != nil {
			log.Fatal("creating account", err)
		}
	}

	decisionmakerOptions := decisionmaker.DecisionMakerOptions{0.01, 0.01, 0.1}

	account = accounts.NewAccountService(accountDocument.ID, accountsRepository)
	decisionMaker := decisionmaker.NewDecisionMaker(dbTrader, account, assetsRepository, decisionmakerOptions)

	krakenCollector := collectors.NewKrakenCollector(krakenAPI)
	application := app.NewApp(notificationsService, decisionMaker, eventLogsRepository, 0.01)

	krakenCollector.Start(application.OnTickerChange)
}
