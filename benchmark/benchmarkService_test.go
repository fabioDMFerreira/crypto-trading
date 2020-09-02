package benchmark_test

import (
	"reflect"
	"testing"

	"github.com/fabiodmferreira/crypto-trading/benchmark"
	adadatahistory "github.com/fabiodmferreira/crypto-trading/data-history/ada"
	btcdatahistory "github.com/fabiodmferreira/crypto-trading/data-history/btc"
	btccashdatahistory "github.com/fabiodmferreira/crypto-trading/data-history/btc-cash"
	eosdatahistory "github.com/fabiodmferreira/crypto-trading/data-history/eos"
	etcdatahistory "github.com/fabiodmferreira/crypto-trading/data-history/etc"
	ethdatahistory "github.com/fabiodmferreira/crypto-trading/data-history/eth"
	ltcdatahistory "github.com/fabiodmferreira/crypto-trading/data-history/ltc"
	monerodatahistory "github.com/fabiodmferreira/crypto-trading/data-history/monero"
	stellardatahistory "github.com/fabiodmferreira/crypto-trading/data-history/stellar"
	xrpdatahistory "github.com/fabiodmferreira/crypto-trading/data-history/xrp"
	"github.com/fabiodmferreira/crypto-trading/domain"
	"github.com/fabiodmferreira/crypto-trading/mocks"
	"github.com/golang/mock/gomock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestBenchmarkServiceFindAll(t *testing.T) {
	service, benchmarkRepository, _, _ := NewBenchmarkService(t)

	service.FindAll()

	got := benchmarkRepository.FindAllCalls
	want := 1

	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestHandleBenchmark(t *testing.T) {
	service, benchmarkRepository, _, _ := NewBenchmarkService(t)

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
	service, _, _, _ := NewBenchmarkService(t)

	got := service.GetDataSources()
	want := map[string]map[string]string{
		"btc": {
			"2019":         btcdatahistory.Twenty19,
			"2020 H1":      btcdatahistory.TwentyTwentyH1,
			"2019-2020 H1": btcdatahistory.Twenty1920H1,
		},
		"btc-cash": {
			"2019":    btccashdatahistory.Twenty19,
			"2020 H1": btccashdatahistory.TwentyTwentyH1,
		},
		"eos": {
			"2019":    eosdatahistory.Twenty19,
			"2020 H1": eosdatahistory.TwentyTwentyH1,
		},
		"etc": {
			"2019":    etcdatahistory.Twenty19,
			"2020 H1": etcdatahistory.TwentyTwentyH1,
		},
		"ltc": {
			"2019":    ltcdatahistory.Twenty19,
			"2020 H1": ltcdatahistory.TwentyTwentyH1,
		},
		"monero": {
			"2019":    monerodatahistory.Twenty19,
			"2020 H1": monerodatahistory.TwentyTwentyH1,
		},
		"stellar": {
			"2019":    stellardatahistory.Twenty19,
			"2020 H1": stellardatahistory.TwentyTwentyH1,
		},
		"xrp": {
			"2019":    xrpdatahistory.Twenty19,
			"2020 H1": xrpdatahistory.TwentyTwentyH1,
		},
		"eth": {
			"2019":    ethdatahistory.Twenty19,
			"2020 H1": ethdatahistory.TwentyTwentyH1,
		},
		"ada": {
			"2019":    adadatahistory.Twenty19,
			"2020 H1": adadatahistory.TwentyTwentyH1,
		},
	}

	if reflect.DeepEqual(got, want) != true {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestServiceAggregateApplicationState(t *testing.T) {
	service, _, _, applicationExecutionStatesRepository := NewBenchmarkService(t)

	service.AggregateApplicationState(mongo.Pipeline{})

	if len(applicationExecutionStatesRepository.AggregateCalls) != 1 {
		t.Errorf("Expected ApplicationExecutionStatesRepository.Aggregate to be called 1 time")
	}
}

func TestServiceCreate(t *testing.T) {
	service, benchmarkRepository, _, _ := NewBenchmarkService(t)

	service.Create(*NewBenchmarkInput())

	if len(benchmarkRepository.InsertOneCalls) != 1 {
		t.Errorf("Expected BenchmarkRepository.InsertOne to be called 1 time")
	}
}

func TestServiceDeletedById(t *testing.T) {
	service, benchmarkRepository, _, _ := NewBenchmarkService(t)

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

func NewBenchmarkService(t *testing.T) (*benchmark.Service, *mocks.BenchmarkRepositorySpy, *mocks.MockAssetPriceRepository, *mocks.ApplicationExecutionStatesRepositorySpy) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	repository := &mocks.BenchmarkRepositorySpy{}
	assetsPriceRepo := mocks.NewMockAssetPriceRepository(ctrl)
	applicationExecutionStatesRepository := &mocks.ApplicationExecutionStatesRepositorySpy{}

	return benchmark.NewService(repository, assetsPriceRepo, applicationExecutionStatesRepository), repository, assetsPriceRepo, applicationExecutionStatesRepository
}
