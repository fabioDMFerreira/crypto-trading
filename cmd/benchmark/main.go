package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fabiodmferreira/crypto-trading/assetsprices"
	"github.com/fabiodmferreira/crypto-trading/benchmark"
	btcdatahistory "github.com/fabiodmferreira/crypto-trading/data-history/btc"
	"github.com/fabiodmferreira/crypto-trading/domain"
)

// ExecuteBenchmark create params and execute benchmarks
func ExecuteBenchmark(done chan benchmark.BenchmarkResult) int {
	filesPaths := []string{
		// fmt.Sprintf("ada/%v", adadatahistory.LastYearMinute),
		fmt.Sprintf("%v", btcdatahistory.LastYearMinute),
		// fmt.Sprintf("btc-cash/%v", btccashdatahistory.LastYearMinute),
		// fmt.Sprintf("btc-sv/%v", btcsvdatahistory.LastYearMinute),
		// fmt.Sprintf("eos/%v", eosdatahistory.LastYearMinute),
		// fmt.Sprintf("etc/%v", etcdatahistory.LastYearMinute),
		// fmt.Sprintf("eth/%v", ethdatahistory.LastYearMinute),
		// fmt.Sprintf("ltc/%v", ltcdatahistory.LastYearMinute),
		// fmt.Sprintf("monero/%v", monerodatahistory.LastYearMinute),
		// fmt.Sprintf("stellar/%v", stellardatahistory.LastYearMinute),
		// fmt.Sprintf("xrp/%v", xrpdatahistory.LastYearMinute),
	}
	initialAmount := []float64{2000}
	maximumBuyAmount := []float32{0.1}
	pretendedProfitPerSold := []float32{0.01, 0.03}
	priceDropToBuy := []float32{0.01}
	priceVariationDetection := []float32{0.01}
	totalPointsHolding := []int{500}

	cases := []benchmark.Input{}

	for _, ia := range initialAmount {
		for _, mba := range maximumBuyAmount {
			for _, pfps := range pretendedProfitPerSold {
				for _, pdtb := range priceDropToBuy {
					for _, pvd := range priceVariationDetection {
						for _, tph := range totalPointsHolding {
							for _, fp := range filesPaths {
								input := domain.BenchmarkInput{
									DecisionMakerOptions: domain.DecisionMakerOptions{MaximumBuyAmount: mba, MinimumProfitPerSold: pfps, MinimumPriceDropToBuy: pdtb},
									StatisticsOptions:    domain.StatisticsOptions{NumberOfPointsHold: tph},
									CollectorOptions:     domain.CollectorOptions{PriceVariationDetection: pvd},
									AccountInitialAmount: ia,
									DataSourceFilePath:   fp,
								}
								cases = append(cases, input)
							}
						}
					}
				}
			}
		}
	}

	benchmark := benchmark.NewService(benchmark.NewRepositoryInMemory(), new(assetsprices.RepositoryInMemory))

	benchmark.BulkRun(cases, done)

	return len(cases)
}

func main() {
	start := time.Now()

	reportsCh := make(chan benchmark.BenchmarkResult)

	iterations := ExecuteBenchmark(reportsCh)

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

		if br.Err == nil {
			result := br.Output
			input := br.Input
			profit := float32(((float64(result.FinalAmount) - input.AccountInitialAmount) * 100) / input.AccountInitialAmount)
			f.WriteString(fmt.Sprintf("%+v,%d,%d,%d,%.2f,%.2f,%.2f%%\n", input, result.Buys, result.Sells, result.SellsPending, input.AccountInitialAmount, result.FinalAmount, profit))
			fOrders, err := os.Create(fmt.Sprintf("./reports/orders-reports/benchmark-%v-orders-%v.csv", startDate, i))

			fOrders.WriteString(fmt.Sprintf("Buy Date,Sell Date,Amount,Buy Price,Buy Value,Sell Price,Sell Value,Return\n"))
			for _, asset := range *result.Assets {
				buyValue := asset.Amount * asset.BuyPrice
				sellValue := asset.Amount * asset.SellPrice
				fOrders.WriteString(fmt.Sprintf("%v,%v,%.2f,%.2f,%.2f,%.2f,%.2f,%.2f\n", asset.BuyTime, asset.SellTime, asset.Amount, asset.BuyPrice, buyValue, asset.SellPrice, sellValue, sellValue-buyValue))
			}

			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal(br.Err)
		}

	}

	fmt.Printf("\n%v", time.Since(start))

}
