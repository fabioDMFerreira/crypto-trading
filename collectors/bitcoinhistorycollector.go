package collectors

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type BitcoinHistoryCollector struct {
	priceVariationDetection float32
	lastTickerPrice         float32
}

func NewBitcoinHistoryCollector(priceVariationDetection float32) *BitcoinHistoryCollector {
	return &BitcoinHistoryCollector{priceVariationDetection, 0}
}

func (bhc *BitcoinHistoryCollector) Start(historyFile *csv.Reader, onChange OnTickerChange) {
	// read header
	_, err := historyFile.Read()
	if err == io.EOF {
		log.Fatalf("Error on reading bitcoin history file", err)
	}

	for {
		// Read each record from csv
		record, err := historyFile.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		priceStr := strings.ReplaceAll(record[1], ",", "")
		price, err := strconv.ParseFloat(priceStr, 32)

		if err != nil {
			log.Fatalf("Error on converting bitcoin price from file", err)
		}

		unixTime, err := strconv.ParseInt(record[0], 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		changeVariance := float32(bhc.lastTickerPrice * bhc.priceVariationDetection)

		if bhc.lastTickerPrice == 0 ||
			float32(price) > bhc.lastTickerPrice+changeVariance ||
			float32(price) < bhc.lastTickerPrice-changeVariance {
			onChange(float32(price), float32(price), time.Unix(unixTime/1000, 0))
		}
	}
}

func GetCsv(file string) (*csv.Reader, error) {
	// Open the file
	filename, err := filepath.Abs(file)

	if err != nil {
		log.Fatalln(err)
	}

	csvfile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	// Parse the file
	r := csv.NewReader(csvfile)
	r.LazyQuotes = true
	//r := csv.NewReader(bufio.NewReader(csvfile))

	return r, nil
}
