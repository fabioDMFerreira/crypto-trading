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
	Amount float32
}

func (a *AccountServiceMock) Deposit(amount float32) error {
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
	return &ar.Assets, nil
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

type BenchmarkInputArgs struct {
	decisionMakerOptions    decisionmaker.DecisionMakerOptions
	PriceVariationDetection float32
}

type BenchmarkResult struct {
	input         BenchmarkInputArgs
	buys          int
	sells         int
	sellsPending  int
	initialAmount float32
	finalAmount   float32
	profit        float32
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
	fmt.Printf("%+v %+v\n", br.input.decisionMakerOptions, br.input.PriceVariationDetection)
	fmt.Printf("Buys %v\n", br.buys)
	fmt.Printf("Sells %v\n", br.sells)
	fmt.Printf("Sells Pending %v\n", br.sellsPending)
	fmt.Printf("Initial amount %v\n", br.initialAmount)
	fmt.Printf("Final amount %v\n", br.finalAmount)
	fmt.Printf("Returns %0.2f%% \n", br.profit)
	fmt.Println("=======")
}

func benchmark(decisionMakerOptions decisionmaker.DecisionMakerOptions, priceVariationDetection float32, done chan BenchmarkResult) {
	bitcoinHistoryCollector := collectors.NewBitcoinHistoryCollector(priceVariationDetection)

	notificationsService := &NotificationsMock{}
	logService := &LogMock{}
	accountService := &AccountServiceMock{5000}
	broker := broker.NewBrokerMock()
	assetsRepository := &AssetsRepositoryMock{}
	trader := trader.NewTrader(assetsRepository, accountService, broker)

	decisionMaker := decisionmaker.NewDecisionMaker(assetsRepository, decisionMakerOptions)

	application := app.NewApp(notificationsService, decisionMaker, logService, assetsRepository, trader, accountService)

	historyFile, err := collectors.GetCsv(fmt.Sprintf("./data-history/btc/%v", btcdatahistory.May2020))

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
		BenchmarkInputArgs{decisionMakerOptions, priceVariationDetection},
		len(assetsRepository.Assets),
		sells,
		len(assetsRepository.Assets) - sells,
		5000,
		accountService.Amount,
		RoundDown(((float64(accountService.Amount)-5000)*100)/5000, 2),
	}

	done <- benchmarkResult
	// for _, asset := range assetsRepository.Assets {
	// 	fmt.Printf("%v %v\n", asset.BuyTime, asset.BuyPrice)
	// }
}

func main() {
	start := time.Now()

	maximumBuyAmount := []float32{0.01}
	pretendedProfitPerSold := []float32{0.005, 0.01, 0.02, 0.03, 0.04, 0.05, 0.06, 0.07, 0.08, 0.09, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6}
	priceDropToBuy := []float32{0.005, 0.01, 0.02, 0.03, 0.04, 0.05, 0.06, 0.07, 0.08, 0.09, 0.1}
	priceVariationDetection := []float32{0.005, 0.007}

	cases := []BenchmarkInputArgs{}

	for _, mba := range maximumBuyAmount {
		for _, pfps := range pretendedProfitPerSold {
			for _, pdtb := range priceDropToBuy {
				for _, pvd := range priceVariationDetection {
					cases = append(cases, BenchmarkInputArgs{decisionmaker.DecisionMakerOptions{mba, pfps, pdtb}, pvd})
				}
			}
		}
	}

	reportsCh := make(chan BenchmarkResult)

	for _, options := range cases {
		go benchmark(options.decisionMakerOptions, options.PriceVariationDetection, reportsCh)
	}

	f, err := os.Create(fmt.Sprintf("./benchmark-reports/benchmark-%v.csv", time.Now().Format("2006-01-02T15:04:05Z07:00")))
	if err != nil {
		log.Fatal(err)
	}

	f.Write([]byte("Case,Buys,Sells,Sells Pending,Initial Amount,Final Amount,Profit\n"))
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(cases); i++ {
		br := <-reportsCh
		fmt.Printf("\r%d/%d", i+1, len(cases))
		f.WriteString(fmt.Sprintf("%+v,%d,%d,%d,%.2f,%.2f,%.2f%%\n", br.input, br.buys, br.sells, br.sellsPending, br.initialAmount, br.finalAmount, br.profit))
	}

	fmt.Printf("\n%v", time.Since(start))

}
