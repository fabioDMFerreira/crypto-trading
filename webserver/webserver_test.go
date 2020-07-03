package webserver

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fabiodmferreira/crypto-trading/benchmark"
)

func TestGetBenchmarkDataSources(t *testing.T) {
	res := MakeRequest(http.MethodGet, "/benchmark/data-sources")

	AssertRequestStatusCode(t, res, http.StatusOK)

	AssertRequestResponse(t, res, fmt.Sprintf("%v\n", `["btc/2019-november.csv","btc/2019-september-crash.csv","btc/2020-mar.csv","btc/2020-april.csv","btc/2020-may.csv","btc/last-year-minute.csv","ada/last-year-minute.csv","btc-cash/last-year-minute.csv","btc-sv/last-year-minute.csv","eos/last-year-minute.csv","etc/last-year-minute.csv","eth/last-year-minute.csv","ltc/last-year-minute.csv","monero/last-year-minute.csv","stellar/last-year-minute.csv","xrp/last-year-minute.csv"]`))
}

func TestGetBenchmarkList(t *testing.T) {
	res := MakeRequest(http.MethodGet, "/benchmark")

	AssertRequestStatusCode(t, res, http.StatusOK)

	AssertRequestResponse(t, res, "[]\n")
}

func TestDeleteBenchmarkResource(t *testing.T) {
	res := MakeRequest(http.MethodDelete, "/benchmark/random-id")

	AssertRequestStatusCode(t, res, http.StatusOK)

	AssertRequestResponse(t, res, "random-id")
}

func MakeRequest(method string, url string) *httptest.ResponseRecorder {
	repo := benchmark.NewRepositoryInMemory()
	benchmarkService := benchmark.NewService(repo)
	server, _ := NewCryptoTradingServer(benchmarkService)

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
