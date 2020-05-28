package main

import (
	"fmt"

	"github.com/fabiodmferreira/crypto-trading/collector"
)

func main() {
	fmt.Println("Starting Coin Historical Collector v0")
	bhc := collector.NewBitcoinHistoricalCollector()
	bhc.Start()
}
