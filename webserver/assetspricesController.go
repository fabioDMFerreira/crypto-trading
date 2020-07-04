package webserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/fabiodmferreira/crypto-trading/domain"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// AssetsPricesController has the handlers of benchmark routes
type AssetsPricesController struct {
	repo domain.AssetPriceRepository
}

// NewAssetsPricesController returns an instance of AssetsPricesController
func NewAssetsPricesController(repo domain.AssetPriceRepository) *AssetsPricesController {
	return &AssetsPricesController{repo}
}

// GetAssetPrices returns prices of the asset between a start date and an end date
func (a *AssetsPricesController) GetAssetPrices(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	queryVars := r.URL.Query()

	if queryVars["startDate"] == nil || queryVars["endDate"] == nil || len(queryVars["startDate"]) == 0 || len(queryVars["endDate"]) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "startDate and endDate parameters are required")
		return
	}

	asset := strings.ToUpper(vars["asset"])

	// TODO: Validate query parameters.

	startDate, _ := time.Parse("2006-01-02T15:04:05", queryVars["startDate"][0])
	endDate, _ := time.Parse("2006-01-02T15:04:05", queryVars["endDate"][0])

	days := endDate.Sub(startDate).Hours() / 24
	// fmt.Printf("%v %v", startDate, endDate)

	var pipelineOptions mongo.Pipeline

	if days > 30 {
		pipelineOptions = mongo.Pipeline{
			{{
				"$match",
				bson.D{
					{"asset", asset},
					{"date", bson.D{{"$gte", startDate}}},
					{"date", bson.D{{"$lte", endDate}}},
				},
			}},
			{{
				"$group",
				bson.D{
					{
						"_id", bson.D{
							{"year", bson.D{{"$year", "$date"}}},
							{"month", bson.D{{"$month", "$date"}}},
							{"day", bson.D{{"$dayOfMonth", "$date"}}},
						}},
					{"price", bson.D{{"$last", "$value"}}},
				},
			}},
		}
	} else if days > 5 {
		pipelineOptions = mongo.Pipeline{
			{{
				"$match",
				bson.D{
					{"asset", asset},
					{"date", bson.D{{"$gte", startDate}}},
					{"date", bson.D{{"$lte", endDate}}},
				},
			}},
			{{
				"$group",
				bson.D{
					{
						"_id", bson.D{
							{"year", bson.D{{"$year", "$date"}}},
							{"month", bson.D{{"$month", "$date"}}},
							{"day", bson.D{{"$dayOfMonth", "$date"}}},
							{"hour", bson.D{{"$hour", "$date"}}},
						}},
					{"price", bson.D{{"$last", "$value"}}},
				},
			}},
		}
	} else {
		pipelineOptions = mongo.Pipeline{
			{{
				"$match",
				bson.D{
					{"asset", asset},
					{"date", bson.D{{"$gte", startDate}}},
					{"date", bson.D{{"$lte", endDate}}},
				},
			}},
			{{
				"$group",
				bson.D{
					{
						"_id", bson.D{
							{"year", bson.D{{"$year", "$date"}}},
							{"month", bson.D{{"$month", "$date"}}},
							{"day", bson.D{{"$dayOfMonth", "$date"}}},
							{"hour", bson.D{{"$hour", "$date"}}},
							{"minute", bson.D{{"$minute", "$date"}}},
						}},
					{"price", bson.D{{"$last", "$value"}}},
				},
			}},
		}
	}
	assetsPrices, err := a.repo.Aggregate(pipelineOptions)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	json.NewEncoder(w).Encode(*assetsPrices)
}
