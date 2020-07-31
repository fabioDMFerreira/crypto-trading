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

	"github.com/fabiodmferreira/crypto-trading/domain"
)

// FileTickerCollector collects data from a csv file
type FileTickerCollector struct {
	options              domain.CollectorOptions
	lastTickerPrice      float32
	observables          []domain.OnTickerChange
	lastPricePublishDate time.Time
}

// NewFileTickerCollector returns an instance of FileTickerCollector
func NewFileTickerCollector(options domain.CollectorOptions) *FileTickerCollector {
	return &FileTickerCollector{options, 0, []domain.OnTickerChange{}, time.Time{}}
}

// Start starts collecting data from data source
func (ftc *FileTickerCollector) Start() {
	// read header
	_, err := ftc.options.DataSource.Read()
	if err == io.EOF {
		log.Fatalf("Error on reading header file: %v", err)
	}

	for {
		// Read each record from csv
		record, err := ftc.options.DataSource.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		priceStr := strings.ReplaceAll(record[1], ",", "")
		price, err := strconv.ParseFloat(priceStr, 32)

		if err != nil {
			log.Fatalf("Error on converting price from file:\n%v", err)
		}

		unixTime, err := strconv.ParseInt(record[0], 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		date := time.Unix(unixTime/1000, 0)

		changeVariance := float32(ftc.lastTickerPrice * ftc.options.PriceVariationDetection)

		timeSinceLastPricePublished := date.Sub(ftc.lastPricePublishDate).Minutes()

		if timeSinceLastPricePublished > float64(ftc.options.NewPriceTimeRate) && (ftc.lastTickerPrice == 0 ||
			float32(price) > ftc.lastTickerPrice+changeVariance ||
			float32(price) < ftc.lastTickerPrice-changeVariance) {
			for _, observable := range ftc.observables {
				observable(float32(price), float32(price), date)

				ftc.lastPricePublishDate = date
			}
		}

	}
}

// Regist add function to be executed when a new price is received
func (ftc *FileTickerCollector) Regist(observable domain.OnTickerChange) {
	ftc.observables = append(ftc.observables, observable)
}

// GetCsv returns the pointer for a csv file from the filepath
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
