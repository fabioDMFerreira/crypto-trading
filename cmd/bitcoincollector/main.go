package main

import (
	"fmt"

	"github.com/fabiodmferreira/crypto-trading/collectors"
)

func main() {
	fmt.Println("Starting Coin Historical Collector v0")
	bhc := collectors.NewBitcoinHistoricalCollector()
	bhc.Start()
}
