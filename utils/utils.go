package utils

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// DollarEuroRate is a value used to convert dollars to euros
const DollarEuroRate = 0.8749291093063332

// GetGroupByDatesIDClause returns group by date clause to be used on aggregate query
func GetGroupByDatesIDClause(startDate time.Time, endDate time.Time) bson.M {
	days := endDate.Sub(startDate).Hours() / 24

	if days > 30 {
		return bson.M{
			"year":  bson.M{"$year": "$date"},
			"month": bson.M{"$month": "$date"},
			"day":   bson.M{"$dayOfMonth": "$date"},
		}
	} else if days > 5 {
		return bson.M{
			"year":  bson.M{"$year": "$date"},
			"month": bson.M{"$month": "$date"},
			"day":   bson.M{"$dayOfMonth": "$date"},
			"hour":  bson.M{"$hour": "$date"},
		}
	} else {
		return bson.M{
			"year":   bson.M{"$year": "$date"},
			"month":  bson.M{"$month": "$date"},
			"day":    bson.M{"$dayOfMonth": "$date"},
			"hour":   bson.M{"$hour": "$date"},
			"minute": bson.M{"$minute": "$date"},
		}
	}
}
