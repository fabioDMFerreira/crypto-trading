package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/fabiodmferreira/crypto-trading/domain"
)

func TestWriteCoindeskResponse(t *testing.T) {
	response := domain.CoindeskHTTPResponse{
		StatusCode: 200,
		Message:    "OK",
		Data: domain.CoindeskResponse{
			Iso:      "BTC",
			Name:     "Bitcoin",
			Slug:     "bitcoin",
			Interval: "1-minute",
			Entries: [][]float64{
				[]float64{1577836859000, 7163.0865663227},
				[]float64{1577836919000, 7164.6556003506},
				[]float64{1577836979000, 7159.9213071402},
			},
		},
	}

	t.Run("should write coindesk response to a writer", func(t *testing.T) {
		buffer := &bytes.Buffer{}

		WriteCoindeskResponse(response.Data, buffer)

		got := buffer.String()
		want := `1577836859000,7163.086566
		1577836919000,7164.655600
		1577836979000,7159.921307`

		if strings.Compare(got, want) == 0 {
			t.Errorf("got %v want %v", got, want)
		}
	})
}
