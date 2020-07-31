package main

import (
	"bytes"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"time"

	krakenapi "github.com/beldur/kraken-go-api-client"
	"github.com/fabiodmferreira/crypto-trading/accounts"
	"github.com/fabiodmferreira/crypto-trading/app"
	"github.com/fabiodmferreira/crypto-trading/assets"
	"github.com/fabiodmferreira/crypto-trading/assetsprices"
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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

	var brokerService domain.Broker
	if appEnv == "production" {
		brokerService = broker.NewKrakenBroker(krakenAPI)
	} else {
		fmt.Println("Broker mocked!")
		brokerService = broker.NewBrokerMock()
	}

	notificationOptions := domain.NotificationOptions{
		Receiver:       notificationsReceiver,
		Sender:         notificationsSender,
		SenderPassword: notificationsSenderPassword,
	}

	assetsCollection := mongoDatabase.Collection(db.ASSETS_COLLECTION)
	assetsRepository := assets.NewRepository(db.NewRepository(assetsCollection))

	eventLogsCollection := mongoDatabase.Collection(db.EVENT_LOGS_COLLECTION)

	accountService := setupAccountService(mongoDatabase, assetsRepository)

	notificationsService := setupNotificationsService(mongoDatabase, notificationOptions)
	decisionMaker := setupDecisionMaker(assetsRepository, mongoDatabase)
	eventLogsRepository := eventlogs.NewEventLogsRepository(eventLogsCollection)
	dbTrader := trader.NewTrader(assetsRepository, accountService, brokerService)
	krakenCollector := collectors.NewKrakenCollector(domain.CollectorOptions{PriceVariationDetection: 0.01}, krakenAPI)

	application := app.NewApp(
		krakenCollector,
		decisionMaker,
		dbTrader,
		accountService,
	)

	application.SetEventsLog(eventLogsRepository)

	application.RegistOnTickerChange(NotificationJob(notificationsService, eventLogsRepository, accountService, assetsRepository))

	application.Start()
}

func setupNotificationsService(mongoDatabase *mongo.Database, notificationOptions domain.NotificationOptions) domain.NotificationsService {
	notificationsCollection := mongoDatabase.Collection(db.NOTIFICATIONS_COLLECTION)

	notificationsRepository := notifications.NewRepository(db.NewRepository(notificationsCollection))

	return notifications.NewService(
		notificationsRepository,
		notificationOptions,
		smtp.SendMail,
	)
}

func setupDecisionMaker(assetsRepository domain.AssetsRepository, mongoDatabase *mongo.Database) domain.DecisionMaker {
	assetsPricesCollection := mongoDatabase.Collection(db.ASSETS_PRICES_COLLECTION)
	assetsPricesRepository := assetsprices.NewRepository(db.NewRepository(assetsPricesCollection))
	assetsPricesService := assetsprices.NewService(assetsPricesRepository, assetsprices.FetchCoindeskRemotePrices)

	decisionmakerOptions :=
		domain.DecisionMakerOptions{
			MinimumProfitPerSold:     0.01,
			MinimumPriceDropToBuy:    0.01,
			MaximumFIATBuyAmount:     500,
			GrowthDecreaseLimit:      -100,
			GrowthIncreaseLimit:      100,
			MinutesToCollectNewPoint: 15,
		}
	numberOfPointsHold := 5000

	statisticsOptions := domain.StatisticsOptions{NumberOfPointsHold: numberOfPointsHold}
	macdParams := statistics.MACDParams{Fast: 24, Slow: 12, Lag: 9}
	macd := statistics.NewMACDContainer(macdParams)
	pricesStatistics := statistics.NewStatistics(statisticsOptions, macd)
	growthMacd := statistics.NewMACDContainer(macdParams)
	growthStatistics := statistics.NewStatistics(statisticsOptions, growthMacd)

	fmt.Println("Fetching assets prices from coindesk...")
	assetsPricesService.FetchAndStoreAssetPrices("BTC", time.Now())
	fmt.Println("Completed")

	lastAssetsPrices, err := assetsPricesService.GetLastAssetsPrices("BTC", numberOfPointsHold)

	if err != nil {
		log.Panicf("%v", err)
	}

	fmt.Println("Adding statistics points...")
	points := []float64{}
	for _, assetPrice := range *lastAssetsPrices {
		points = append(points, float64(assetPrice.Value))
	}
	nPoints := len(points)
	for i := nPoints - 1; i >= 0; i-- {
		pricesStatistics.AddPoint(points[i])
	}
	fmt.Println("Completed")

	return decisionmaker.NewDecisionMaker(assetsRepository, decisionmakerOptions, pricesStatistics, growthStatistics, assetsPricesService)
}

func setupAccountService(mongoDatabase *mongo.Database, assetsRepository domain.AssetsRepository) domain.AccountService {
	accountsCollection := mongoDatabase.Collection(db.ACCOUNTS_COLLECTION)
	accountsRepository := accounts.NewRepository(accountsCollection)

	accountDocument, err := accountsRepository.FindByBroker("kraken")
	if err != nil {
		accountDocument, err = accountsRepository.Create("kraken", 5000)

		if err != nil {
			log.Fatal("creating account", err)
		}
	}

	return accounts.NewAccountService(accountDocument.ID, accountsRepository, assetsRepository)
}

func NotificationJob(
	notificationsService domain.NotificationsService,
	eventLogsRepository domain.EventsLog,
	accountService domain.AccountService,
	assetsRepository domain.AssetsRepository) func(ask, bid float32, date time.Time) {
	return func(ask, bid float32, date time.Time) {

		shouldSendNotification := notificationsService.ShouldSendNotification()

		if !shouldSendNotification {
			return
		}

		eventLogs, err := eventLogsRepository.FindAllToNotify()

		if err != nil {
			fmt.Println(err)
			return
		}

		pendingAssets, err := accountService.FindPendingAssets()
		if err != nil {
			fmt.Println(err)
			return
		}

		accountAmount, err := accountService.GetAmount()

		if err != nil {
			fmt.Println(err)
			return
		}

		lastNotificationTime, err := notificationsService.FindLastEventLogsNotificationDate()

		if err != nil {
			fmt.Println(err)
			lastNotificationTime = time.Now()
		}

		startDate, endDate := lastNotificationTime, time.Now()

		balance, err := assetsRepository.GetBalance(startDate, endDate)

		if err != nil {
			fmt.Println(err)
			return
		}

		message, err := notifications.GenerateEventlogReportEmail(accountAmount, balance, startDate, endDate, eventLogs, pendingAssets)

		if err != nil {
			fmt.Println(err)
			return
		}

		err = sendReport(notificationsService, message)

		if err != nil {
			markLogsEventsAsNotified(eventLogsRepository, eventLogs)
		}
	}
}

func sendReport(notificationsService domain.NotificationsService, message *bytes.Buffer) error {
	subject := "Crypto-Trading: Report"

	err := notificationsService.CreateEmailNotification(subject, message.String(), "eventlogs")

	if err != nil {
		fmt.Println(err)
	}

	return err
}

func markLogsEventsAsNotified(eventLogsRepository domain.EventsLog, eventLogs *[]domain.EventLog) {
	var eventLogsIds []primitive.ObjectID

	for _, event := range *eventLogs {
		eventLogsIds = append(eventLogsIds, event.ID)
	}

	err := eventLogsRepository.MarkNotified(eventLogsIds)

	if err != nil {
		fmt.Println(err)
	}
}
