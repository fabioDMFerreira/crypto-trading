package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/fabiodmferreira/crypto-trading/accounts"
	"github.com/fabiodmferreira/crypto-trading/app"
	"github.com/fabiodmferreira/crypto-trading/applicationExecutionStates"
	"github.com/fabiodmferreira/crypto-trading/assets"
	"github.com/fabiodmferreira/crypto-trading/assetsprices"
	"github.com/fabiodmferreira/crypto-trading/benchmark"
	"github.com/fabiodmferreira/crypto-trading/db"
	"github.com/fabiodmferreira/crypto-trading/eventlogs"
	"github.com/fabiodmferreira/crypto-trading/notifications"
	"github.com/fabiodmferreira/crypto-trading/webserver"
	"github.com/gorilla/handlers"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	// load environment variables
	err := godotenv.Load()
	if err != nil {
		fmt.Println(".env file does not exist")
	}

	mongoURL := os.Getenv("MONGO_URL")
	mongoDB := os.Getenv("MONGO_DB")

	dbClient, err := db.ConnectDB(mongoURL)

	if err != nil {
		log.Fatal("connecting db", err)
	}

	mongoDatabase := dbClient.Database(mongoDB)
	benchmarksCollection := mongoDatabase.Collection(db.BENCHMARKS_COLLECTION)
	assetspricesCollection := mongoDatabase.Collection(db.ASSETS_PRICES_COLLECTION)
	applicationExecutionStatesCollection := mongoDatabase.Collection(db.APPLICATION_EXECUTION_STATES_COLLECTION)

	benchmarkRepository := benchmark.NewRepository(db.NewRepository(benchmarksCollection))
	assetspricesRepository := assetsprices.NewRepository(db.NewRepository(assetspricesCollection))
	applicationExecutionStatesRepository := applicationExecutionStates.NewRepository(db.NewRepository(applicationExecutionStatesCollection))
	benchmarkService := benchmark.NewService(benchmarkRepository, assetspricesRepository, applicationExecutionStatesRepository)

	accountsCollection := mongoDatabase.Collection(db.ACCOUNTS_COLLECTION)
	accountsRepository := accounts.NewRepository(db.NewRepository(accountsCollection))

	assetsCollection := mongoDatabase.Collection(db.ASSETS_COLLECTION)
	assetsRepository := assets.NewRepository(db.NewRepository(assetsCollection))

	logEventsCollection := mongoDatabase.Collection(db.EVENT_LOGS_COLLECTION)
	logEventsRepository := eventlogs.NewEventLogsRepository(db.NewRepository(logEventsCollection), primitive.NewObjectID())

	notificationsCollection := mongoDatabase.Collection(db.NOTIFICATIONS_COLLECTION)
	notificationsRepository := notifications.NewRepository(db.NewRepository(notificationsCollection))

	applicationsCollection := mongoDatabase.Collection(db.APPLICATIONS_COLLECTION)
	applicationsRepository := app.NewRepository(db.NewRepository(applicationsCollection))
	applicationsService := app.NewService(applicationsRepository, applicationExecutionStatesRepository, logEventsRepository, notificationsRepository)

	server, err := webserver.NewCryptoTradingServer(benchmarkService, assetspricesRepository, accountsRepository, assetsRepository, applicationsService)

	if err != nil {
		log.Fatalf("problem creating server, %v ", err)
	}

	if err := http.ListenAndServe(":4000", handlers.LoggingHandler(os.Stdout, server)); err != nil {
		log.Fatalf("could not listen on port 4000 %v", err)
	}
}
