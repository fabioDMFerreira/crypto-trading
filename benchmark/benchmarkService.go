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
	var LastPrice float32

	benchmarkApplication.RegistOnNewAssetPrice(func(ohlc *domain.OHLC) {
		if benchmarkID != nil {
			states = append(states, bson.M{
				"date":        ohlc.Time,
				"executionId": *benchmarkID,
				"state":       benchmarkApplication.GetState(),
			})

			if len(states) == 1000 {
				s.applicationExecutionStatesRepository.BulkCreate(&states)
				states = []bson.M{}
			}

			LastPrice = ohlc.Close
		}
	})

	benchmarkApplication.Start()

	assetsDocs, _ := benchmarkApplication.FetchAssets()

	benchmarkAssetsInfo := assets.GroupAssetsByState(assetsDocs)

	amount, _ := benchmarkApplication.GetAccountAmount()

	output := domain.BenchmarkOutput{
		Buys:                benchmarkAssetsInfo.Buys,
		Sells:               benchmarkAssetsInfo.Sells,
		SellsPending:        benchmarkAssetsInfo.SellsPending,
		AssetsAmountPending: benchmarkAssetsInfo.AssetsAmountPending,
		AssetsValuePending:  LastPrice * benchmarkAssetsInfo.AssetsAmountPending,
		LastPrice:           LastPrice,
		FinalAmount:         amount,
		Assets:              assetsDocs,
	}

	return &output, nil
}

// setupApplication create the necessary application to run the benchmark
func (s *Service) setupApplication(input Input) (*app.App, error) {
	macd := statistics.NewMACDContainer(statistics.MACDParams{Fast: 12, Slow: 26, Lag: 9}, []float64{})
	statisticsService := statistics.NewStatistics(input.StatisticsOptions, macd)

	statisticsOptions := domain.StatisticsOptions{
		NumberOfPointsHold: input.StatisticsOptions.NumberOfPointsHold / 2,
	}

	growthStatisticsService := statistics.NewStatistics(statisticsOptions, macd)
	accelerationStatisticsService := statistics.NewStatistics(statisticsOptions, macd)
	volumeStatistics := statistics.NewStatistics(statisticsOptions, macd)

	assetsRepository := &assets.AssetsRepositoryInMemory{}
	accountService := accounts.NewAccountServiceInMemory(float32(input.AccountInitialAmount), assetsRepository)

	decisionMakerOptions := input.DecisionMakerOptions
	decisionMaker := decisionmaker.NewDecisionMaker(accountService, decisionMakerOptions, statisticsService, growthStatisticsService, accelerationStatisticsService, volumeStatistics)

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
}

// AggregateApplicationState returns an aggregate of application state
func (s *Service) AggregateApplicationState(pipeline mongo.Pipeline) (*[]bson.M, error) {
	return s.applicationExecutionStatesRepository.Aggregate(pipeline)
}
