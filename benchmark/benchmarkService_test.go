package benchmark_test

import (
	"reflect"
	"testing"

	"github.com/fabiodmferreira/crypto-trading/benchmark"
	adadatahistory "github.com/fabiodmferreira/crypto-trading/data-history/ada"
	btcdatahistory "github.com/fabiodmferreira/crypto-trading/data-history/btc"
	btccashdatahistory "github.com/fabiodmferreira/crypto-trading/data-history/btc-cash"
	btcsvdatahistory "github.com/fabiodmferreira/crypto-trading/data-history/btc-sv"
	eosdatahistory "github.com/fabiodmferreira/crypto-trading/data-history/eos"
	etcdatahistory "github.com/fabiodmferreira/crypto-trading/data-history/etc"
	ethdatahistory "github.com/fabiodmferreira/crypto-trading/data-history/eth"
	ltcdatahistory "github.com/fabiodmferreira/crypto-trading/data-history/ltc"
	monerodatahistory "github.com/fabiodmferreira/crypto-trading/data-history/monero"
	stellardatahistory "github.com/fabiodmferreira/crypto-trading/data-history/stellar"
	xrpdatahistory "github.com/fabiodmferreira/crypto-trading/data-history/xrp"
	"github.com/fabiodmferreira/crypto-trading/domain"
	"github.com/fabiodmferreira/crypto-trading/mocks"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestBenchmarkServiceFindAll(t *testing.T) {
	service, benchmarkRepository, _, _ := NewBenchmarkService()

	service.FindAll()

	got := benchmarkRepository.FindAllCalls
	want := 1

	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestHandleBenchmark(t *testing.T) {
	service, benchmarkRepository, _, _ := NewBenchmarkService()

	benchmarkID := primitive.NewObjectID()

	benchmark := &domain.Benchmark{
		ID:    benchmarkID,
		Input: *NewBenchmarkInput(),
	}

	err := service.HandleBenchmark(benchmark)

	if err != nil {
		t.Errorf("Not expected Run to return error: %v", err)
	}

	if len(benchmarkRepository.UpdateBenchmarkCompletedCalls) != 1 {
		t.Errorf("Expected updateBenchmarkCompleted to be called 1 time")
	}

	got := benchmarkRepository.UpdateBenchmarkCompletedCalls[0].ID
	want := benchmarkID.Hex()

	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestServiceGetDatSources(t *testing.T) {
	service, _, _, _ := NewBenchmarkService()

	got := service.GetDataSources()
	want := map[string]map[string]string{
		"btc": map[string]string{
			"Last Year Minute": btcdatahistory.LastYearMinute,
			"2019 - Current":   btcdatahistory.Twenty19Current,
		},
		"btc-cash": map[string]string{
			"Last Year Minute": btccashdatahistory.LastYearMinute,
			"2019 - Current":   btccashdatahistory.Twenty19Current,
		},
		"btc-sv": map[string]string{
			"Last Year Minute": btcsvdatahistory.LastYearMinute,
			"2019 - Current":   btcsvdatahistory.Twenty19Current,
		},
		"eos": map[string]string{
			"Last Year Minute": eosdatahistory.LastYearMinute,
			"2019 - Current":   eosdatahistory.Twenty19Current,
		},
		"etc": map[string]string{
			"Last Year Minute": etcdatahistory.LastYearMinute,
			"2019 - Current":   etcdatahistory.Twenty19Current,
		},
		"ltc": map[string]string{
			"Last Year Minute": ltcdatahistory.LastYearMinute,
			"2019 - Current":   ltcdatahistory.Twenty19Current,
		},
		"monero": map[string]string{
			"Last Year Minute": monerodatahistory.LastYearMinute,
			"2019 - Current":   monerodatahistory.Twenty19Current,
		},
		"stellar": map[string]string{
			"Last Year Minute": stellardatahistory.LastYearMinute,
			"2019 - Current":   stellardatahistory.Twenty19Current,
		},
		"xrp": map[string]string{
			"Last Year Minute": xrpdatahistory.LastYearMinute,
			"2019 - Current":   xrpdatahistory.Twenty19Current,
		},
		"eth": map[string]string{
			"Last Year Minute": ethdatahistory.LastYearMinute,
			"2019 - Current":   ethdatahistory.Twenty19Current,
		},
		"ada": map[string]string{
			"Last Year Minute": adadatahistory.LastYearMinute,
			"2019 - Current":   adadatahistory.Twenty19Current,
		},
	}

	if reflect.DeepEqual(got, want) != true {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestServiceAggregateApplicationState(t *testing.T) {
	service, _, _, applicationExecutionStatesRepository := NewBenchmarkService()

	service.AggregateApplicationState(mongo.Pipeline{})

	if len(applicationExecutionStatesRepository.AggregateCalls) != 1 {
		t.Errorf("Expected ApplicationExecutionStatesRepository.Aggregate to be called 1 time")
	}
}

func TestServiceCreate(t *testing.T) {
	service, benchmarkRepository, _, _ := NewBenchmarkService()

	service.Create(*NewBenchmarkInput())

	if len(benchmarkRepository.InsertOneCalls) != 1 {
		t.Errorf("Expected BenchmarkRepository.InsertOne to be called 1 time")
	}
}

func TestServiceDeletedById(t *testing.T) {
	service, benchmarkRepository, _, _ := NewBenchmarkService()

	ID := "test-id-1"

	err := service.DeleteByID(ID)

	if err != nil {
		t.Errorf("Not expected DeleteByID to throw error: %v", err)
	}

	got := benchmarkRepository.DeleteByIdCalls
	want := []string{ID}

	if reflect.DeepEqual(got, want) != true {
		t.Errorf("BenchmarkRepositoy.DeletedById: got %v want %v", got, want)
	}

	// TODO: find a way to wait for go routine inside DeleteByID to finish
	// got2 := applicationExecutionStatesRepository.BulkDeleteCalls
	// want2 := []string{ID}

	// if reflect.DeepEqual(got2, want2) != true {
	// 	t.Errorf("ApplicationExecutionStatesRepository.BulkDelete: got %v want %v", got2, want2)
	// }
}

func NewBenchmarkInput() *domain.BenchmarkInput {
	return &domain.BenchmarkInput{
		DecisionMakerOptions: domain.DecisionMakerOptions{MaximumBuyAmount: 0.1, MinimumProfitPerSold: 0.1, MinimumPriceDropToBuy: 0.02},
		StatisticsOptions:    domain.StatisticsOptions{NumberOfPointsHold: 10},
		CollectorOptions:     domain.CollectorOptions{PriceVariationDetection: 0.1},
		AccountInitialAmount: 5000,
		DataSourceFilePath:   btcdatahistory.March2020,
	}
}

func NewBenchmarkService() (*benchmark.Service, *mocks.BenchmarkRepositorySpy, *mocks.AssetPriceRepositorySpy, *mocks.ApplicationExecutionStatesRepositorySpy) {
	repository := &mocks.BenchmarkRepositorySpy{}
	assetPriceRepository := &mocks.AssetPriceRepositorySpy{}
	applicationExecutionStatesRepository := &mocks.ApplicationExecutionStatesRepositorySpy{}

	return benchmark.NewService(repository, assetPriceRepository, applicationExecutionStatesRepository), repository, assetPriceRepository, applicationExecutionStatesRepository
}
