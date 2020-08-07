package webserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/fabiodmferreira/crypto-trading/benchmark"
	"github.com/fabiodmferreira/crypto-trading/domain"
	"github.com/fabiodmferreira/crypto-trading/utils"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// BenchmarkController has the handlers of benchmark routes
type BenchmarkController struct {
	benchmark domain.BenchmarkService
}

// NewBenchmarkController returns an instance of BenchmarkController
func NewBenchmarkController(benchmark domain.BenchmarkService) *BenchmarkController {
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

	w.WriteHeader(http.StatusCreated)
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

// GetBenchmarkExecutionStateHandler returns state of application on each price change
func (b *BenchmarkController) GetBenchmarkExecutionStateHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	queryVars := r.URL.Query()

	if queryVars["startDate"] == nil || queryVars["endDate"] == nil || len(queryVars["startDate"]) == 0 || len(queryVars["endDate"]) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "startDate and endDate parameters are required")
		return
	}

	benchmarkID, err := primitive.ObjectIDFromHex(vars["id"])

	if err != nil {
		fmt.Fprint(w, "invalid benchmark id")
		return
	}

	// TODO: Validate query parameters.

	startDate, _ := time.Parse("2006-01-02T15:04:05", queryVars["startDate"][0])
	endDate, _ := time.Parse("2006-01-02T15:04:05", queryVars["endDate"][0])

	var pipelineOptions mongo.Pipeline

	groupByDatesClause := utils.GetGroupByDatesIDClause(startDate, endDate)

	pipelineOptions = mongo.Pipeline{
		{
			primitive.E{
				Key: "$match",
				Value: bson.M{
					"executionId": benchmarkID,
					"date":        bson.M{"$gte": startDate, "$lte": endDate}},
			},
		},
		{
			primitive.E{
				Key: "$group",
				Value: bson.M{
					"_id":                 groupByDatesClause,
					"average":             bson.M{"$avg": "$state.average"},
					"standardDeviation":   bson.M{"$avg": "$state.standardDeviation"},
					"higherBollingerBand": bson.M{"$avg": "$state.higherBollingerBand"},
					"lowerBollingerBand":  bson.M{"$avg": "$state.lowerBollingerBand"},
					"currentChange":       bson.M{"$avg": "$state.currentChange"},
					"accountAmount":       bson.M{"$avg": "$state.accountAmount"},
				},
			}},
	}

	benchmarkStates, err := b.benchmark.AggregateApplicationState(pipelineOptions)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	json.NewEncoder(w).Encode(*benchmarkStates)
}
