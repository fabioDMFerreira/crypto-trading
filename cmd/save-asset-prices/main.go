package main

import (
	"encoding/csv"
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
	adadatahistory "github.com/fabiodmferreira/crypto-trading/data-history/ada"
	btcdatahistory "github.com/fabiodmferreira/crypto-trading/data-history/btc"
	ethdatahistory "github.com/fabiodmferreira/crypto-trading/data-history/eth"
	"github.com/fabiodmferreira/crypto-trading/db"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
)

// bulkUpdateElements is the number of elements that are used on bulk upsert
const bulkUpdateElements = 1000

// const csvFile = btcdatahistory.LastYearMinute
// const asset = "BTC"

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

	assetsPricesFiles := map[string]string{
		"ADA": adadatahistory.Twenty1920H1,
		"BTC": btcdatahistory.Twenty1920H1,
		// "BTC-CASH": btccashdatahistory.Twenty19,
		// "EOS":      eosdatahistory.Twenty19,
		// "ETC":      etcdatahistory.Twenty19,
		"ETH": ethdatahistory.Twenty1920H1,
		// "LTC":      ltcdatahistory.Twenty19,
		// "MONERO":   monerodatahistory.Twenty19,
		// "STELLAR":  stellardatahistory.Twenty19,
		// "XRP":      xrpdatahistory.Twenty19,
	}

	repo := db.NewRepository(collection)

	for asset, csvFile := range assetsPricesFiles {
		err := repo.BulkDelete(bson.M{"asset": asset})

		if err != nil {
			log.Fatalf("error on bulk deleting: %v", err)
		}

		historyFile, err := readHistoryFile(csvFile)

		if err != nil {
			log.Fatalf("error on getting cv: %v", err)
		}

		// read header
		_, err = historyFile.Read()
		if err == io.EOF {
			log.Fatalf("Error on reading header file: %v", err)
		}

		assetsPrices, err := getFileAssetsPrices(asset, historyFile)

		err = db.BatchBulkCreate(repo.BulkCreate, assetsPrices, bulkUpdateElements)

		if err != nil {
			log.Fatalf("Error on bulk creating assets prices: %v", err)
		}
	}
}

func getFileAssetsPrices(asset string, historyFile *csv.Reader) (*[]bson.M, error) {

	var documents []bson.M

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
				"c":     price,
				"date":  time.Unix(unixTime, 0),
			},
		)
	}

	return &documents, nil
}

func readHistoryFile(fileName string) (*csv.Reader, error) {
	_, currentFilePath, _, _ := runtime.Caller(0)
	currentDir := path.Dir(currentFilePath)
	return collectors.GetCsv(fmt.Sprintf("%v/../../data-history/%v", currentDir, fileName))
}
