package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/fabiodmferreira/crypto-trading/collectors"
	btcdatahistory "github.com/fabiodmferreira/crypto-trading/data-history/btc"
	"github.com/fabiodmferreira/crypto-trading/db"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
)

// bulkUpdateElements is the number of elements that are used on bulk upsert
const bulkUpdateElements = 1000
const csvFile = btcdatahistory.LastYearMinute
const asset = "BTC"

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
		log.Fatal("connecting db: ", err)
	}

	mongoDatabase := dbClient.Database(mongoDB)

	collection := mongoDatabase.Collection(db.ASSETS_PRICES_COLLECTION)

	repo := db.NewRepository(collection)

	err = repo.BulkDelete(bson.M{"asset": asset})

	if err != nil {
		log.Fatal("error on bulk deleting: ", err)
	}

	_, currentFilePath, _, _ := runtime.Caller(0)
	currentDir := path.Dir(currentFilePath)
	historyFile, err := collectors.GetCsv(fmt.Sprintf("%v/../../data-history/%v", currentDir, csvFile))

	if err != nil {
		log.Fatal("error on getting cv: ", err)
	}

	// read header
	_, err = historyFile.Read()
	if err == io.EOF {
		log.Fatalf("Error on reading header file: %v", err)
	}

	var documents []bson.M
	counter := 0

	for {
		// Read each record from csv
		record, err := historyFile.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		priceStr := strings.ReplaceAll(record[1], ",", "")
		price, err := strconv.ParseFloat(priceStr, 32)

		if err != nil {
			log.Fatalf("Error on converting price from file:\n%v", err)
		}

		unixTime, err := strconv.ParseInt(record[0], 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		documents = append(documents,
			bson.M{
				"asset": asset,
				"value": price,
				"date":  time.Unix(unixTime/1000, 0),
			},
		)

		if len(documents) == bulkUpdateElements {
			err := repo.BulkCreate(&documents)

			if err != nil {
				log.Fatal("error on bulking create: ", err)
			}

			counter += bulkUpdateElements
			fmt.Printf("\rCreated: %d", counter)
			documents = []bson.M{}
		}
	}
}
