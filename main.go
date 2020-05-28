package main

import (
	"fmt"
	"log"
	"os"

	krakenapi "github.com/beldur/kraken-go-api-client"
	"github.com/fabiodmferreira/crypto-trading/accounts"
	"github.com/fabiodmferreira/crypto-trading/assets"
	"github.com/fabiodmferreira/crypto-trading/collectors"
	"github.com/fabiodmferreira/crypto-trading/db"
	"github.com/fabiodmferreira/crypto-trading/decisionmaker"
	"github.com/fabiodmferreira/crypto-trading/eventlogs"
	"github.com/fabiodmferreira/crypto-trading/notifications"
	"github.com/fabiodmferreira/crypto-trading/trader"
	"github.com/joho/godotenv"
)

const (
	MaximumBuyAmount        = 0.01
	PretendedProfitPerSold  = 0.1
	PriceDropToBuy          = 0.1
	PriceVariationDetection = 0.01
)

func main() {
	// load environment variables
	err := godotenv.Load()
	if err != nil {
		fmt.Println(".env file does not exist")
	}

	krakenKey := os.Getenv("KRAKEN_API_KEY")
	krakenPrivateKey := os.Getenv("KRAKEN_PRIVATE_KEY")
	mongoURL := os.Getenv("MONGO_URL")
	mongoDB := os.Getenv("MONGO_DB")
	notificationsReceiver := os.Getenv("NOTIFICATIONS_RECEIVER")
	notificationsSender := os.Getenv("NOTIFICATIONS_SENDER")
	notificationsSenderPassword := os.Getenv("NOTIFICATIONS_SENDER_PASSWORD")

	// initialize third party instances
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
	dbTrader := trader.NewDBTrader(assetsRepository, eventLogsRepository)
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

	account = accounts.NewAccountService(accountDocument.ID, accountsRepository)
	decisionMaker := decisionmaker.NewDecisionMaker(dbTrader, account)

	var lastPrice float32

	onTickerChange := func(ask, bid float32) {
		if lastPrice == 0 ||
			ask > lastPrice+(lastPrice*PriceVariationDetection) ||
			ask < lastPrice-(lastPrice*PriceVariationDetection) {
			lastPrice = ask
			eventLogsRepository.Create("btc price change", fmt.Sprintf("BTC PRICE: %v", lastPrice))

			assets, err := assetsRepository.FindAll()
			if err == nil {
				decisionMaker.DecideToSell(ask, assets, PretendedProfitPerSold)
			}

			cheaperAssetPrice, err := assetsRepository.FindCheaperAssetPrice()
			if err == nil {
				decisionMaker.DecideToBuy(ask, cheaperAssetPrice, PriceDropToBuy, MaximumBuyAmount)
			} else {
				log.Fatal(err)
			}
		}

		notificationsService.CheckEventLogs()
	}

	collectors.KrakenTickerCollector(krakenAPI, onTickerChange)
}
