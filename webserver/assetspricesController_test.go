package webserver_test

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/fabiodmferreira/crypto-trading/mocks"
	"github.com/fabiodmferreira/crypto-trading/webserver"
)

func TestGetAssetPrices(t *testing.T) {

	t.Run("should return 400 if no date paremeNter is passed", func(t *testing.T) {
		assetspricesController, _ := NewAssetsPricesController()

		req, err := http.NewRequest("GET", "/assets-prices", nil)

		if err != nil {
			t.Fatal(err)
		}

		rr := NewHttpResponse(NewGetAssetsPricesHandler(assetspricesController), req)

		AssertResponseStatusCode(t, rr, http.StatusBadRequest)
	})

	t.Run("should return 200 if start date and end date are passed as parameters", func(t *testing.T) {
		assetspricesController, assetspricesRepository := NewAssetsPricesController()

		params := url.Values{"startDate": {"2006-01-02T15:04:05"}, "endDate": {"2006-01-02T15:04:05"}}

		req, err := http.NewRequest("GET", "/assets-prices?"+params.Encode(), nil)

		if err != nil {
			t.Fatal(err)
		}

		rr := NewHttpResponse(NewGetAssetsPricesHandler(assetspricesController), req)

		AssertResponseStatusCode(t, rr, http.StatusOK)

		got := len(assetspricesRepository.AggregateCalls)
		want := 1

		if reflect.DeepEqual(got, want) != true {
			t.Errorf("got %v want %v", got, want)
		}
	})
}

func NewAssetsPricesController() (*webserver.AssetsPricesController, *mocks.AssetPriceRepositorySpy) {
	assetsRepository := &mocks.AssetPriceRepositorySpy{}
	assetspricesController := webserver.NewAssetsPricesController(assetsRepository)

	return assetspricesController, assetsRepository
}

func NewGetAssetsPricesHandler(assetspricesController *webserver.AssetsPricesController) http.HandlerFunc {
	return http.HandlerFunc(assetspricesController.GetAssetPrices)
}
