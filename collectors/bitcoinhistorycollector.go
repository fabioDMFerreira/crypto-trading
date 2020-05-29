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

type BitcoinHistoryCollector struct{}

func NewBitcoinHistoryCollector() *BitcoinHistoryCollector {
	return &BitcoinHistoryCollector{}
}

func (bhc *BitcoinHistoryCollector) Start(onChange OnTickerChange) {
	bitcoinHistoryFile, err := GetBitcoinHistoryFile()

	if err != nil {
		log.Fatalf("Error on getting bitcoin history file", err)
	}

	// read header
	_, err = bitcoinHistoryFile.Read()
	if err == io.EOF {
		log.Fatalf("Error on reading bitcoin history file", err)
	}

	for {
		// Read each record from csv
		record, err := bitcoinHistoryFile.Read()
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

		onChange(float32(price), float32(price), time.Unix(unixTime/1000, 0))
	}
}

func GetBitcoinHistoryFile() (*csv.Reader, error) {
	// Open the file
	filename, err := filepath.Abs("./collectors/2020.csv")

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
