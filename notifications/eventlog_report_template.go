package notifications

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/fabiodmferreira/crypto-trading/domain"
)

type EventsLogEmailFormat struct {
	Date    string
	Message string
	ID      string
}

type PendingAssetsEmailFormat struct {
	Date     string
	Amount   string
	BuyPrice string
	ID       string
}

func GenerateEventlogReportEmail(
	amount float32,
	balance float32,
	startDate time.Time,
	endDate time.Time,
	eventsLog *[]domain.EventLog,
	pendingAssets *[]domain.Asset,
) (*bytes.Buffer, error) {
	appEnv := os.Getenv("APP_ENV")

	eventLogFormatted := []EventsLogEmailFormat{}

	for _, event := range *eventsLog {
		eventLogFormatted = append(eventLogFormatted, EventsLogEmailFormat{event.CreatedAt.Format("02-Jan-2006 15:04"), event.Message, event.ID.Hex()})
	}

	pendingAssetsFormatted := []PendingAssetsEmailFormat{}

	for _, asset := range *pendingAssets {
		pendingAssetsFormatted = append(pendingAssetsFormatted, PendingAssetsEmailFormat{asset.BuyTime.Format("02-Jan-2006 15:04"), fmt.Sprintf("%.2fBTC", asset.Amount), fmt.Sprintf("%.2f€", asset.BuyPrice), asset.ID.Hex()})

	}

	var currentDir string

	if appEnv == "production" {
		currentDir = "notifications"
	} else {
		_, currentFilePath, _, _ := runtime.Caller(0)
		currentDir = path.Dir(currentFilePath)
	}
	t, err := template.ParseFiles(currentDir + "/eventlogreport.html")

	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)

	err = t.Execute(buf, struct {
		AccountAmount      string
		TotalAssetsPending string
		Balance            string
		StartDate          string
		EndDate            string
		EventsLog          []EventsLogEmailFormat
		AssetsPending      []PendingAssetsEmailFormat
	}{
		AccountAmount:      fmt.Sprintf("%.2f€", amount),
		TotalAssetsPending: fmt.Sprintf("%d", len(*pendingAssets)),
		Balance:            fmt.Sprintf("%.2f€", balance),
		StartDate:          startDate.Format("02-Jan-2006 15:04"),
		EndDate:            endDate.Format("02-Jan-2006 15:04"),
		EventsLog:          eventLogFormatted,
		AssetsPending:      pendingAssetsFormatted,
	})

	return buf, err
}
