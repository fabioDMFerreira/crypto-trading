package collectors

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/fabiodmferreira/crypto-trading/domain"
)

// FileTickerCollector collects data from a csv file
type FileTickerCollector struct {
	options              domain.CollectorOptions
	lastTickerPrice      float32
	observables          []domain.OnNewAssetPrice
	lastPricePublishDate time.Time
	indicators           *[]domain.Indicator
}

// NewFileTickerCollector returns an instance of FileTickerCollector
func NewFileTickerCollector(options domain.CollectorOptions, indicators *[]domain.Indicator) *FileTickerCollector {
	return &FileTickerCollector{
		options:    options,
		indicators: indicators,
	}
}

// GetTicker is a stub
func (ftc *FileTickerCollector) GetTicker(tickerSymbol string) (float32, error) {
	return 0, nil
}

func (ftc *FileTickerCollector) SetIndicators(indicators *[]domain.Indicator) {
	ftc.indicators = indicators
}

// Start starts collecting data from data source
func (ftc *FileTickerCollector) Start() {
	for {
		// Read each record from csv
		record, err := ftc.options.DataSource.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if len(record) < 6 {
			break
		}

		open, err := strconv.ParseFloat(record[1], 32)
		high, err1 := strconv.ParseFloat(record[2], 32)
		low, err2 := strconv.ParseFloat(record[3], 32)
		close, err3 := strconv.ParseFloat(record[4], 32)
		volume, err4 := strconv.ParseFloat(record[5], 32)

		if err != nil || err1 != nil || err2 != nil || err3 != nil || err4 != nil {
			log.Fatalf("Error on converting price from file:\n%v %v %v %v %v", err, err1, err2, err3, err4)
		}

		unixTime, err := strconv.ParseInt(record[0], 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		date := time.Unix(unixTime, 0)

		ohlc := &domain.OHLC{
			Time:    date,
			EndTime: date,
			Open:    float32(open),
			Close:   float32(close),
			High:    float32(high),
			Low:     float32(low),
			Volume:  float32(volume),
		}

		for _, indicator := range *ftc.indicators {
			indicator.AddValue(ohlc)
		}

		for _, observable := range ftc.observables {
			observable(ohlc)
		}

	}
}

// Stop is a stub
func (ftc *FileTickerCollector) Stop() {}

// Regist add function to be executed when a new price is received
func (ftc *FileTickerCollector) Regist(observable domain.OnNewAssetPrice) {
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
