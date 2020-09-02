package main

import (
	"fmt"
	"log"
	"os"
	"time"

	btcdatahistory "github.com/fabiodmferreira/crypto-trading/data-history/btc"
	"github.com/fabiodmferreira/crypto-trading/domain"
	"github.com/fabiodmferreira/crypto-trading/statistics"

	"github.com/fabiodmferreira/crypto-trading/collectors"
)

func main() {
	historyFile, err := collectors.GetCsv(fmt.Sprintf("./data-history/%v", btcdatahistory.TwentyTwentyH1))

	if err != nil {
		log.Fatal(err)
	}

	startDate := time.Now().Format("2006-01-02T15:04:05Z07:00")
	f, err := os.Create(fmt.Sprintf("./reports/change-reports/change-%v.csv", startDate))
	if err != nil {
		log.Fatal(err)
	}

	bitcoinHistoryCollector := collectors.NewFileTickerCollector(domain.CollectorOptions{PriceVariationDetection: 0.01, DataSource: historyFile})

	statisticsOptions := domain.StatisticsOptions{NumberOfPointsHold: 38000}
	macd := statistics.NewMACDContainer(statistics.MACDParams{Fast: 12, Slow: 26, Lag: 9}, []float64{})
	stats := statistics.NewStatistics(statisticsOptions, macd)

	f.Write([]byte("Date,Price,Change,ChangeofChange,VelocityDirection,AccelerationDirection,Histogram\n"))

	var lastPrice float32
	var lastChange float32
	var velocityCurrentDirection bool
	var timesSameVelocityDirection int
	var accelerationCurrentDirection bool
	var timesSameAccelerationDirection int

	bitcoinHistoryCollector.Regist(func(ohlc *domain.OHLC) {
		stats.AddPoint(float64(ohlc.Close))
		if lastPrice > 0 {
			change := ohlc.Close - lastPrice
			changeOfChange := change - lastChange
			if accelerationCurrentDirection && changeOfChange > 0 || !accelerationCurrentDirection && changeOfChange < 0 {
				timesSameAccelerationDirection++
			} else {
				timesSameAccelerationDirection = 0
				accelerationCurrentDirection = !accelerationCurrentDirection
			}

			if velocityCurrentDirection && change > 0 || !velocityCurrentDirection && change < 0 {
				timesSameVelocityDirection++
			} else {
				timesSameVelocityDirection = 0
				velocityCurrentDirection = !velocityCurrentDirection
			}

			f.WriteString(fmt.Sprintf("%v,%.2f,%.2f,%.2f,%d,%d,%.2f\n", ohlc.Time.Format("2006-01-02T15:04:05"), ohlc.Close, change, changeOfChange, timesSameVelocityDirection, timesSameAccelerationDirection, stats.MACD.GetLastHistogramPoint()))

			lastPrice = ohlc.Close
			lastChange = change
		} else {
			f.WriteString(fmt.Sprintf("%v,%.2f,%d,%d\n", ohlc.Time, ohlc.Close, 0, 0))
			lastPrice = ohlc.Close
		}
	})

	bitcoinHistoryCollector.Start()
}
