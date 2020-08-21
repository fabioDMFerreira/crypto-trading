package benchmark

import (
	"fmt"
	"path"
	"runtime"
	"time"

	"github.com/fabiodmferreira/crypto-trading/accounts"
	"github.com/fabiodmferreira/crypto-trading/app"
	"github.com/fabiodmferreira/crypto-trading/assets"
	"github.com/fabiodmferreira/crypto-trading/broker"
	"github.com/fabiodmferreira/crypto-trading/collectors"
	"github.com/fabiodmferreira/crypto-trading/decisionmaker"
	"github.com/fabiodmferreira/crypto-trading/domain"
	"github.com/fabiodmferreira/crypto-trading/statistics"
	"github.com/fabiodmferreira/crypto-trading/trader"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

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
)

// Input is an alias for BenchmarkInput
type Input = domain.BenchmarkInput

// Output is an alias for BenchmarkOutput
type Output = domain.BenchmarkOutput

// Service is a service with all methods to interact with benchmark related functions
type Service struct {
	repository                           domain.BenchmarksRepository
	assetpriceRepository                 domain.AssetPriceRepository
	applicationExecutionStatesRepository domain.ApplicationExecutionStateRepository
}

// NewService returns an instance of Service
func NewService(repo domain.BenchmarksRepository, assetpriceRepository domain.AssetPriceRepository, applicationExecutionStatesRepository domain.ApplicationExecutionStateRepository) *Service {
	return &Service{repo, assetpriceRepository, applicationExecutionStatesRepository}
}

// Create inserts one benchmark in database
func (s *Service) Create(input domain.BenchmarkInput) (*domain.Benchmark, error) {
	benchmark := &domain.Benchmark{ID: primitive.NewObjectID(), Input: input, Status: "Pending", CreatedAt: time.Now()}
	return benchmark, s.repository.InsertOne(benchmark)
}

// DeleteByID removes one benchmark from database
func (s *Service) DeleteByID(id string) error {
	err := s.repository.DeleteByID(id)

	if err != nil {
		return err
	}

	go s.applicationExecutionStatesRepository.BulkDeleteByExecutionID(id)

	return nil
}

// FindAll returns every benchmark
func (s *Service) FindAll() (*[]domain.Benchmark, error) {
	return s.repository.FindAll()
}

// BulkRun runs multiple benchmarks concurrently
func (s *Service) BulkRun(inputs []Input, c chan domain.BenchmarkResult) {
	for _, input := range inputs {
		go s.routineRun(&input, c)
	}
}

// ChannelService passes a benchmark output to a channel. Useful to run benchmarks in routines.
func (s *Service) routineRun(input *Input, done chan domain.BenchmarkResult) {
	result, err := s.Run(*input, nil)

	done <- domain.BenchmarkResult{Input: input, Output: result, Err: err}
}

// Run executes benchmark and returns performance results
func (s *Service) Run(input Input, benchmarkID *primitive.ObjectID) (*Output, error) {
	benchmarkApplication, err := s.setupApplication(input)

	if err != nil {
		return nil, err
	}

	var states []bson.M

	benchmarkApplication.RegistOnTickerChange(func(ask, bid float32, time time.Time) {
		if benchmarkID != nil {
			states = append(states, bson.M{
				"date":        time,
				"executionId": *benchmarkID,
				"state":       benchmarkApplication.GetState(),
			})

			if len(states) == 1000 {
				s.applicationExecutionStatesRepository.BulkCreate(&states)
				states = []bson.M{}
			}
		}
	})

	benchmarkApplication.Start()

	assetsDocs, _ := benchmarkApplication.FetchAssets()

	benchmarkAssetsInfo := assets.GroupAssetsByState(assetsDocs)

	amount, _ := benchmarkApplication.GetAccountAmount()

	output := domain.BenchmarkOutput{
		Buys:         benchmarkAssetsInfo.Buys,
		Sells:        benchmarkAssetsInfo.Sells,
		SellsPending: benchmarkAssetsInfo.SellsPending,
		FinalAmount:  amount,
		Assets:       assetsDocs,
	}

	return &output, nil
}

// setupApplication create the necessary application to run the benchmark
func (s *Service) setupApplication(input Input) (*app.App, error) {
	macd := statistics.NewMACDContainer(statistics.MACDParams{Fast: 12, Slow: 26, Lag: 9}, []float64{})
	statisticsService := statistics.NewStatistics(input.StatisticsOptions, macd)
	growthStatisticsService := statistics.NewStatistics(input.StatisticsOptions, macd)
	accelerationStatisticsService := statistics.NewStatistics(input.StatisticsOptions, macd)
	assetsRepository := &assets.AssetsRepositoryInMemory{}
	accountService := accounts.NewAccountServiceInMemory(float32(input.AccountInitialAmount), assetsRepository)

	decisionMakerOptions := input.DecisionMakerOptions
	decisionMaker := decisionmaker.NewDecisionMaker(accountService, decisionMakerOptions, statisticsService, growthStatisticsService, accelerationStatisticsService)

	_, currentFilePath, _, _ := runtime.Caller(0)
	currentDir := path.Dir(currentFilePath)
	historyFile, err := collectors.GetCsv(fmt.Sprintf("%v/../data-history/%v", currentDir, input.DataSourceFilePath))

	if err != nil {
		return nil, err
	}

	input.CollectorOptions.DataSource = historyFile

	collector := collectors.NewFileTickerCollector(input.CollectorOptions)

	broker := broker.NewBrokerMock()
	trader := trader.NewTrader(accountService, broker)

	application := app.NewApp(collector, decisionMaker, trader, accountService)

	return application, err
}

// HandleBenchmark executes benchmark and updates database accordingly
func (s *Service) HandleBenchmark(benchmark *domain.Benchmark) error {

	output, err := s.Run(benchmark.Input, &benchmark.ID)

	if err != nil {
		return err
	}

	// updates benchmark status and benchmark output
	err = s.repository.UpdateBenchmarkCompleted(benchmark.ID.Hex(), output)

	if err != nil {
		return err
	}

	return nil
}

// GetDataSources returns all available data sources
func (s *Service) GetDataSources() map[string]map[string]string {
	return map[string]map[string]string{
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
}

// AggregateApplicationState returns an aggregate of application state
func (s *Service) AggregateApplicationState(pipeline mongo.Pipeline) (*[]bson.M, error) {
	return s.applicationExecutionStatesRepository.Aggregate(pipeline)
}
