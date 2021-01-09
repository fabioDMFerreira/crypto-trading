package utils

import (
	"encoding/json"
	"math"
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

// merge merges the two JSON-marshalable values x1 and x2,
// preferring x1 over x2 except where x1 and x2 are
// JSON objects, in which case the keys from both objects
// are included and their values merged recursively.
//
// It returns an error if x1 or x2 cannot be JSON-marshaled.
func Merge(x1, x2 interface{}) (interface{}, error) {
	data1, err := json.Marshal(x1)
	if err != nil {
		return nil, err
	}
	data2, err := json.Marshal(x2)
	if err != nil {
		return nil, err
	}
	var j1 interface{}
	err = json.Unmarshal(data1, &j1)
	if err != nil {
		return nil, err
	}
	var j2 interface{}
	err = json.Unmarshal(data2, &j2)
	if err != nil {
		return nil, err
	}
	return merge1(j1, j2), nil
}

func merge1(x1, x2 interface{}) interface{} {
	switch x1 := x1.(type) {
	case map[string]interface{}:
		x2, ok := x2.(map[string]interface{})
		if !ok {
			return x1
		}
		for k, v2 := range x2 {
			if v1, ok := x1[k]; ok {
				x1[k] = merge1(v1, v2)
			} else {
				x1[k] = v2
			}
		}
	case nil:
		// merge(nil, map[string]interface{...}) -> map[string]interface{...}
		x2, ok := x2.(map[string]interface{})
		if ok {
			return x2
		}
	}
	return x1
}

// RoundFloorDecimals rounds number specified of decimal places
func RoundFloorDecimals(n float32, decimalPlaces int) float32 {
	factor := math.Pow10(decimalPlaces)

	return float32(math.Floor(float64(n)*factor) / factor)
}
