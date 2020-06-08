package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func GetNextDateParams(startDate, endDate time.Time) (time.Time, time.Time) {
	return startDate.AddDate(0, 0, 1), endDate.AddDate(0, 0, 1)
}

func SerializeDate(date time.Time) string {
	return strings.Replace(date.Format("2006-01-02 15:04"), " ", "T", 1)
}

func TransverseDatesRanges(startDate, endDate time.Time, handle func(time.Time, time.Time) error) {
	startDateCursor := startDate
	endDateCursor := startDateCursor.Add(
		23*time.Hour +
			59*time.Minute +
			59*time.Second)

	for startDateCursor.Before(endDate) {
		handle(startDateCursor, endDateCursor)
		startDateCursor, endDateCursor = GetNextDateParams(startDateCursor, endDateCursor)
	}
}

func FetchCoindeskData(startDate, endDate string, coin string, target *CoindeskHTTPResponse) error {
	r, err := http.Get(fmt.Sprintf("https://production.api.coindesk.com/v2/price/values/%v?start_date=%v&end_date=%v&ohlc=false", coin, startDate, endDate))

	if err != nil {
		return err
	}
	defer r.Body.Close()

	body, readErr := ioutil.ReadAll(r.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	// var data map[string]interface{}
	return json.Unmarshal(body, &target)
}

type CoindeskResponse struct {
	Iso      string      `json:"iso"`
	Name     string      `json:"name"`
	Slug     string      `json:"slug"`
	Interval string      `json:"interval"`
	Entries  [][]float64 `json:"entries"`
}

type CoindeskHTTPResponse struct {
	StatusCode int              `json:"statusCode"`
	Message    string           `json:"message"`
	Data       CoindeskResponse `json:"data"`
}

func WriteCoindeskResponse(response CoindeskResponse, out io.Writer) {
	for _, entry := range response.Entries {
		fmt.Fprintf(out, "%.0f,%f\n", entry[0], entry[1])
	}
}

func GetPrices(startDate, endDate time.Time, coin string) (*CoindeskResponse, error) {

	response := CoindeskHTTPResponse{}
	err := FetchCoindeskData(SerializeDate(startDate), SerializeDate(endDate), coin, &response)

	time.Sleep(2 * time.Second)

	if err != nil {
		return nil, err
	}

	return &response.Data, nil
}

func GetAndStoreData(f io.Writer, coin string) func(time.Time, time.Time) error {
	counter := 0
	return func(startDate, endDate time.Time) error {
		response, err := GetPrices(startDate, endDate, coin)

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

	startDate, _ := time.Parse(dateLayout, "2019-06-01 00:00:00")
	endDate, _ := time.Parse(dateLayout, "2020-06-01 23:59:59")

	iterations := []struct {
		filePath string
		coin     string
	}{
		// {"ada/last-year-minute.csv", "ADA"},
		// {"btc/last-year-minute.csv", "BTC"},
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

	for _, i := range iterations {
		fmt.Printf("\nfetching %v...\n", i.coin)
		f, err := os.Create(fmt.Sprintf("./data-history/%v", i.filePath))

		if err != nil {
			log.Fatal(err)
		}

		f.Write([]byte("Date,Price\n"))

		TransverseDatesRanges(startDate, endDate, GetAndStoreData(f, i.coin))

		if err != nil {
			log.Fatal(err)
		}
	}

}
