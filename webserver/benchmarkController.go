package webserver

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fabiodmferreira/crypto-trading/benchmark"
	"github.com/fabiodmferreira/crypto-trading/domain"
	"github.com/gorilla/mux"
)

// BenchmarkController has the handlers of benchmark routes
type BenchmarkController struct {
	benchmark *benchmark.Service
}

// NewBenchmarkController returns an instance of BenchmarkController
func NewBenchmarkController(benchmark *benchmark.Service) *BenchmarkController {
	return &BenchmarkController{benchmark}
}

// GetBenchmarkDataSources returns list of all available data sources
func (b *BenchmarkController) GetBenchmarkDataSourcesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(b.benchmark.GetDataSources())
}

// BenchmarkHandler handles benchmark routes
func (b *BenchmarkController) BenchmarkHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		b.CreateBenchmark(w, r)
	case http.MethodGet:
		b.GetBenchmarks(w, r)
	}
}

// ResourceHandler handles benchmark routes releated with a benchmark result
func (b *BenchmarkController) ResourceHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		b.DeleteBenchmark(w, r)
	case http.MethodDelete:
		b.DeleteBenchmark(w, r)
	}
}

// BenchmarkAPIOutput is a easy format to be read by charts
type BenchmarkAPIOutput struct {
	Prices  [][]float32    `json:"prices"`
	Balance [][]float32    `json:"balance"`
	Buys    [][]float32    `json:"buys"`
	Sells   [][]float32    `json:"sells"`
	Assets  []domain.Asset `json:"assets"`
}

// GetBenchmarks returns all existing benchmarks
func (b *BenchmarkController) GetBenchmarks(w http.ResponseWriter, r *http.Request) {
	benchmarks, err := b.benchmark.FindAll()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(benchmarks)
}

// CreateBenchmark handles request for creating benchmark in database and starting the benchmark execution
func (b *BenchmarkController) CreateBenchmark(w http.ResponseWriter, r *http.Request) {
	var input benchmark.Input

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// input := benchmark.Input{
	// 	DecisionMakerOptions: decisionmaker.Options{MaximumBuyAmount: 0.1, MinimumProfitPerSold: 0.03, MinimumPriceDropToBuy: 0.01},
	// 	StatisticsOptions:    statistics.Options{NumberOfPointsHold: 200},
	// 	CollectorOptions:     collectors.Options{0.01, nil},
	// 	AccountInitialAmount: 2000,
	// 	DataSourceFilePath:   btcdatahistory.May2020,
	// }

	benchmark, err := b.benchmark.Create(input)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(benchmark)

	go b.benchmark.HandleBenchmark(benchmark)
	// if err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	fmt.Fprint(w, err)
	// } else {
	// 	var response BenchmarkAPIOutput

	// 	for _, asset := range *output.Assets {
	// 		response.Buys = append(response.Buys, []float32{float32(asset.BuyTime.Unix()) * 1000, asset.BuyPrice})

	// 		if asset.Sold {
	// 			response.Sells = append(response.Sells, []float32{float32(asset.SellTime.Unix()) * 1000, asset.SellPrice})
	// 		}
	// 	}

	// 	response.Prices = output.Prices
	// 	response.Balance = output.Balances
	// 	response.Assets = *output.Assets

	// 	w.Header().Set("content-type", "application/json")
	// 	json.NewEncoder(w).Encode(response)
	// }
}

// DeleteBenchmark handles request for deleting benchmark
func (b *BenchmarkController) DeleteBenchmark(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	err := b.benchmark.DeleteByID(vars["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, vars["id"])
}
