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
	btccashdatahistory "github.com/fabiodmferreira/crypto-trading/data-history/btc-cash"
	btcsvdatahistory "github.com/fabiodmferreira/crypto-trading/data-history/btc-sv"
	eosdatahistory "github.com/fabiodmferreira/crypto-trading/data-history/eos"
	etcdatahistory "github.com/fabiodmferreira/crypto-trading/data-history/etc"
	ethdatahistory "github.com/fabiodmferreira/crypto-trading/data-history/eth"
	ltcdatahistory "github.com/fabiodmferreira/crypto-trading/data-history/ltc"
	monerodatahistory "github.com/fabiodmferreira/crypto-trading/data-history/monero"
	stellardatahistory "github.com/fabiodmferreira/crypto-trading/data-history/stellar"
	xrpdatahistory "github.com/fabiodmferreira/crypto-trading/data-history/xrp"
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
		"ADA":      adadatahistory.LastYearMinute,
		"BTC":      btcdatahistory.Twenty19Current,
		"BTC-CASH": btccashdatahistory.LastYearMinute,
		"BTC-SV":   btcsvdatahistory.LastYearMinute,
		"EOS":      eosdatahistory.LastYearMinute,
		"ETC":      etcdatahistory.LastYearMinute,
		"ETH":      ethdatahistory.LastYearMinute,
		"LTC":      ltcdatahistory.LastYearMinute,
		"MONERO":   monerodatahistory.LastYearMinute,
		"STELLAR":  stellardatahistory.LastYearMinute,
		"XRP":      xrpdatahistory.LastYearMinute,
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
				"value": price,
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
