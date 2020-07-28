package webserver_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/fabiodmferreira/crypto-trading/webserver"
	"go.mongodb.org/mongo-driver/bson"
)

func TestGroupByDatesID(t *testing.T) {

	t.Run("should return query with year,month and day if dates difference is bigger than 30 days", func(t *testing.T) {
		got := webserver.GetGroupByDatesIDClause(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2020, 2, 20, 0, 0, 0, 0, time.UTC))
		want := bson.M{
			"year":  bson.M{"$year": "$date"},
			"month": bson.M{"$month": "$date"},
			"day":   bson.M{"$dayOfMonth": "$date"},
		}

		if reflect.DeepEqual(got, want) != true {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("should return query with year,month, day and hour if dates difference is lower than 30 days but bigger than 5", func(t *testing.T) {
		got := webserver.GetGroupByDatesIDClause(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2020, 1, 15, 0, 0, 0, 0, time.UTC))
		want := bson.M{
			"year":  bson.M{"$year": "$date"},
			"month": bson.M{"$month": "$date"},
			"day":   bson.M{"$dayOfMonth": "$date"},
			"hour":  bson.M{"$hour": "$date"},
		}

		if reflect.DeepEqual(got, want) != true {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("should return query with year,month, day, hour and minute if dates difference is lower than 5 days", func(t *testing.T) {
		got := webserver.GetGroupByDatesIDClause(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC))
		want := bson.M{
			"year":   bson.M{"$year": "$date"},
			"month":  bson.M{"$month": "$date"},
			"day":    bson.M{"$dayOfMonth": "$date"},
			"hour":   bson.M{"$hour": "$date"},
			"minute": bson.M{"$minute": "$date"},
		}

		if reflect.DeepEqual(got, want) != true {
			t.Errorf("got %v want %v", got, want)
		}
	})
}
