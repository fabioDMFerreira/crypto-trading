package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/fabiodmferreira/crypto-trading/assetsprices"
	"github.com/fabiodmferreira/crypto-trading/domain"
)

// WriteCoindeskResponse write remote response in write stream
func WriteCoindeskResponse(response domain.CoindeskResponse, out io.Writer) {
	for _, entry := range response.Entries {
		fmt.Fprintf(out, "%.0f,%f\n", entry[0], entry[1])
	}
}

// GetAndStoreData fetches remotes data and saves it in a file
func GetAndStoreData(f io.Writer, coin string, service *assetsprices.Service) func(time.Time, time.Time) error {
	counter := 0
	return func(startDate, endDate time.Time) error {
		response, err := service.GetRemotePrices(startDate, endDate, coin)

		if err != nil {
			return err
		}

		WriteCoindeskResponse(*response, f)

		counter++
		fmt.Printf("%v\r", startDate)

		return nil
	}
}

func main() {
	dateLayout := "2006-01-02 15:04:05"

	startDate, _ := time.Parse(dateLayout, "2019-01-01 00:00:00")
	endDate, _ := time.Parse(dateLayout, "2020-07-15 23:59:59")

	iterations := []struct {
		filePath string
		coin     string
	}{
		{"ada/last-year-minute.csv", "ADA"},
		{"btc/2019-current.csv", "BTC"},
		{"btc-cash/last-year-minute.csv", "BCH"},
		{"btc-sv/last-year-minute.csv", "BSV"},
		{"eos/last-year-minute.csv", "EOS"},
		{"etc/last-year-minute.csv", "ETC"},
		{"eth/last-year-minute.csv", "ETH"},
		{"ltc/last-year-minute.csv", "LTC"},
		{"monero/last-year-minute.csv", "XMR"},
		{"stellar/last-year-minute.csv", "XLM"},
		{"xrp/last-year-minute.csv", "XRP"},
	}

	assetsPricesRepository := assetsprices.NewRepositoryInMemory()
	assetsPricesService := assetsprices.NewService(assetsPricesRepository)

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
