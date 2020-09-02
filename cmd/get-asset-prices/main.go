package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/fabiodmferreira/crypto-trading/assetsprices"
)

// GetAndStoreData fetches remotes data and saves it in a file
func GetAndStoreData(f io.Writer, coin string, service *assetsprices.Service) func(time.Time, time.Time) error {
	counter := 0
	return func(startDate, endDate time.Time) error {
		assetsprices, err := service.FetchRemotePrices(startDate, endDate, coin)

		if err != nil {
			return err
		}

		for _, entry := range *assetsprices {
			fmt.Fprintf(f, "%d,%f\n", entry["date"].(time.Time).Unix(), entry["value"])
		}

		counter++
		fmt.Printf("%v\r", startDate)

		return nil
	}
}

func main() {
	dateLayout := "2006-01-02 15:04:05"

	startDate, _ := time.Parse(dateLayout, "2019-01-01 00:00:00")
	endDate, _ := time.Parse(dateLayout, "2020-08-15 23:59:59")

	iterations := []struct {
		filePath string
		coin     string
	}{
		{"ada/2019-current.csv", "ADA"},
		{"btc/2019-current.csv", "BTC"},
		{"btc-cash/2019-current.csv", "BCH"},
		{"eos/2019-current.csv", "EOS"},
		{"etc/2019-current.csv", "ETC"},
		{"eth/2019-current.csv", "ETH"},
		{"ltc/2019-current.csv", "LTC"},
		{"monero/2019-current.csv", "XMR"},
		{"stellar/2019-current.csv", "XLM"},
		{"xrp/2019-current.csv", "XRP"},
	}

	assetsPricesRepository := assetsprices.NewRepositoryInMemory()
	assetsPricesService := assetsprices.NewService(assetsPricesRepository, assetsprices.NewCoindeskRemoteSource(http.Get).FetchRemoteAssetsPrices)

	for _, i := range iterations {
		fmt.Printf("\nfetching %v...\n", i.coin)
		f, err := os.Create(fmt.Sprintf("./data-history/%v", i.filePath))

		if err != nil {
			log.Fatal(err)
		}

		f.Write([]byte("Date,Price\n"))

		assetsprices.TransverseDatesRange(startDate, endDate, GetAndStoreData(f, i.coin, assetsPricesService))

		if err != nil {
			log.Fatal(err)
		}
	}

}
