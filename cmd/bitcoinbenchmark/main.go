package main

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/fabiodmferreira/crypto-trading/app"
	"github.com/fabiodmferreira/crypto-trading/assets"
	"github.com/fabiodmferreira/crypto-trading/broker"
	"github.com/fabiodmferreira/crypto-trading/collectors"
	"github.com/fabiodmferreira/crypto-trading/decisionmaker"
	"github.com/fabiodmferreira/crypto-trading/trader"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotificationsMock struct {
	checks int
}

func (n *NotificationsMock) CheckEventLogs() error {
	n.checks++
	return nil
}

type LogMock struct {
	logs [][]string
}

func (l *LogMock) Create(logType, message string) error {
	l.logs = append(l.logs, []string{logType, message})
	return nil
}

type AccountMock struct {
	Amount float32
}

func (a *AccountMock) Deposit(amount float32) error {
	a.Amount += amount
	// if a.Amount < 5000 {
	// 	fmt.Printf("%v ", a.Amount)
	// }
	return nil
}

func (a *AccountMock) Withdraw(amount float32) error {
	if amount > a.Amount {
		return errors.New("Insufficient Funds")
	}

	a.Amount -= amount
	// if a.Amount < 5000 {
	// 	fmt.Printf("%v ", a.Amount)
	// }

	return nil
}

type AssetsRepositoryMock struct {
	Assets []assets.Asset
}

func (ar *AssetsRepositoryMock) FindAll() (*[]assets.Asset, error) {
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

func (ar *AssetsRepositoryMock) Create(asset *assets.Asset) error {
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

func RoundDown(input float64, places int) (newVal float32) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * input
	round = math.Floor(digit)
	newVal = float32(round / pow)
	return
}

func benchmark(decisionMakerOptions decisionmaker.DecisionMakerOptions, PriceVariationDetection float32) {
	bitcoinHistoryCollector := collectors.NewBitcoinHistoryCollector()

	notificationsService := &NotificationsMock{}
	logService := &LogMock{}
	account := &AccountMock{5000}
	broker := broker.NewBrokerMock()
	assetsRepository := &AssetsRepositoryMock{}
	trader := trader.NewTrader(assetsRepository, logService, broker)

	decisionMaker := decisionmaker.NewDecisionMaker(trader, account, assetsRepository, decisionMakerOptions)

	application := app.NewApp(notificationsService, decisionMaker, logService, PriceVariationDetection)

	bitcoinHistoryCollector.Start(func(ask, bid float32, buyTime time.Time) {
		// fmt.Printf("* %v\n", ask)
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

	fmt.Println("======")
	fmt.Printf("%+v %+v\n", decisionMakerOptions, PriceVariationDetection)
	fmt.Printf("Buys %v\n", len(assetsRepository.Assets))
	fmt.Printf("Sells %v\n", sells)
	fmt.Printf("Sells Pending %v\n", len(assetsRepository.Assets)-sells)
	fmt.Printf("Initial amount 5000\n")
	fmt.Printf("Final amount %v\n", account.Amount)
	fmt.Printf("Returns %0.2f%% \n", RoundDown(((float64(account.Amount)-5000)*100)/5000, 2))
	fmt.Println("=======")
	// for _, asset := range assetsRepository.Assets {
	// 	fmt.Printf("%v %v\n", asset.BuyTime, asset.BuyPrice)
	// }
}

func main() {
	start := time.Now()

	cases := []struct {
		decisionMakerOptions    decisionmaker.DecisionMakerOptions
		PriceVariationDetection float32
	}{
		{
			decisionmaker.DecisionMakerOptions{0.01, 0.01, 0.1},
			float32(0.01),
		},
		{
			decisionmaker.DecisionMakerOptions{0.01, 0.001, 0.1},
			float32(0.001),
		},
		{
			decisionmaker.DecisionMakerOptions{0.01, 0.001, 0.1},
			float32(0.0001),
		},
	}

	for _, options := range cases {
		benchmark(options.decisionMakerOptions, options.PriceVariationDetection)
	}

	fmt.Println(time.Since(start))

}
