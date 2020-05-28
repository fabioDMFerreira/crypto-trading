package collectors

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type BitcoinHistoricalCollector struct {
}

func write(date, price string, wg *sync.WaitGroup) {
	defer wg.Done()

	t, _ := time.Parse("Jan 02, 2006", date)
	fmt.Printf("%v %v \n", t.Format("2006-1-02"), price)
}

func (bhc *BitcoinHistoricalCollector) Start() {
	// Open the file
	filename, err := filepath.Abs("./collector/btc-historical-data-2012-2020.csv")

	if err != nil {
		log.Fatalln(err)
	}

	csvfile, err := os.Open(filename)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	// Parse the file
	r := csv.NewReader(csvfile)
	r.LazyQuotes = true
	//r := csv.NewReader(bufio.NewReader(csvfile))

	// read header
	header, err := r.Read()
	if err == io.EOF {
		return
	}

	fmt.Printf("%v\n", header)

	var wg sync.WaitGroup

	// Iterate through the records
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		wg.Add(1)
		go write(record[0], record[1], &wg)
	}

	wg.Wait()
}

func NewBitcoinHistoricalCollector() *BitcoinHistoricalCollector {
	return &BitcoinHistoricalCollector{}
}
