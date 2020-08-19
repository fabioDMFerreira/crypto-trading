package main

import (
	"fmt"
	"log"
	"os"

	krakenapi "github.com/beldur/kraken-go-api-client"
	"github.com/fabiodmferreira/crypto-trading/appfactory"
	"github.com/fabiodmferreira/crypto-trading/db"
	"github.com/fabiodmferreira/crypto-trading/domain"
	"github.com/joho/godotenv"
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

	// initialize third party instances
	krakenKey := os.Getenv("KRAKEN_API_KEY")
	krakenPrivateKey := os.Getenv("KRAKEN_PRIVATE_KEY")
	krakenAPI := krakenapi.New(krakenKey, krakenPrivateKey)

	dbClient, err := db.ConnectDB(env.MongoURL)

	if err != nil {
		log.Fatal("connecting db", err)
	}

	mongoDatabase := dbClient.Database(env.MongoDB)

	application, err := appfactory.SetupApplication(env, mongoDatabase, krakenAPI)

	if err != nil {
		log.Fatal(err)
	}

	application.Start()
}
