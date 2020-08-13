package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
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
	appID := os.Getenv("APP_ID")

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

	accountsCollection := mongoDatabase.Collection(db.ACCOUNTS_COLLECTION)
	accountsRepository := accounts.NewRepository(db.NewRepository(accountsCollection))

	applicationsCollection := mongoDatabase.Collection(db.APPLICATIONS_COLLECTION)
	applicationsRepository := app.NewRepository(db.NewRepository(applicationsCollection))

	var appMetaData *domain.Application
	var accountService domain.AccountService
	if appID == "" {
		appMetaData, err = createDefaultApplicationMetaData(notificationOptions, applicationsRepository, accountsRepository)

		if err != nil {
			log.Fatalf("Not able to create a new application due to %v", err)
		}

		accountService = accounts.NewAccountService(appMetaData.AccountID, accountsRepository, assetsRepository)

		fmt.Printf("Application created with id %v\n", appMetaData.ID)
	} else {
		appMetaData, err = applicationsRepository.FindByID(appID)

		if err != nil {
			log.Fatalf("Not able to get application with id %v due to %v", appID, err)
		}

		account, err := accountsRepository.FindById(appMetaData.AccountID)

		if err != nil {
			log.Fatalf("Not able to get account with id %v due to %v", appMetaData.AccountID, err)
		}

		accountService = accounts.NewAccountService(account.ID, accountsRepository, assetsRepository)
	}

	eventLogsCollection := mongoDatabase.Collection(db.EVENT_LOGS_COLLECTION)

	assetsPricesCollection := mongoDatabase.Collection(db.ASSETS_PRICES_COLLECTION)
	assetsPricesRepository := assetsprices.NewRepository(db.NewRepository(assetsPricesCollection))
	assetsPricesService := assetsprices.NewService(assetsPricesRepository, assetsprices.NewCoindeskRemoteSource(http.Get).FetchRemoteAssetsPrices)

	fmt.Println("Fetching assets prices from coindesk...")
	assetsPricesService.FetchAndStoreAssetPrices(appMetaData.Asset, time.Now())
	fmt.Println("Completed")

	statisticsOptions := appMetaData.Options.StatisticsOptions
	lastAssetsPrices, err := assetsPricesService.GetLastAssetsPrices(appMetaData.Asset, statisticsOptions.NumberOfPointsHold)

	if err != nil {
		log.Panicf("%v", err)
	}

	pricesStatistics := setupStatistics(statisticsOptions)
	growthStatistics := setupStatistics(statisticsOptions)
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

	notificationsService := setupNotificationsService(mongoDatabase, notificationOptions, appMetaData.ID)
	decisionmakerOptions := appMetaData.Options.DecisionMakerOptions
	decisionMaker := setupDecisionMaker(decisionmakerOptions, pricesStatistics, growthStatistics, accountService)
	eventLogsRepository := eventlogs.NewEventLogsRepository(db.NewRepository(eventLogsCollection), appMetaData.ID)
	dbTrader := trader.NewTrader(accountService, brokerService)
	krakenCollector := collectors.NewKrakenCollector(appMetaData.Asset, appMetaData.Options.CollectorOptions, krakenAPI)

	application := app.NewApp(
		krakenCollector,
		decisionMaker,
		dbTrader,
		accountService,
	)
	application.Asset = appMetaData.Asset

	application.SetEventsLog(eventLogsRepository)

	application.RegistOnTickerChange(NotificationJob(notificationsService, eventLogsRepository, accountService))

	application.RegistOnTickerChange(func(ask, bid float32, date time.Time) {
		assetsPricesService.Create(date, ask, appMetaData.Asset)
	})

	applicationExecutionStateCollection := mongoDatabase.Collection(db.APPLICATION_EXECUTION_STATES_COLLECTION)
	applicationExecutionStateRepository := db.NewRepository(applicationExecutionStateCollection)

	application.RegistOnTickerChange(func(ask, bid float32, date time.Time) {
		state := domain.ApplicationExecutionState{
			ID:          primitive.NewObjectID(),
			ExecutionID: appMetaData.ID,
			Date:        date,
			State:       application.GetState(),
		}
		applicationExecutionStateRepository.InsertOne(state)
	})

	application.Start()
}

func createDefaultApplicationMetaData(notificationOptions domain.NotificationOptions, repository domain.ApplicationRepository, accountsRepository *accounts.Repository) (*domain.Application, error) {
	options := domain.ApplicationOptions{
		NotificationOptions: notificationOptions,
		StatisticsOptions:   domain.StatisticsOptions{NumberOfPointsHold: 5000},
		DecisionMakerOptions: domain.DecisionMakerOptions{
			MinimumProfitPerSold:  0.01,
			MinimumPriceDropToBuy: 0.01,
			MaximumFIATBuyAmount:  500,
			GrowthDecreaseLimit:   -100,
			GrowthIncreaseLimit:   100,
		},
		CollectorOptions: domain.CollectorOptions{
			PriceVariationDetection: 0.01,
			NewPriceTimeRate:        15,
		},
	}

	account, err := accountsRepository.Create("kraken", 5000)

	if err != nil {
		return nil, err
	}

	return repository.Create("BTC", options, account.ID)
}

func setupNotificationsService(mongoDatabase *mongo.Database, notificationOptions domain.NotificationOptions, appID primitive.ObjectID) domain.NotificationsService {
	notificationsCollection := mongoDatabase.Collection(db.NOTIFICATIONS_COLLECTION)

	notificationsRepository := notifications.NewRepository(db.NewRepository(notificationsCollection))

	return notifications.NewService(
		notificationsRepository,
		notificationOptions,
		smtp.SendMail,
		appID,
	)
}

func setupDecisionMaker(
	decisionmakerOptions domain.DecisionMakerOptions,
	pricesStatistics domain.Statistics,
	growthStatistics domain.Statistics,
	accountService domain.AccountService,
) domain.DecisionMaker {
	return decisionmaker.NewDecisionMaker(accountService, decisionmakerOptions, pricesStatistics, growthStatistics)
}

func setupStatistics(statisticsOptions domain.StatisticsOptions) domain.Statistics {
	macdParams := statistics.MACDParams{Fast: 24, Slow: 12, Lag: 9}
	macd := statistics.NewMACDContainer(macdParams)
	return statistics.NewStatistics(statisticsOptions, macd)
}

func NotificationJob(
	notificationsService domain.NotificationsService,
	eventLogsRepository domain.EventsLog,
	accountService domain.AccountService,
) func(ask, bid float32, date time.Time) {
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

		balance, err := accountService.GetBalance(startDate, endDate)

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
			fmt.Println(err)
			return
		}

		markLogsEventsAsNotified(eventLogsRepository, eventLogs)
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
