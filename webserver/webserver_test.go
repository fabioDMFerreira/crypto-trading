package webserver_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fabiodmferreira/crypto-trading/applicationExecutionStates"
	"github.com/fabiodmferreira/crypto-trading/assets"
	"github.com/fabiodmferreira/crypto-trading/assetsprices"
	"github.com/fabiodmferreira/crypto-trading/benchmark"
	btcdatahistory "github.com/fabiodmferreira/crypto-trading/data-history/btc"
	"github.com/fabiodmferreira/crypto-trading/domain"
	"github.com/fabiodmferreira/crypto-trading/mocks"
	"github.com/fabiodmferreira/crypto-trading/webserver"
	"github.com/golang/mock/gomock"
)

func TestGetBenchmarkDataSources(t *testing.T) {
	res := MakeRequest(t, http.MethodGet, "/api/benchmark/data-sources", nil)

	AssertResponseStatusCode(t, res, http.StatusOK)

	AssertRequestResponse(t, res, fmt.Sprintf("%v\n", `{"ada":{"2019":"ada/2019.csv","2020 H1":"ada/2020-h1.csv"},"btc":{"2019":"btc/2019.csv","2019-2020 H1":"btc/2019-2020-h1.csv","2020 H1":"btc/2020-h1.csv"},"btc-cash":{"2019":"btc-cash/2019.csv","2020 H1":"btc-cash/2020-h1.csv"},"eos":{"2019":"eos/2019.csv","2020 H1":"eos/2020-h1.csv"},"etc":{"2019":"etc/2019.csv","2020 H1":"etc/2020-h1.csv"},"eth":{"2019":"eth/2019.csv","2020 H1":"eth/2020-h1.csv"},"ltc":{"2019":"ltc/2019.csv","2020 H1":"ltc/2020-h1.csv"},"monero":{"2019":"monero/2019.csv","2020 H1":"monero/2020-h1.csv"},"stellar":{"2019":"stellar/2019.csv","2020 H1":"stellar/2020-h1.csv"},"xrp":{"2019":"xrp/2019.csv","2020 H1":"xrp/2020-h1.csv"}}`))
}

func TestGetBenchmarkList(t *testing.T) {
	res := MakeRequest(t, http.MethodGet, "/api/benchmark", nil)

	AssertResponseStatusCode(t, res, http.StatusOK)

	AssertRequestResponse(t, res, "[]\n")
}

func TestDeleteBenchmarkResource(t *testing.T) {
	res := MakeRequest(t, http.MethodDelete, "/api/benchmark/random-id", nil)

	AssertResponseStatusCode(t, res, http.StatusOK)

	AssertRequestResponse(t, res, "random-id")
}

func TestCreateBenchmarkResource(t *testing.T) {
	input := benchmark.Input{
		DecisionMakerOptions: domain.DecisionMakerOptions{MaximumBuyAmount: 0.1, MinimumProfitPerSold: 0.03, MinimumPriceDropToBuy: 0.01},
		StatisticsOptions:    domain.StatisticsOptions{NumberOfPointsHold: 200},
		CollectorOptions:     domain.CollectorOptions{PriceVariationDetection: 0.01, DataSource: nil},
		AccountInitialAmount: 2000,
		DataSourceFilePath:   btcdatahistory.March2020,
	}

	body, _ := json.Marshal(input)
	reader := bytes.NewReader(body)

	res := MakeRequest(t, http.MethodPost, "/api/benchmark", reader)

	AssertResponseStatusCode(t, res, http.StatusCreated)
}

func MakeRequest(t *testing.T, method string, url string, body *bytes.Reader) *httptest.ResponseRecorder {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	repo := benchmark.NewRepositoryInMemory()
	assetsPricesRepo := assetsprices.NewRepositoryInMemory()
	applicationExecutionsStatesRepo := applicationExecutionStates.NewRepositoryInMemory()
	accountsRepo := mocks.NewMockAccountsRepository(ctrl)
	appService := mocks.NewMockApplicationService(ctrl)
	assetsRepo := &assets.AssetsRepositoryInMemory{}
	benchmarkService := benchmark.NewService(repo, assetsPricesRepo, applicationExecutionsStatesRepo)
	server, _ := webserver.NewCryptoTradingServer(benchmarkService, assetsPricesRepo, accountsRepo, assetsRepo, appService)

	var req *http.Request

	if body != nil {
		req, _ = http.NewRequest(method, url, body)
	} else {
		req, _ = http.NewRequest(method, url, nil)
	}

	res := httptest.NewRecorder()

	server.ServeHTTP(res, req)

	return res
}

func AssertResponseStatusCode(t *testing.T, res *httptest.ResponseRecorder, want int) {
	t.Helper()

	got := res.Result().StatusCode

	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func AssertRequestResponse(t *testing.T, res *httptest.ResponseRecorder, want string) {
	t.Helper()

	got := res.Body.String()

	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func NewHttpResponse(handler http.HandlerFunc, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	return rr
}
