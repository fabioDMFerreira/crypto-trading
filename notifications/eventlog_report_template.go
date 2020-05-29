package notifications

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/fabiodmferreira/crypto-trading/assets"
	"github.com/fabiodmferreira/crypto-trading/eventlogs"
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

func GenerateEventlogReportEmail(amount float32, totalAssetsPending int, eventsLog *[]eventlogs.EventLog, assets *[]assets.Asset) (*bytes.Buffer, error) {
	eventLogFormatted := []EventsLogEmailFormat{}

	for _, event := range *eventsLog {
		eventLogFormatted = append(eventLogFormatted, EventsLogEmailFormat{event.DateCreated.Format("02-Jan-2006 15:04"), event.Message, event.ID.Hex()})
	}

	pendingAssetsFormatted := []PendingAssetsEmailFormat{}

	for _, asset := range *assets {
		pendingAssetsFormatted = append(pendingAssetsFormatted, PendingAssetsEmailFormat{asset.BuyTime.Format("02-Jan-2006 15:04"), fmt.Sprintf("%.2fBTC", asset.Amount), fmt.Sprintf("%.2f€", asset.BuyPrice), asset.ID.Hex()})

	}

	t, err := template.ParseFiles("notifications/eventlogreport.html")

	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)

	err = t.Execute(buf, struct {
		AccountAmount      string
		TotalAssetsPending string
		EventsLog          []EventsLogEmailFormat
		AssetsPending      []PendingAssetsEmailFormat
	}{
		AccountAmount:      fmt.Sprintf("%.2f€", amount),
		TotalAssetsPending: fmt.Sprintf("%d", totalAssetsPending),
		EventsLog:          eventLogFormatted,
		AssetsPending:      pendingAssetsFormatted,
	})

	return buf, err
}
