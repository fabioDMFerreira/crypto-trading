package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/fabiodmferreira/crypto-trading/assetsprices"
	"github.com/fabiodmferreira/crypto-trading/benchmark"
	"github.com/fabiodmferreira/crypto-trading/db"
	"github.com/fabiodmferreira/crypto-trading/webserver"
	"github.com/gorilla/handlers"
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

	dbClient, err := db.ConnectDB(mongoURL)

	if err != nil {
		log.Fatal("connecting db", err)
	}

	mongoDatabase := dbClient.Database(mongoDB)
	benchmarksCollection := mongoDatabase.Collection(db.BENCHMARKS_COLLECTION)
	assetspricesCollection := mongoDatabase.Collection(db.ASSETS_PRICES_COLLECTION)

	benchmarkRepository := benchmark.NewRepository(db.NewRepository(benchmarksCollection))
	assetspricesRepository := assetsprices.NewRepository(db.NewRepository(assetspricesCollection))
	benchmarkService := benchmark.NewService(benchmarkRepository, assetspricesRepository)

	server, err := webserver.NewCryptoTradingServer(benchmarkService, assetspricesRepository)

	if err != nil {
		log.Fatalf("problem creating server, %v ", err)
	}

	if err := http.ListenAndServe(":5000", handlers.LoggingHandler(os.Stdout, server)); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
