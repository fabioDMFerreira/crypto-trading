package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fabiodmferreira/crypto-trading/assets"
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
	"github.com/fabiodmferreira/crypto-trading/decisionmaker"
	"github.com/fabiodmferreira/crypto-trading/statistics"
)

type Algo0BenchmarkInputArgs struct {
	decisionMakerOptions    decisionmaker.DecisionMaker0Options
	PriceVariationDetection float32
	InitialAmount           float64
	filePath                string
}

type Algo1BenchmarkInputArgs struct {
	decisionMakerOptions    decisionmaker.Options
	PriceVariationDetection float32
	InitialAmount           float64
	TotalPointsHolding      int
	filePath                string
}

func BenchmarkAlgo1(done chan benchmark.Output) int {
	filesPaths := []string{
		fmt.Sprintf("ada/%v", adadatahistory.LastYearMinute),
		fmt.Sprintf("btc/%v", btcdatahistory.LastYearMinute),
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
	initialAmount := []float64{2000}
	maximumBuyAmount := []float32{0.1}
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
							for _, fp := range filesPaths {
								cases = append(cases, Algo1BenchmarkInputArgs{decisionmaker.Options{mba, pfps, pdtb}, pvd, ia, tph, fp})
							}
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
		assetsRepository := &assets.AssetsRepositoryInMemory{}
		decisionMaker := decisionmaker.NewDecisionMaker(assetsRepository, options.decisionMakerOptions, statisticsService)
		go benchmark.Benchmark(decisionMaker, assetsRepository, options, options.PriceVariationDetection, options.InitialAmount, options.filePath, done)
	}

	return len(cases)
}

func BenchmarkAlgo0(done chan benchmark.Output) int {
	filesPaths := []string{
		fmt.Sprintf("./data-history/ada/%v", adadatahistory.LastYearMinute),
		fmt.Sprintf("./data-history/btc/%v", btcdatahistory.LastYearMinute),
		fmt.Sprintf("./data-history/btc-cash/%v", btccashdatahistory.LastYearMinute),
		fmt.Sprintf("./data-history/btc-sv/%v", btcsvdatahistory.LastYearMinute),
		fmt.Sprintf("./data-history/eos/%v", eosdatahistory.LastYearMinute),
		fmt.Sprintf("./data-history/etc/%v", etcdatahistory.LastYearMinute),
		fmt.Sprintf("./data-history/eth/%v", ethdatahistory.LastYearMinute),
		fmt.Sprintf("./data-history/ltc/%v", ltcdatahistory.LastYearMinute),
		fmt.Sprintf("./data-history/monero/%v", monerodatahistory.LastYearMinute),
		fmt.Sprintf("./data-history/stellar/%v", stellardatahistory.LastYearMinute),
		fmt.Sprintf("./data-history/xrp/%v", xrpdatahistory.LastYearMinute),
	}
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
						for _, fp := range filesPaths {
							cases = append(cases, Algo0BenchmarkInputArgs{decisionmaker.DecisionMaker0Options{mba, pfps, pdtb}, pvd, ia, fp})
						}
					}
				}
			}
		}
	}

	for _, options := range cases {
		assetsRepository := &assets.AssetsRepositoryInMemory{}
		decisionMaker := decisionmaker.NewDecisionMaker0(assetsRepository, options.decisionMakerOptions)
		go benchmark.Benchmark(decisionMaker, assetsRepository, options, options.PriceVariationDetection, options.InitialAmount, options.filePath, done)
	}

	return len(cases)
}

func main() {
	start := time.Now()

	reportsCh := make(chan benchmark.Output)

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
		f.WriteString(fmt.Sprintf("%+v,%d,%d,%d,%.2f,%.2f,%.2f%%\n", br.Input, br.Buys, br.Sells, br.SellsPending, br.InitialAmount, br.FinalAmount, br.Profit))
		fOrders, err := os.Create(fmt.Sprintf("./reports/orders-reports/benchmark-%v-orders-%v.csv", startDate, i))

		fOrders.WriteString(fmt.Sprintf("Buy Date,Sell Date,Amount,Buy Price,Buy Value,Sell Price,Sell Value,Return\n"))
		for _, asset := range *br.Assets {
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
