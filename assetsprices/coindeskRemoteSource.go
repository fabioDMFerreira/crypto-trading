package assetsprices

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/fabiodmferreira/crypto-trading/domain"
	"github.com/fabiodmferreira/crypto-trading/utils"
	"go.mongodb.org/mongo-driver/bson"
)

// FetchRemotePrices uses remote source to get asset prices
func FetchCoindeskRemotePrices(startDate, endDate time.Time, asset string) (*[]bson.M, error) {
	response := domain.CoindeskHTTPResponse{}
	err := fetchCoindeskData(SerializeDate(startDate), SerializeDate(endDate), asset, &response)

	time.Sleep(2 * time.Second)

	if err != nil {
		return nil, err
	}

	var assetsPrices []bson.M

	for _, entry := range response.Data.Entries {
		assetsPrices = append(assetsPrices,
			bson.M{
				"asset": asset,
				"date":  time.Unix(int64(entry[0])/1000, 0),
				"value": float32(entry[1]) * utils.DollarEuroRate,
			})
	}

	return &assetsPrices, nil
}

// fetchCoindeskData uses public API provided by Coindesk
func fetchCoindeskData(startDate, endDate string, coin string, target *domain.CoindeskHTTPResponse) error {
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
