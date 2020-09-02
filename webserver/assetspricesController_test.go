package webserver_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/fabiodmferreira/crypto-trading/mocks"
	"github.com/fabiodmferreira/crypto-trading/webserver"
	"github.com/golang/mock/gomock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetAssetPrices(t *testing.T) {

	t.Run("should return 400 if no date paremeNter is passed", func(t *testing.T) {
		assetspricesController, _ := NewAssetsPricesController(t)

		req, err := http.NewRequest("GET", "/assets-prices", nil)

		if err != nil {
			t.Fatal(err)
		}

		rr := NewHttpResponse(NewGetAssetsPricesHandler(assetspricesController), req)

		AssertResponseStatusCode(t, rr, http.StatusBadRequest)
	})

	t.Run("should return 200 if start date and end date are passed as parameters", func(t *testing.T) {
		assetspricesController, assetspricesRepository := NewAssetsPricesController(t)

		assetspricesRepository.EXPECT().Aggregate(gomock.Any()).Return(&[]primitive.M{}, nil).Times(1)

		params := url.Values{"startDate": {"2006-01-02T15:04:05"}, "endDate": {"2006-01-02T15:04:05"}}

		req, err := http.NewRequest("GET", "/assets-prices?"+params.Encode(), nil)

		if err != nil {
			t.Fatal(err)
		}

		rr := NewHttpResponse(NewGetAssetsPricesHandler(assetspricesController), req)

		AssertResponseStatusCode(t, rr, http.StatusOK)
	})
}

func NewAssetsPricesController(t *testing.T) (*webserver.AssetsPricesController, *mocks.MockAssetPriceRepository) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	assetsPricesRepository := mocks.NewMockAssetPriceRepository(ctrl)
	assetspricesController := webserver.NewAssetsPricesController(assetsPricesRepository)

	return assetspricesController, assetsPricesRepository
}

func NewGetAssetsPricesHandler(assetspricesController *webserver.AssetsPricesController) http.HandlerFunc {
	return http.HandlerFunc(assetspricesController.GetAssetPrices)
}
