package appfactory

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"time"

	krakenapi "github.com/beldur/kraken-go-api-client"
	"github.com/fabiodmferreira/crypto-trading/accounts"
	"github.com/fabiodmferreira/crypto-trading/app"
	"github.com/fabiodmferreira/crypto-trading/assets"
	"github.com/fabiodmferreira/crypto-trading/assetsprices"
	"github.com/fabiodmferreira/crypto-trading/broker"
	"github.com/fabiodmferreira/crypto-trading/db"
	"github.com/fabiodmferreira/crypto-trading/decisionmaker"
	"github.com/fabiodmferreira/crypto-trading/domain"
	"github.com/fabiodmferreira/crypto-trading/eventlogs"
	"github.com/fabiodmferreira/crypto-trading/notifications"
	"github.com/fabiodmferreira/crypto-trading/statistics"
	"github.com/fabiodmferreira/crypto-trading/trader"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupApplication(appMetaData *domain.Application, mongoDatabase *mongo.Database, broker domain.Broker, collector domain.Collector) (*app.App, error) {
	// Setup repositories
	assetsCollection := mongoDatabase.Collection(db.ASSETS_COLLECTION)
	assetsRepository := assets.NewRepository(db.NewRepository(assetsCollection))

	accountsCollection := mongoDatabase.Collection(db.ACCOUNTS_COLLECTION)
	accountsRepository := accounts.NewRepository(db.NewRepository(accountsCollection))

	applicationExecutionStateCollection := mongoDatabase.Collection(db.APPLICATION_EXECUTION_STATES_COLLECTION)
	applicationExecutionStateRepository := db.NewRepository(applicationExecutionStateCollection)

	// Setup services
	accountService, err := accounts.NewAccountService(appMetaData.AccountID.Hex(), accountsRepository, assetsRepository)

	if err != nil {
		return nil, err
	}

	eventLogsCollection := mongoDatabase.Collection(db.EVENT_LOGS_COLLECTION)
	eventLogsRepository := eventlogs.NewEventLogsRepository(db.NewRepository(eventLogsCollection), appMetaData.ID)

	assetsPricesService := setupAssetsPricesService(mongoDatabase)
	pricesStatistics, growthStatistics, err := setupPricesStatistics(assetsPricesService, appMetaData.Asset, appMetaData.Options.StatisticsOptions)

	if err != nil {
		return nil, err
	}

	notificationsService := setupNotificationsService(mongoDatabase, appMetaData.Options.NotificationOptions, appMetaData.ID)
	decisionMaker := setupDecisionMaker(appMetaData.Options.DecisionMakerOptions, pricesStatistics, growthStatistics, accountService)
	dbTrader := trader.NewTrader(accountService, broker)

	// Create application
	application := app.NewApp(collector, decisionMaker, dbTrader, accountService)
	application.Asset = appMetaData.Asset
	application.SetEventsLog(eventLogsRepository)

	// Regist events
	application.RegistOnTickerChange(NotificationJob(notificationsService, eventLogsRepository, accountService))
	application.RegistOnTickerChange(SaveAssetPrice(appMetaData.Asset, assetsPricesService))
	application.RegistOnTickerChange(SaveApplicationState(appMetaData.ID, application, applicationExecutionStateRepository))

	return application, nil
}

func FindOrCreateAppMetaData(env domain.Env, applicationsRepository domain.ApplicationRepository, accountsRepository domain.AccountsRepository) (*domain.Application, error) {
	var appMetaData *domain.Application
	var err error

	if env.AppID == "" {
		notificationOptions := domain.NotificationOptions{
			Receiver:       env.NotificationsReceiver,
			Sender:         env.NotificationsSender,
			SenderPassword: env.NotificationsSenderPassword,
		}

		appMetaData, err = createDefaultAppMetadata(notificationOptions, applicationsRepository, accountsRepository)

		if err != nil {
			log.Fatalf("Not able to create a new application due to %v", err)
		}

		fmt.Printf("Application created with id %v\n", appMetaData.ID)

	} else {
		appMetaData, err = applicationsRepository.FindByID(env.AppID)

		if err != nil {
			return nil, fmt.Errorf("Not able to get application with id %v due to %v", env.AppID, err)
		}

	}

	return appMetaData, nil
}

func SaveAssetPrice(asset string, assetsPricesService domain.AssetsPricesService) domain.OnTickerChange {
	return func(ask, bid float32, date time.Time) {
		assetsPricesService.Create(date, ask, asset)
	}
}

func SaveApplicationState(ID primitive.ObjectID, application *app.App, applicationExecutionStateRepository domain.Repository) domain.OnTickerChange {
	return func(ask, bid float32, date time.Time) {
		state := domain.ApplicationExecutionState{
			ID:          primitive.NewObjectID(),
			ExecutionID: ID,
			Date:        date,
			State:       application.GetState(),
		}
		applicationExecutionStateRepository.InsertOne(state)
	}
}

func createDefaultAppMetadata(notificationOptions domain.NotificationOptions, repository domain.ApplicationRepository, accountsRepository domain.AccountsRepository) (*domain.Application, error) {
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
			NewPriceTimeRate:        1,
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

func getLastAssetsPrices(asset string, numberOfPoints int, assetsPricesService domain.AssetsPricesService) (*[]domain.AssetPrice, error) {
	fmt.Println("Fetching assets prices from coindesk...")
	assetsPricesService.FetchAndStoreAssetPrices(asset, time.Now())
	fmt.Println("Completed")

	return assetsPricesService.GetLastAssetsPrices(asset, numberOfPoints)
}

func appendAssetsPricesToStatistics(statistics domain.Statistics, lastAssetsPrices *[]domain.AssetPrice) {
	points := []float64{}
	for _, assetPrice := range *lastAssetsPrices {
		points = append(points, float64(assetPrice.Value))
	}
	nPoints := len(points)
	for i := nPoints - 1; i >= 0; i-- {
		statistics.AddPoint(points[i])
	}
}

func setupAssetsPricesService(mongoDatabase *mongo.Database) domain.AssetsPricesService {
	assetsPricesCollection := mongoDatabase.Collection(db.ASSETS_PRICES_COLLECTION)
	assetsPricesRepository := assetsprices.NewRepository(db.NewRepository(assetsPricesCollection))

	return assetsprices.NewService(assetsPricesRepository, assetsprices.NewCoindeskRemoteSource(http.Get).FetchRemoteAssetsPrices)
}

func setupPricesStatistics(assetsPricesService domain.AssetsPricesService, asset string, statisticsOptions domain.StatisticsOptions) (domain.Statistics, domain.Statistics, error) {
	lastAssetsPrices, err := getLastAssetsPrices(asset, statisticsOptions.NumberOfPointsHold, assetsPricesService)

	if err != nil {
		return nil, nil, fmt.Errorf("%v", err)
	}

	pricesStatistics := setupStatistics(statisticsOptions)
	growthStatistics := setupStatistics(statisticsOptions)

	appendAssetsPricesToStatistics(pricesStatistics, lastAssetsPrices)

	return pricesStatistics, growthStatistics, nil
}

func setupStatistics(statisticsOptions domain.StatisticsOptions) domain.Statistics {
	macdParams := statistics.MACDParams{Fast: 24, Slow: 12, Lag: 9}
	macd := statistics.NewMACDContainer(macdParams)
	return statistics.NewStatistics(statisticsOptions, macd)
}

func GetBroker(appEnv string, krakenAPI *krakenapi.KrakenAPI) domain.Broker {
	var brokerService domain.Broker
	if appEnv == "production" {
		brokerService = broker.NewKrakenBroker(krakenAPI)
	} else {
		fmt.Println("Broker mocked!")
		brokerService = broker.NewBrokerMock()
	}
	return brokerService
}
