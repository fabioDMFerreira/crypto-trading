package webserver

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// GroupByDatesID returns grioup by date clause to be used on aggregate query
func GroupByDatesID(startDate time.Time, endDate time.Time) bson.D {
	days := endDate.Sub(startDate).Hours() / 24

	if days > 30 {
		return bson.D{
			{"year", bson.D{{"$year", "$date"}}},
			{"month", bson.D{{"$month", "$date"}}},
			{"day", bson.D{{"$dayOfMonth", "$date"}}},
		}
	} else if days > 5 {
		return bson.D{
			{"year", bson.D{{"$year", "$date"}}},
			{"month", bson.D{{"$month", "$date"}}},
			{"day", bson.D{{"$dayOfMonth", "$date"}}},
			{"hour", bson.D{{"$hour", "$date"}}},
		}
	} else {
		return bson.D{
			{"year", bson.D{{"$year", "$date"}}},
			{"month", bson.D{{"$month", "$date"}}},
			{"day", bson.D{{"$dayOfMonth", "$date"}}},
			{"hour", bson.D{{"$hour", "$date"}}},
			{"minute", bson.D{{"$minute", "$date"}}},
		}
	}

}
