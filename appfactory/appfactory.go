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
	"github.com/fabiodmferreira/crypto-trading/indicators"
	"github.com/fabiodmferreira/crypto-trading/notifications"
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
	priceIndicator, volumeIndicator, err := setupIndicators(assetsPricesService, appMetaData.Asset, appMetaData.Options.StatisticsOptions)

	if err != nil {
		return nil, err
	}

	notificationsService := setupNotificationsService(mongoDatabase, appMetaData.Options.NotificationOptions, appMetaData.ID)
	decisionMaker := setupDecisionMaker(priceIndicator, volumeIndicator, accountService, appMetaData.Options.DecisionMakerOptions)
	dbTrader := trader.NewTrader(broker)

	// Create application
	application := app.NewApp(&[]domain.Collector{collector}, decisionMaker, dbTrader, accountService)
	application.Asset = appMetaData.Asset
	application.SetEventsLog(eventLogsRepository)

	// Regist events
	collector.Regist(NotificationJob(notificationsService, eventLogsRepository, accountService))
	collector.Regist(SaveAssetPrice(appMetaData.Asset, assetsPricesService))
	collector.Regist(SaveApplicationState(appMetaData.ID, application, applicationExecutionStateRepository))

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

		appMetaData, err = CreateDefaultAppMetadata(notificationOptions, applicationsRepository, accountsRepository)

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

func SaveAssetPrice(asset string, assetsPricesService domain.AssetsPricesService) domain.OnNewAssetPrice {
	return func(ohlc *domain.OHLC) {
		assetsPricesService.Create(ohlc, asset)
	}
}

func SaveApplicationState(ID primitive.ObjectID, application *app.App, applicationExecutionStateRepository domain.Repository) domain.OnNewAssetPrice {
	return func(ohlc *domain.OHLC) {
		state := domain.ApplicationExecutionState{
			ID:          primitive.NewObjectID(),
			ExecutionID: ID,
			Date:        ohlc.Time,
			State:       application.GetState(),
		}
		applicationExecutionStateRepository.InsertOne(state)
	}
}

func CreateDefaultAppMetadata(notificationOptions domain.NotificationOptions, repository domain.ApplicationRepository, accountsRepository domain.AccountsRepository) (*domain.Application, error) {
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
	priceIndicator *indicators.PriceIndicator,
	volumeIndicator *indicators.VolumeIndicator,
	accountService domain.AccountService,
	options domain.DecisionMakerOptions,
) domain.DecisionMaker {
	buyStrategy := decisionmaker.NewBuyStrategy(priceIndicator, volumeIndicator, accountService, options)
	sellStrategy := decisionmaker.NewSellStrategy(priceIndicator, volumeIndicator, accountService, options)

	return decisionmaker.NewDecisionMaker(buyStrategy, sellStrategy)
}

func NotificationJob(
	notificationsService domain.NotificationsService,
	eventLogsRepository domain.EventsLog,
	accountService domain.AccountService,
) func(ohlc *domain.OHLC) {
	return func(ohlc *domain.OHLC) {

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

func appendAssetsPricesToStatistics(priceIndicator *indicators.PriceIndicator, lastAssetsPrices *[]domain.AssetPrice) {
	points := []float64{}
	for _, assetPrice := range *lastAssetsPrices {
		points = append(points, float64(assetPrice.Close))
	}
	nPoints := len(points)
	for i := nPoints - 1; i >= 0; i-- {
		priceIndicator.AddValue(
			&domain.OHLC{
				Close:  float32(points[i]),
				Open:   float32(points[i]),
				High:   float32(points[i]),
				Low:    float32(points[i]),
				Volume: 0,
			},
		)
	}
}

func setupAssetsPricesService(mongoDatabase *mongo.Database) domain.AssetsPricesService {
	assetsPricesCollection := mongoDatabase.Collection(db.ASSETS_PRICES_COLLECTION)
	assetsPricesRepository := assetsprices.NewRepository(db.NewRepository(assetsPricesCollection))

	return assetsprices.NewService(assetsPricesRepository, assetsprices.NewCoindeskRemoteSource(http.Get).FetchRemoteAssetsPrices)
}

func setupIndicators(assetsPricesService domain.AssetsPricesService, asset string, statisticsOptions domain.StatisticsOptions) (*indicators.PriceIndicator, *indicators.VolumeIndicator, error) {
	lastAssetsPrices, err := getLastAssetsPrices(asset, statisticsOptions.NumberOfPointsHold, assetsPricesService)

	if err != nil {
		return nil, nil, fmt.Errorf("%v", err)
	}

	priceIndicator := indicators.NewPriceIndicator(indicators.NewMetricStatisticsIndicator(statisticsOptions))
	volumeIndicator := indicators.NewVolumeIndicator(indicators.NewMetricStatisticsIndicator(statisticsOptions))

	appendAssetsPricesToStatistics(priceIndicator, lastAssetsPrices)

	return priceIndicator, volumeIndicator, nil
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
