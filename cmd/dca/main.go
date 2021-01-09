package main

import (
	"fmt"
	"log"
	"os"
	"time"

	krakenapi "github.com/beldur/kraken-go-api-client"
	"github.com/fabiodmferreira/crypto-trading/broker"
	"github.com/fabiodmferreira/crypto-trading/collectors"
	"github.com/fabiodmferreira/crypto-trading/db"
	"github.com/fabiodmferreira/crypto-trading/dca"
	"github.com/fabiodmferreira/crypto-trading/domain"
	"github.com/joho/godotenv"
)

var (
	// WeekSeconds is the number of seconds in a week
	WeekSeconds = 604800
)

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

	dbClient, err := db.ConnectDB(env.MongoURL)

	if err != nil {
		log.Fatal("connecting db", err)
	}

	mongoDatabase := dbClient.Database(env.MongoDB)
	dcaJobsCollection := mongoDatabase.Collection(db.DCA_JOBS_COLLECTION)
	dcaJobsRepo := dca.NewJobsRepository(db.NewRepository(dcaJobsCollection))

	dcaAssetsCollection := mongoDatabase.Collection(db.DCA_ASSETS_COLLECTION)
	dcaAssetsRepo := dca.NewAssetsRepository(db.NewRepository(dcaAssetsCollection))

	// initialize third party instances
	krakenKey := os.Getenv("KRAKEN_API_KEY")
	krakenPrivateKey := os.Getenv("KRAKEN_PRIVATE_KEY")
	krakenAPI := krakenapi.New(krakenKey, krakenPrivateKey)

	collector := collectors.NewKrakenCollector("BTC", domain.CollectorOptions{}, krakenAPI, &[]domain.Indicator{})
	trader := broker.NewKrakenBroker(krakenAPI)

	service := dca.NewService(trader, collector, dcaJobsRepo, dcaAssetsRepo)

	if len(os.Args) > 1 && os.Args[1] == "create" {
		dcaJob := &domain.DCAJob{
			NextExecution: time.Now().Unix(),
			Period:        int64(WeekSeconds),
			Options: domain.DCAJobOptions{
				TotalFIATAmount: 200,
				CoinsProportion: map[string]float32{
					"BTC":  60,
					"ETH":  15,
					"DOT":  10,
					"ADA":  10,
					"ATOM": 5,
				},
			},
		}

		err = service.CreateDCA(dcaJob)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err = service.DrainDCA()
		if err != nil {
			log.Fatal(err)
		}
	}

}
