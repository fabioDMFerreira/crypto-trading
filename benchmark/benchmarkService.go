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
	"github.com/fabiodmferreira/crypto-trading/eventlogs"
	"github.com/fabiodmferreira/crypto-trading/notifications"
	"github.com/fabiodmferreira/crypto-trading/statistics"
	"github.com/fabiodmferreira/crypto-trading/trader"
	"go.mongodb.org/mongo-driver/bson/primitive"

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

// Service is a service with all methods to interact with benchmark related functions
type Service struct {
	repository           domain.BenchmarksRepository
	assetpriceRepository domain.AssetPriceRepository
}

// NewService returns an instance of Service
func NewService(repo domain.BenchmarksRepository, assetpriceRepository domain.AssetPriceRepository) *Service {
	return &Service{repo, assetpriceRepository}
}

// Create inserts one benchmark in database
func (s *Service) Create(input domain.BenchmarkInput) (*domain.Benchmark, error) {
	benchmark := &domain.Benchmark{ID: primitive.NewObjectID(), Input: input, Status: "Pending", CreatedAt: time.Now()}
	return benchmark, s.repository.InsertOne(benchmark)
}

// DeleteByID removes one benchmark from database
func (s *Service) DeleteByID(id string) error {
	return s.repository.DeleteByID(id)
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
	result, err := s.Run(*input)

	done <- BenchmarkResult{Input: input, Output: result, Err: err}
}

// Run executes benchmark and returns performance results
func (s *Service) Run(input Input) (*Output, error) {
	benchmarkApplication, err := s.setupApplication(input)

	if err != nil {
		return nil, err
	}

	balances := [][]float32{}
	var currentAmount float32

	benchmarkApplication.RegistOnTickerChange(func(ask, bid float32, time time.Time) {
		unixTime := float32(time.Unix()) * 1000
		amount, _ := benchmarkApplication.GetAccountAmount()
		if currentAmount != amount {
			balances = append(balances, []float32{unixTime, amount})
			currentAmount = amount
		}
	})

	benchmarkApplication.Start()

	var sells int
	assets, _ := benchmarkApplication.FetchAssets()
	for _, asset := range *assets {
		if asset.Sold {
			sells++
		} else {
		}
	}

	Buys := [][]float32{}
	Sells := [][]float32{}

	for _, asset := range *assets {
		Buys = append(Buys, []float32{float32(asset.BuyTime.Unix()) * 1000, asset.BuyPrice})

		if asset.Sold {
			Sells = append(Sells, []float32{float32(asset.SellTime.Unix()) * 1000, asset.SellPrice})
		}
	}

	amount, _ := benchmarkApplication.GetAccountAmount()

	output := domain.BenchmarkOutput{
		Buys:         Buys,
		Sells:        Sells,
		SellsPending: len(*assets) - sells,
		FinalAmount:  amount,
		Assets:       assets,
		Balances:     balances,
	}

	return &output, nil
}

// setupApplication create the necessary application to run the benchmark
func (s *Service) setupApplication(input Input) (*app.App, error) {
	macd := statistics.NewMACDContainer(statistics.MACDParams{Fast: 12, Slow: 26, Lag: 9}, []float64{})
	statisticsService := statistics.NewStatistics(input.StatisticsOptions, macd)
	assetsRepository := &assets.AssetsRepositoryInMemory{}
	decisionMakerOptions := input.DecisionMakerOptions
	decisionMaker := decisionmaker.NewDecisionMaker(assetsRepository, decisionMakerOptions, statisticsService)

	_, currentFilePath, _, _ := runtime.Caller(0)
	currentDir := path.Dir(currentFilePath)
	historyFile, err := collectors.GetCsv(fmt.Sprintf("%v/../data-history/%v", currentDir, input.DataSourceFilePath))

	if err != nil {
		return nil, err
	}

	input.CollectorOptions.DataSource = historyFile

	collector := collectors.NewFileTickerCollector(input.CollectorOptions)

	notificationsService := &notifications.NotificationsMock{}
	logService := &eventlogs.EventLogsServiceMock{}
	accountService := accounts.NewAccountServiceInMemory(float32(input.AccountInitialAmount))
	broker := broker.NewBrokerMock()
	trader := trader.NewTrader(assetsRepository, accountService, broker)

	application := app.NewApp(notificationsService, decisionMaker, logService, assetsRepository, trader, accountService, collector)

	return application, err
}

// HandleBenchmark executes benchmark and updates database accordingly
func (s *Service) HandleBenchmark(benchmark *domain.Benchmark) error {

	output, err := s.Run(benchmark.Input)

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
func (s *Service) GetDataSources() []string {
	return []string{
		btcdatahistory.November2019,
		btcdatahistory.SeptemberCrash2019,
		btcdatahistory.March2020,
		btcdatahistory.April2020,
		btcdatahistory.May2020,
		btcdatahistory.LastYearMinute,
		fmt.Sprintf("btc-cash/%v", btccashdatahistory.LastYearMinute),
		fmt.Sprintf("btc-sv/%v", btcsvdatahistory.LastYearMinute),
		fmt.Sprintf("eos/%v", eosdatahistory.LastYearMinute),
		fmt.Sprintf("etc/%v", etcdatahistory.LastYearMinute),
		fmt.Sprintf("eth/%v", ethdatahistory.LastYearMinute),
		fmt.Sprintf("ltc/%v", ltcdatahistory.LastYearMinute),
		fmt.Sprintf("monero/%v", monerodatahistory.LastYearMinute),
		fmt.Sprintf("stellar/%v", stellardatahistory.LastYearMinute),
		fmt.Sprintf("xrp/%v", xrpdatahistory.LastYearMinute),
	}
}
