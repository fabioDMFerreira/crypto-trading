package main

import (
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/fabiodmferreira/crypto-trading/app"
	"github.com/fabiodmferreira/crypto-trading/broker"
	"github.com/fabiodmferreira/crypto-trading/collectors"
	btcdatahistory "github.com/fabiodmferreira/crypto-trading/data-history/btc"
	"github.com/fabiodmferreira/crypto-trading/decisionmaker"
	"github.com/fabiodmferreira/crypto-trading/domain"
	"github.com/fabiodmferreira/crypto-trading/statistics"
	"github.com/fabiodmferreira/crypto-trading/trader"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotificationsMock struct {
	emailsNotifications int
}

func (n *NotificationsMock) CreateEmailNotification(subject, message, notificationType string) error {
	n.emailsNotifications++
	return nil
}

func (n *NotificationsMock) FindLastEventLogsNotificationDate() (time.Time, error) {
	return time.Now(), nil
}

type LogMock struct {
	logs [][]string
}

func (l *LogMock) Create(logType, message string) error {
	l.logs = append(l.logs, []string{logType, message})
	return nil
}

func (l *LogMock) FindAllToNotify() (*[]domain.EventLog, error) {
	return &[]domain.EventLog{}, nil
}

func (l *LogMock) MarkNotified(ids []primitive.ObjectID) error {
	return nil
}

type AccountServiceMock struct {
	Amount    float32
	withdraws int
	deposits  int
}

func (a *AccountServiceMock) Deposit(amount float32) error {
	a.deposits++
	a.Amount += amount
	// if a.Amount < 5000 {
	// 	fmt.Printf("%v ", a.Amount)
	// }
	return nil
}

func (a *AccountServiceMock) Withdraw(amount float32) error {
	if amount > a.Amount {
		return errors.New("Insufficient Funds")
	}

	a.Amount -= amount
	a.withdraws++
	// if a.Amount < 5000 {
	// 	fmt.Printf("%v ", a.Amount)
	// }

	return nil
}

func (a *AccountServiceMock) GetAmount() (float32, error) {
	return a.Amount, nil
}

type AssetsRepositoryMock struct {
	Assets []domain.Asset
}

func (ar *AssetsRepositoryMock) FindAll() (*[]domain.Asset, error) {

	pendingAssets := []domain.Asset{}

	for _, asset := range ar.Assets {
		if !asset.Sold {
			pendingAssets = append(pendingAssets, asset)
		}
	}
	return &pendingAssets, nil
}

func (ar *AssetsRepositoryMock) FindCheaperAssetPrice() (float32, error) {
	var minimumPrice float32

	for _, asset := range ar.Assets {
		if asset.Sold == false && minimumPrice > asset.BuyPrice {
			minimumPrice = asset.BuyPrice
		}
	}

	return minimumPrice, nil
}

func (ar *AssetsRepositoryMock) GetBalance(startDate, endDate time.Time) (float32, error) {
	return 0, nil
}

func (ar *AssetsRepositoryMock) Create(asset *domain.Asset) error {
	ar.Assets = append(ar.Assets, *asset)
	return nil
}

func (ar *AssetsRepositoryMock) Sell(id primitive.ObjectID, price float32) error {

	for index, asset := range ar.Assets {
		if asset.ID == id {
			ar.Assets[index].SellPrice = price
			ar.Assets[index].Sold = true
			ar.Assets[index].SellTime = time.Now()
			break
		}
	}

	return nil
}

type Algo0BenchmarkInputArgs struct {
	decisionMakerOptions    decisionmaker.DecisionMaker0Options
	PriceVariationDetection float32
	InitialAmount           float64
}

type Algo1BenchmarkInputArgs struct {
	decisionMakerOptions    decisionmaker.Options
	PriceVariationDetection float32
	InitialAmount           float64
	TotalPointsHolding      int
}

type BenchmarkResult struct {
	input         interface{}
	buys          int
	sells         int
	sellsPending  int
	initialAmount float32
	finalAmount   float32
	profit        float32
	assets        *[]domain.Asset
}

func RoundDown(input float64, places int) (newVal float32) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * input
	round = math.Floor(digit)
	newVal = float32(round / pow)
	return
}

func WriteResult(br BenchmarkResult) {
	fmt.Println("======")
	fmt.Printf("%+v\n", br.input)
	fmt.Printf("Buys %v\n", br.buys)
	fmt.Printf("Sells %v\n", br.sells)
	fmt.Printf("Sells Pending %v\n", br.sellsPending)
	fmt.Printf("Initial amount %v\n", br.initialAmount)
	fmt.Printf("Final amount %v\n", br.finalAmount)
	fmt.Printf("Returns %0.2f%% \n", br.profit)
	fmt.Println("=======")
}

func benchmark(decisionMaker domain.DecisionMaker, assetsRepository *AssetsRepositoryMock, inputArgs interface{}, priceVariationDetection float32, initialAmount float64, done chan BenchmarkResult) {
	bitcoinHistoryCollector := collectors.NewBitcoinHistoryCollector(priceVariationDetection)

	notificationsService := &NotificationsMock{}
	logService := &LogMock{}
	accountService := &AccountServiceMock{float32(initialAmount), 0, 0}
	broker := broker.NewBrokerMock()
	trader := trader.NewTrader(assetsRepository, accountService, broker)

	application := app.NewApp(notificationsService, decisionMaker, logService, assetsRepository, trader, accountService)

	historyFile, err := collectors.GetCsv(fmt.Sprintf("./data-history/btc/%v", btcdatahistory.LastYearMinute))

	if err != nil {
		log.Fatal(err)
	}

	bitcoinHistoryCollector.Start(historyFile, func(ask, bid float32, buyTime time.Time) {
		application.OnTickerChange(ask, bid, buyTime)
	})

	var sells int
	for _, asset := range assetsRepository.Assets {
		if asset.Sold {
			// fmt.Printf("#%v B:%v S:%v\n", index, asset.BuyPrice, asset.SellPrice)
			sells++
		} else {
			// fmt.Printf("#%v B:%v\n", index, asset.BuyPrice, asset.SellPrice)
		}
	}

	benchmarkResult := BenchmarkResult{
		inputArgs,
		len(assetsRepository.Assets),
		sells,
		len(assetsRepository.Assets) - sells,
		float32(initialAmount),
		accountService.Amount,
		RoundDown(((float64(accountService.Amount)-initialAmount)*100)/initialAmount, 2),
		&assetsRepository.Assets,
	}

	done <- benchmarkResult
}

func BenchmarkAlgo1(done chan BenchmarkResult) int {
	initialAmount := []float64{2000}
	maximumBuyAmount := []float32{0.1, 0.2, 0.3, 0.4, 0.5}
	pretendedProfitPerSold := []float32{0.01}
	priceDropToBuy := []float32{0.01}
	priceVariationDetection := []float32{0.01}
	totalPointsHolding := []int{38000}

	cases := []Algo1BenchmarkInputArgs{}

	for _, ia := range initialAmount {
		for _, mba := range maximumBuyAmount {
			for _, pfps := range pretendedProfitPerSold {
				for _, pdtb := range priceDropToBuy {
					for _, pvd := range priceVariationDetection {
						for _, tph := range totalPointsHolding {
							cases = append(cases, Algo1BenchmarkInputArgs{decisionmaker.Options{mba, pfps, pdtb}, pvd, ia, tph})
						}
					}
				}
			}
		}
	}

	for _, options := range cases {
		statisticsOptions := statistics.Options{options.TotalPointsHolding}
		macd := statistics.NewMACDContainer(statistics.MACDParams{12, 26, 9}, []float64{})
		statisticsService := statistics.NewStatistics(statisticsOptions, macd)
		assetsRepository := &AssetsRepositoryMock{}
		decisionMaker := decisionmaker.NewDecisionMaker(assetsRepository, options.decisionMakerOptions, statisticsService)
		go benchmark(decisionMaker, assetsRepository, options, options.PriceVariationDetection, options.InitialAmount, done)
	}

	return len(cases)
}

func BenchmarkAlgo0(done chan BenchmarkResult) int {
	initialAmount := []float64{500}
	maximumBuyAmount := []float32{0.01}
	pretendedProfitPerSold := []float32{0.01}
	priceDropToBuy := []float32{0.01}
	priceVariationDetection := []float32{0.01}

	cases := []Algo0BenchmarkInputArgs{}

	for _, ia := range initialAmount {
		for _, mba := range maximumBuyAmount {
			for _, pfps := range pretendedProfitPerSold {
				for _, pdtb := range priceDropToBuy {
					for _, pvd := range priceVariationDetection {
						cases = append(cases, Algo0BenchmarkInputArgs{decisionmaker.DecisionMaker0Options{mba, pfps, pdtb}, pvd, ia})
					}
				}
			}
		}
	}

	for _, options := range cases {
		assetsRepository := &AssetsRepositoryMock{}
		decisionMaker := decisionmaker.NewDecisionMaker0(assetsRepository, options.decisionMakerOptions)
		go benchmark(decisionMaker, assetsRepository, options, options.PriceVariationDetection, options.InitialAmount, done)
	}

	return len(cases)
}

func main() {
	start := time.Now()

	reportsCh := make(chan BenchmarkResult)

	iterations := BenchmarkAlgo1(reportsCh)

	startDate := time.Now().Format("2006-01-02T15:04:05Z07:00")
	f, err := os.Create(fmt.Sprintf("./reports/benchmark-reports/benchmark-%v.csv", startDate))
	if err != nil {
		log.Fatal(err)
	}

	f.Write([]byte("Case,Buys,Sells,Sells Pending,Initial Amount,Final Amount,Profit\n"))
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < iterations; i++ {
		br := <-reportsCh
		fmt.Printf("\r%d/%d", i+1, iterations)
		f.WriteString(fmt.Sprintf("%+v,%d,%d,%d,%.2f,%.2f,%.2f%%\n", br.input, br.buys, br.sells, br.sellsPending, br.initialAmount, br.finalAmount, br.profit))
		fOrders, err := os.Create(fmt.Sprintf("./reports/orders-reports/benchmark-%v-orders-%v.csv", startDate, i))

		fOrders.WriteString(fmt.Sprintf("Buy Date,Sell Date,Amount,Buy Price,Buy Value,Sell Price,Sell Value,Return\n"))
		for _, asset := range *br.assets {
			buyValue := asset.Amount * asset.BuyPrice
			sellValue := asset.Amount * asset.SellPrice
			fOrders.WriteString(fmt.Sprintf("%v,%v,%.2f,%.2f,%.2f,%.2f,%.2f,%.2f\n", asset.BuyTime, asset.SellTime, asset.Amount, asset.BuyPrice, buyValue, asset.SellPrice, sellValue, sellValue-buyValue))
		}

		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Printf("\n%v", time.Since(start))

}
