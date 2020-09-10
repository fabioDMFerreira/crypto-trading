package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	krakenapi "github.com/beldur/kraken-go-api-client"
	"github.com/fabiodmferreira/crypto-trading/accounts"
	"github.com/fabiodmferreira/crypto-trading/app"
	"github.com/fabiodmferreira/crypto-trading/appfactory"
	"github.com/fabiodmferreira/crypto-trading/collectors"
	"github.com/fabiodmferreira/crypto-trading/db"
	"github.com/fabiodmferreira/crypto-trading/domain"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type StartApplicationType func(metadata *domain.Application, appEnv string) (*app.App, error)

type mongoQuery struct {
	ID primitive.ObjectID `bson:"_id"  json:"_id"`
}

func main() {
	// load environment variables
	err := godotenv.Load()
	if err != nil {
		fmt.Println(".env file does not exist")
	}

	env := domain.Env{
		MongoURL:                    os.Getenv("MONGO_URL"),
		MongoDB:                     os.Getenv("MONGO_DB"),
		NotificationsReceiver:       os.Getenv("NOTIFICATIONS_RECEIVER"),
		NotificationsSender:         os.Getenv("NOTIFICATIONS_SENDER"),
		NotificationsSenderPassword: os.Getenv("NOTIFICATIONS_SENDER_PASSWORD"),
		AppEnv:                      os.Getenv("APP_ENV"),
		AppID:                       os.Getenv("APP_ID"),
	}

	// initialize third party instances
	krakenKey := os.Getenv("KRAKEN_API_KEY")
	krakenPrivateKey := os.Getenv("KRAKEN_PRIVATE_KEY")
	krakenAPI := krakenapi.New(krakenKey, krakenPrivateKey)

	dbClient, err := db.ConnectDB(env.MongoURL)

	if err != nil {
		log.Fatal("connecting db", err)
	}

	mongoDatabase := dbClient.Database(env.MongoDB)

	applicationsCollection := mongoDatabase.Collection(db.APPLICATIONS_COLLECTION)
	applicationsRepository := app.NewRepository(db.NewRepository(applicationsCollection))

	applications, err := applicationsRepository.FindAll()

	startApplication := startApplicationFactory(mongoDatabase, krakenAPI)

	appManager := make(map[string]*app.App)

	for _, metadata := range *applications {
		application, err := startApplication(&metadata, env.AppEnv)

		if err != nil {
			log.Fatal(err)
		}

		appManager[metadata.ID.Hex()] = application
	}

	if len(appManager) == 0 {
		fmt.Printf("Creating a default application")

		notificationOptions := domain.NotificationOptions{
			Receiver:       env.NotificationsReceiver,
			Sender:         env.NotificationsSender,
			SenderPassword: env.NotificationsSenderPassword,
		}

		accountsRepository := accounts.NewRepository(db.NewRepository(mongoDatabase.Collection(db.ACCOUNTS_COLLECTION)))

		metadata, err := appfactory.CreateDefaultAppMetadata(notificationOptions, applicationsRepository, accountsRepository)

		if err != nil {
			log.Fatal(err)
		}

		application, err := startApplication(metadata, env.AppEnv)

		if err != nil {
			log.Fatal(err)
		}

		appManager[metadata.ID.Hex()] = application
	}

	var waitGroup sync.WaitGroup
	waitGroup.Add(1)
	go watchApplicationsChanges(applicationsCollection, &waitGroup, func(change bson.M) {
		switch operationType := change["operationType"]; operationType {
		case "delete":
			metadata := parseBsonToAppMetadata(change["fullDocument"])
			if application, ok := appManager[metadata.ID.Hex()]; ok {
				application.Stop()
				appManager[metadata.ID.Hex()] = nil
			}
		case "update":
			var query mongoQuery
			bsonBytes, _ := bson.Marshal(change["documentKey"])
			bson.Unmarshal(bsonBytes, &query)
			if query.ID.Hex() != "000000000000000000000000" {
				metadata, err := applicationsRepository.FindByID(query.ID.Hex())
				if err != nil {
					fmt.Printf("Not able to start application with ID %v due to next error: %v", metadata.ID, err)
				}
				restartApp(metadata, appManager, env.AppEnv, startApplication)
			} else {
				fmt.Printf("Query should have an _id (%+v).", change["documentKey"])
			}
		case "replace", "insert":
			metadata := parseBsonToAppMetadata(change["fullDocument"])
			restartApp(metadata, appManager, env.AppEnv, startApplication)
		default:
			fmt.Printf("new event => %v", operationType)
		}
	})

	waitGroup.Wait()
}

func restartApp(metadata *domain.Application, appManager map[string]*app.App, appEnv string, startApplication StartApplicationType) {
	var err error
	if application, ok := appManager[metadata.ID.Hex()]; ok {
		application.Stop()
		newApp, err := startApplication(metadata, appEnv)

		if err != nil {
			appManager[metadata.ID.Hex()] = nil
		} else {
			appManager[metadata.ID.Hex()] = newApp
		}
	} else {
		newApp, err := startApplication(metadata, appEnv)

		if err != nil {
			appManager[metadata.ID.Hex()] = newApp
		}
	}

	if err != nil {
		fmt.Printf("Not able to start application with ID %v due to next error: %v", metadata.ID, err)
	}
}

func parseBsonToAppMetadata(data interface{}) *domain.Application {
	var metadata domain.Application
	bsonBytes, _ := bson.Marshal(data)
	bson.Unmarshal(bsonBytes, &metadata)
	return &metadata
}

func startApplicationFactory(mongoDatabase *mongo.Database, krakenAPI *krakenapi.KrakenAPI) StartApplicationType {
	return func(metadata *domain.Application, appEnv string) (*app.App, error) {
		brokerService := appfactory.GetBroker(appEnv, krakenAPI)

		collector := collectors.NewKrakenCollector(metadata.Asset, domain.CollectorOptions{NewPriceTimeRate: 1}, krakenAPI, &[]domain.Indicator{})

		application, err := appfactory.SetupApplication(metadata, mongoDatabase, brokerService, collector)

		if err != nil {
			return nil, err
		}

		go application.Start()

		return application, nil
	}
}

func watchApplicationsChanges(applicationsCollection *mongo.Collection, waitGroup *sync.WaitGroup, callback func(bson.M)) {
	defer waitGroup.Done()

	applicationsStream, err := applicationsCollection.Watch(context.TODO(), mongo.Pipeline{})
	if err != nil {
		panic(err)
	}

	defer applicationsStream.Close(context.TODO())

	for applicationsStream.Next(context.TODO()) {
		var data bson.M
		if err := applicationsStream.Decode(&data); err != nil {
			panic(err)
		}
		callback(data)
	}
}
