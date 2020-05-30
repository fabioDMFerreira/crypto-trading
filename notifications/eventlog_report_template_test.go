package notifications

import (
	"strings"
	"testing"
	"time"

	"github.com/fabiodmferreira/crypto-trading/assets"
	"github.com/fabiodmferreira/crypto-trading/eventlogs"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGenerateEventlogReportEmail(t *testing.T) {
	t.Run("should return template", func(t *testing.T) {
		id1, _ := primitive.ObjectIDFromHex("5ed1836b3032202176ddf411")
		id2, _ := primitive.ObjectIDFromHex("5ed1836b3032202176ddf412")

		date := time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
		got, err := GenerateEventlogReportEmail(
			5000,
			10,
			-20,
			date,
			date,
			&[]eventlogs.EventLog{
				{
					ID:          id1,
					EventName:   "bought",
					Message:     "asset bought",
					Notified:    false,
					DateCreated: date,
				},
			},
			&[]assets.Asset{
				{
					ID:       id2,
					Amount:   0.01,
					BuyPrice: 9000,
					BuyTime:  date,
				},
			},
		)

		if err != nil {
			t.Error(err)
		}

		if !strings.Contains(got.String(), "5ed1836b3032202176ddf411") {
			t.Errorf("%v\nSomething is wrong with the template", got)
		}
	})
}
