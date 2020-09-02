package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	dateLayout := "2006-01-02 15:04:05"

	startDate, _ := time.Parse(dateLayout, "2019-01-01 00:00:00")
	endDate, _ := time.Parse(dateLayout, "2020-06-30 23:59:59")
	outputName := "2019-2020-h1.csv"

	rootPaths := []string{
		"ada",
		"btc",
		// "btc-cash",
		// "eos",
		// "etc",
		"eth",
		// "ltc",
		// "monero",
		// "stellar",
		// "xrp",
	}

	for _, path := range rootPaths {
		file, _ := os.Open(fmt.Sprintf("data-history/%v/dataset.csv", path))
		reader := csv.NewReader(file)
		reader.LazyQuotes = true

		writer, _ := os.Create(fmt.Sprintf("data-history/%v/%v.csv", path, outputName))

		fmt.Printf("writing %v/%v...\n", path, outputName)

		FilterFileByDatesInterval(reader, writer, int(startDate.Unix()), int(endDate.Unix()))

		fmt.Printf("%v done\n", path)
	}
}

func FilterFileByDatesInterval(in *csv.Reader, out io.Writer, startDate int, endDate int) {
	first := false

	counter := 0

	for {
		record, err := in.Read()

		counter++

		fmt.Printf("%v\r", counter)

		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		if len(record) > 1 {
			time, _ := strconv.Atoi(record[0])

			if startDate > time {
				continue
			}

			if time > endDate {
				break
			}

			if first {
				fmt.Fprint(out, "\n")
			} else {
				first = true
			}

			fmt.Fprint(out, strings.Join(record, ","))
		}
	}
}
