package benchmark

import (
	"fmt"
	"log"
	"time"

	"github.com/fabiodmferreira/crypto-trading/accounts"
	"github.com/fabiodmferreira/crypto-trading/app"
	"github.com/fabiodmferreira/crypto-trading/assets"
	"github.com/fabiodmferreira/crypto-trading/broker"
	"github.com/fabiodmferreira/crypto-trading/collectors"
	"github.com/fabiodmferreira/crypto-trading/domain"
	"github.com/fabiodmferreira/crypto-trading/eventlogs"
	"github.com/fabiodmferreira/crypto-trading/notifications"
	"github.com/fabiodmferreira/crypto-trading/trader"
)

// Output is the output of the benchmark
type Output struct {
	Input         interface{}
	Buys          int
	Sells         int
	SellsPending  int
	InitialAmount float32
	FinalAmount   float32
	Profit        float32
	Assets        *[]domain.Asset
}

// String displays Output formatted
func (o *Output) String() {
	fmt.Println("======")
	fmt.Printf("%+v\n", o.Input)
	fmt.Printf("Buys %v\n", o.Buys)
	fmt.Printf("Sells %v\n", o.Sells)
	fmt.Printf("Sells Pending %v\n", o.SellsPending)
	fmt.Printf("Initial amount %v\n", o.InitialAmount)
	fmt.Printf("Final amount %v\n", o.FinalAmount)
	fmt.Printf("Returns %0.2f%% \n", o.Profit)
	fmt.Println("=======")
}

// Benchmark runs algorithm and returns performance results
func Benchmark(decisionMaker domain.DecisionMaker, assetsRepository *assets.AssetsRepositoryInMemory, inputArgs interface{}, priceVariationDetection float32, initialAmount float64, filePath string, done chan Output) {
	bitcoinHistoryCollector := collectors.NewBitcoinHistoryCollector(priceVariationDetection)

	notificationsService := &notifications.NotificationsMock{}
	logService := &eventlogs.EventLogsServiceMock{}
	accountService := accounts.NewAccountServiceInMemory(float32(initialAmount))
	broker := broker.NewBrokerMock()
	trader := trader.NewTrader(assetsRepository, accountService, broker)

	application := app.NewApp(notificationsService, decisionMaker, logService, assetsRepository, trader, accountService)

	historyFile, err := collectors.GetCsv(fmt.Sprintf("./data-history/%v", filePath))

	if err != nil {
		log.Fatal(err)
	}

	bitcoinHistoryCollector.Start(historyFile, func(ask, bid float32, buyTime time.Time) {
		application.OnTickerChange(ask, bid, buyTime)
	})

	var sells int
	for _, asset := range assetsRepository.Assets {
		if asset.Sold {
			sells++
		} else {
		}
	}

	Output := Output{
		inputArgs,
		len(assetsRepository.Assets),
		sells,
		len(assetsRepository.Assets) - sells,
		float32(initialAmount),
		accountService.Amount,
		float32(((float64(accountService.Amount) - initialAmount) * 100) / initialAmount),
		&assetsRepository.Assets,
	}

	done <- Output
}
