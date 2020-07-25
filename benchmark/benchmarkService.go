package benchmark

import (
	"fmt"
	"path"
	"runtime"
	"time"

	"github.com/fabiodmferreira/crypto-trading/accounts"
	"github.com/fabiodmferreira/crypto-trading/app"
	"github.com/fabiodmferreira/crypto-trading/assets"
	"github.com/fabiodmferreira/crypto-trading/assetsprices"
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

// BenchmarkResult stores benchmark returned value and possible error
type BenchmarkResult struct {
	Input  *Input
	Output *Output
	Err    error
}

type BenchmarkAssetsInfo struct {
	Buys         [][]float32
	Sells        [][]float32
	SellsPending int
}

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

	go s.applicationExecutionStatesRepository.BulkDelete(id)

	return nil
}

// FindAll returns every benchmark
func (s *Service) FindAll() (*[]domain.Benchmark, error) {
	return s.repository.FindAll()
}

// BulkRun runs multiple benchmarks concurrently
func (s *Service) BulkRun(inputs []Input, c chan BenchmarkResult) {
	for _, input := range inputs {
		go s.routineRun(&input, c)
	}
}

// ChannelService passes a benchmark output to a channel. Useful for run benchmarks in go routines.
func (s *Service) routineRun(input *Input, done chan BenchmarkResult) {
	result, err := s.Run(*input, nil)

	done <- BenchmarkResult{Input: input, Output: result, Err: err}
}

// Run executes benchmark and returns performance results
func (s *Service) Run(input Input, benchmarkID *primitive.ObjectID) (*Output, error) {
	benchmarkApplication, err := s.setupApplication(input)

	if err != nil {
		return nil, err
	}

	balances := [][]float32{}
	var currentAmount float32
	var states []bson.M

	benchmarkApplication.RegistOnTickerChange(func(ask, bid float32, time time.Time) {
		unixTime := float32(time.Unix()) * 1000
		amount, _ := benchmarkApplication.GetAccountAmount()
		if currentAmount != amount {
			balances = append(balances, []float32{unixTime, amount})
			currentAmount = amount
		}

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

	assets, _ := benchmarkApplication.FetchAssets()

	benchmarkAssetsInfo := s.getBenchmarkAssetsInfo(assets)

	amount, _ := benchmarkApplication.GetAccountAmount()

	output := domain.BenchmarkOutput{
		Buys:         benchmarkAssetsInfo.Buys,
		Sells:        benchmarkAssetsInfo.Sells,
		SellsPending: benchmarkAssetsInfo.SellsPending,
		FinalAmount:  amount,
		Assets:       assets,
		Balances:     balances,
	}

	return &output, nil
}

// getBenchmarkAssetsInfo transverse assets and generate info to be saved on benchmark output
func (s *Service) getBenchmarkAssetsInfo(assets *[]domain.Asset) BenchmarkAssetsInfo {
	var sells int

	Buys := [][]float32{}
	Sells := [][]float32{}

	for _, asset := range *assets {
		Buys = append(Buys, []float32{float32(asset.BuyTime.Unix()) * 1000, asset.BuyPrice})

		if asset.Sold {
			Sells = append(Sells, []float32{float32(asset.SellTime.Unix()) * 1000, asset.SellPrice})
			sells++
		}
	}

	return BenchmarkAssetsInfo{
		Buys:         Buys,
		Sells:        Sells,
		SellsPending: len(*assets) - sells,
	}
}

// setupApplication create the necessary application to run the benchmark
func (s *Service) setupApplication(input Input) (*app.App, error) {
	macd := statistics.NewMACDContainer(statistics.MACDParams{Fast: 12, Slow: 26, Lag: 9}, []float64{})
	statisticsService := statistics.NewStatistics(input.StatisticsOptions, macd)
	growthStatisticsService := statistics.NewStatistics(input.StatisticsOptions, macd)
	assetsRepository := &assets.AssetsRepositoryInMemory{}
	decisionMakerOptions := input.DecisionMakerOptions
	assetsPricesServiceMock := &assetsprices.ServiceMock{}
	decisionMaker := decisionmaker.NewDecisionMaker(assetsRepository, decisionMakerOptions, statisticsService, growthStatisticsService, assetsPricesServiceMock)

	_, currentFilePath, _, _ := runtime.Caller(0)
	currentDir := path.Dir(currentFilePath)
	historyFile, err := collectors.GetCsv(fmt.Sprintf("%v/../data-history/%v", currentDir, input.DataSourceFilePath))

	if err != nil {
		return nil, err
	}

	input.CollectorOptions.DataSource = historyFile

	collector := collectors.NewFileTickerCollector(input.CollectorOptions)

	accountService := accounts.NewAccountServiceInMemory(float32(input.AccountInitialAmount), assetsRepository)
	broker := broker.NewBrokerMock()
	trader := trader.NewTrader(assetsRepository, accountService, broker)

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

	// creates prices if they do not exist

	// creates benchmark balances

	// creates benchmark assets

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
		},
		"btc-sv": map[string]string{
			"Last Year Minute": btcsvdatahistory.LastYearMinute,
		},
		"eos": map[string]string{
			"Last Year Minute": eosdatahistory.LastYearMinute,
		},
		"etc": map[string]string{
			"Last Year Minute": etcdatahistory.LastYearMinute,
		},
		"ltc": map[string]string{
			"Last Year Minute": ltcdatahistory.LastYearMinute,
		},
		"monero": map[string]string{
			"Last Year Minute": monerodatahistory.LastYearMinute,
		},
		"stellar": map[string]string{
			"Last Year Minute": stellardatahistory.LastYearMinute,
		},
		"xrp": map[string]string{
			"Last Year Minute": xrpdatahistory.LastYearMinute,
		},
		"eth": map[string]string{
			"Last Year Minute": ethdatahistory.LastYearMinute,
		},
		"ada": map[string]string{
			"Last Year Minute": adadatahistory.LastYearMinute,
		},
	}
}

// AggregateApplicationState returns an aggregate of application state
func (s *Service) AggregateApplicationState(pipeline mongo.Pipeline) (*[]bson.M, error) {
	return s.applicationExecutionStatesRepository.Aggregate(pipeline)
}
