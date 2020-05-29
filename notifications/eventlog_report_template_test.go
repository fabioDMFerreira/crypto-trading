package notifications

import (
	"testing"
	"time"

	"github.com/fabiodmferreira/crypto-trading/assets"
	"github.com/fabiodmferreira/crypto-trading/eventlogs"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	ReportHTML = `<!DOCTYPE html>
	<html>

	<head>
		<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css"
			integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossorigin="anonymous">
	</head>

	<body>
		<p>
			<label>Current Amount:</label> 5000.00€<br />
			<label>Total Orders:</label> 10<br />
		</p>

		<h2>Events Log</h2>
		<table class="table table-bordered table-striped">
			<tr>
				<th>Id</th>
				<th>Date</th>
				<th>Event</th>
			</tr>

				<tr>
					<td>5ed1836b3032202176ddf411</td>
					<td>17-Nov-2009 20:34</td>
					<td>asset bought</td>
				</tr>

		</table>

		<h2>Assets waiting to be sold</h2>
		<table class="table table-bordered table-striped">
			<tr>
				<th>Date</th>
				<th>Amount</th>
				<th>Buy Price</th>
			</tr>

				<tr>
					<td>5ed1836b3032202176ddf412</td>
					<td>17-Nov-2009 20:34</td>
					<td>0.01BTC</td>
					<td>9000.00€</td>
				</tr>

		</table>

	</body>

	</html>`
)

func TestGenerateEventlogReportEmail(t *testing.T) {
	t.Run("should return template", func(t *testing.T) {
		id1, _ := primitive.ObjectIDFromHex("5ed1836b3032202176ddf411")
		id2, _ := primitive.ObjectIDFromHex("5ed1836b3032202176ddf412")

		got, err := GenerateEventlogReportEmail(
			5000,
			10,
			&[]eventlogs.EventLog{
				{
					ID:          id1,
					EventName:   "bought",
					Message:     "asset bought",
					Notified:    false,
					DateCreated: time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
				},
			},
			&[]assets.Asset{
				{
					ID:       id2,
					Amount:   0.01,
					BuyPrice: 9000,
					BuyTime:  time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
				},
			},
		)

		want := ReportHTML

		if err != nil {
			t.Error(err)
		}

		if got.String() != want {
			t.Errorf("got %v want %v", got, want)
		}
	})
}
