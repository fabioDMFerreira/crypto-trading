package webserver

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fabiodmferreira/crypto-trading/applicationExecutionStates"
	"github.com/fabiodmferreira/crypto-trading/assetsprices"
	"github.com/fabiodmferreira/crypto-trading/benchmark"
)

func TestGetBenchmarkDataSources(t *testing.T) {
	res := MakeRequest(http.MethodGet, "/api/benchmark/data-sources")

	AssertRequestStatusCode(t, res, http.StatusOK)

	AssertRequestResponse(t, res, fmt.Sprintf("%v\n", `{"ada":{"Last Year Minute":"ada/last-year-minute.csv"},"btc":{"2019 - Current":"btc/2019-current.csv","Last Year Minute":"btc/last-year-minute.csv"},"btc-cash":{"Last Year Minute":"btc-cash/last-year-minute.csv"},"btc-sv":{"Last Year Minute":"btc-sv/last-year-minute.csv"},"eos":{"Last Year Minute":"eos/last-year-minute.csv"},"etc":{"Last Year Minute":"etc/last-year-minute.csv"},"eth":{"Last Year Minute":"eth/last-year-minute.csv"},"ltc":{"Last Year Minute":"ltc/last-year-minute.csv"},"monero":{"Last Year Minute":"monero/last-year-minute.csv"},"stellar":{"Last Year Minute":"stellar/last-year-minute.csv"},"xrp":{"Last Year Minute":"xrp/last-year-minute.csv"}}`))
}

func TestGetBenchmarkList(t *testing.T) {
	res := MakeRequest(http.MethodGet, "/api/benchmark")

	AssertRequestStatusCode(t, res, http.StatusOK)

	AssertRequestResponse(t, res, "[]\n")
}

func TestDeleteBenchmarkResource(t *testing.T) {
	res := MakeRequest(http.MethodDelete, "/api/benchmark/random-id")

	AssertRequestStatusCode(t, res, http.StatusOK)

	AssertRequestResponse(t, res, "random-id")
}

func MakeRequest(method string, url string) *httptest.ResponseRecorder {
	repo := benchmark.NewRepositoryInMemory()
	assetsPricesRepo := assetsprices.NewRepositoryInMemory()
	applicationExecutionsStatesRepo := applicationExecutionStates.NewRepositoryInMemory()
	benchmarkService := benchmark.NewService(repo, assetsPricesRepo, applicationExecutionsStatesRepo)
	server, _ := NewCryptoTradingServer(benchmarkService, assetsPricesRepo)

	req, _ := http.NewRequest(method, url, nil)
	res := httptest.NewRecorder()

	server.ServeHTTP(res, req)

	return res
}

func AssertRequestStatusCode(t *testing.T, res *httptest.ResponseRecorder, want int) {
	t.Helper()

	got := res.Result().StatusCode

	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func AssertRequestResponse(t *testing.T, res *httptest.ResponseRecorder, want string) {
	t.Helper()

	got := res.Body.String()

	fmt.Printf("'%v'\n'%v'\n", got, want)

	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}
