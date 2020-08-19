package webserver_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/fabiodmferreira/crypto-trading/benchmark"
	btcdatahistory "github.com/fabiodmferreira/crypto-trading/data-history/btc"
	"github.com/fabiodmferreira/crypto-trading/domain"
	"github.com/fabiodmferreira/crypto-trading/mocks"
	"github.com/fabiodmferreira/crypto-trading/webserver"
	"github.com/golang/mock/gomock"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/gorilla/mux"
)

func TestBenchmarkControllerGetDataSources(t *testing.T) {
	benchmarkController, benchmarkService, _ := NewBenchmarkController(t)

	req, err := http.NewRequest("GET", "/benchmark/data-sources", nil)

	if err != nil {
		t.Fatal(err)
	}

	rr := NewHttpResponse(http.HandlerFunc(benchmarkController.GetBenchmarkDataSourcesHandler), req)

	AssertResponseStatusCode(t, rr, http.StatusOK)

	if benchmarkService.GetDataSourcesCalls != 1 {
		t.Errorf("Expected benchmarkService.GetDataSources to have been called 1 time")
	}
}

func TestBenchmarkControllerGetBenchmarks(t *testing.T) {
	benchmarkController, benchmarkService, _ := NewBenchmarkController(t)

	req, err := http.NewRequest("GET", "/benchmark", nil)

	if err != nil {
		t.Fatal(err)
	}

	rr := NewHttpResponse(http.HandlerFunc(benchmarkController.GetBenchmarks), req)

	AssertResponseStatusCode(t, rr, http.StatusOK)

	if benchmarkService.FindAllCalls != 1 {
		t.Errorf("Expected benchmarkService.FindAllCalls to have been called 1 time")
	}
}

func TestBenchmarkControllerCreateBenchmark(t *testing.T) {
	benchmarkController, benchmarkService, _ := NewBenchmarkController(t)

	input := benchmark.Input{
		DecisionMakerOptions: domain.DecisionMakerOptions{MaximumBuyAmount: 0.1, MinimumProfitPerSold: 0.03, MinimumPriceDropToBuy: 0.01},
		StatisticsOptions:    domain.StatisticsOptions{NumberOfPointsHold: 200},
		CollectorOptions:     domain.CollectorOptions{PriceVariationDetection: 0.01, DataSource: nil},
		AccountInitialAmount: 2000,
		DataSourceFilePath:   btcdatahistory.May2020,
	}

	body, _ := json.Marshal(input)

	req, err := http.NewRequest("POST", "/benchmark", bytes.NewReader(body))

	if err != nil {
		t.Fatal(err)
	}

	rr := NewHttpResponse(http.HandlerFunc(benchmarkController.CreateBenchmark), req)

	AssertResponseStatusCode(t, rr, http.StatusCreated)

	if len(benchmarkService.CreateCalls) != 1 {
		t.Errorf("Expected benchmarkService.Create to have been called 1 time")
	}
}

func TestBenchmarkControllerDelete(t *testing.T) {
	benchmarkController, benchmarkService, _ := NewBenchmarkController(t)

	req, err := http.NewRequest("DELETE", "/benchmark", nil)

	if err != nil {
		t.Fatal(err)
	}

	rr := NewHttpResponse(http.HandlerFunc(benchmarkController.DeleteBenchmark), req)

	AssertResponseStatusCode(t, rr, http.StatusOK)

	if len(benchmarkService.DeleteByIDCalls) != 1 {
		t.Errorf("Expected benchmarkService.DeleteByID to have been called 1 time")
	}

	AssertRequestResponse(t, rr, "")
}

func TestBenchmarkControllerGetBenchmarkExecutionState(t *testing.T) {
	t.Run("should return 400 if no parameter is passed", func(t *testing.T) {
		benchmarkController, benchmarkService, _ := NewBenchmarkController(t)

		req, err := http.NewRequest("GET", "/benchmark/id/state", nil)

		if err != nil {
			t.Fatal(err)
		}

		rr := NewHttpResponse(http.HandlerFunc(benchmarkController.GetBenchmarkExecutionStateHandler), req)

		AssertResponseStatusCode(t, rr, http.StatusBadRequest)

		if len(benchmarkService.AggregateApplicationStateCalls) != 0 {
			t.Errorf("Expected benchmarkService.AggregateApplicationState to have been called 0 times")
		}

		AssertRequestResponse(t, rr, "startDate and endDate parameters are required")
	})

	t.Run("should return 200 if parameters start date and end date are passed", func(t *testing.T) {
		benchmarkController, _, applicationService := NewBenchmarkController(t)

		applicationService.
			EXPECT().
			GetStateAggregated(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(&[]bson.M{}, nil)

		params := url.Values{"startDate": {"2006-01-02T15:04:05"}, "endDate": {"2006-01-02T15:04:05"}}

		router := mux.NewRouter()
		router.HandleFunc("/benchmark/{id}/state", http.HandlerFunc(benchmarkController.GetBenchmarkExecutionStateHandler))

		ts := httptest.NewServer(router)
		defer ts.Close()

		url := ts.URL + "/benchmark/5f1e2010abf9760d6c686e32/state?" + params.Encode()
		rr, err := http.Get(url)

		if err != nil {
			t.Fatal(err)
		}

		if status := rr.StatusCode; status != http.StatusOK {
			t.Fatalf("wrong status code: got %d want %d", status, http.StatusOK)
		}

	})
}

func NewBenchmarkController(t *testing.T) (*webserver.BenchmarkController, *mocks.BenchmarkServiceSpy, *mocks.MockApplicationService) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	benchmarkService := &mocks.BenchmarkServiceSpy{}
	applicationService := mocks.NewMockApplicationService(ctrl)
	benchmarkController := webserver.NewBenchmarkController(benchmarkService, applicationService)

	return benchmarkController, benchmarkService, applicationService
}
