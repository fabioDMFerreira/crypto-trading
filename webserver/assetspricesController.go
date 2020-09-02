package webserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/fabiodmferreira/crypto-trading/domain"
	"github.com/fabiodmferreira/crypto-trading/utils"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// AssetsPricesController has the assets prices routes handlers
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

	var pipelineOptions mongo.Pipeline

	groupByDatesClause := utils.GetGroupByDatesIDClause(startDate, endDate)

	pipelineOptions = mongo.Pipeline{
		{primitive.E{
			Key: "$match",
			Value: bson.M{
				"asset": asset,
				"date":  bson.M{"$gte": startDate, "$lte": endDate},
			},
		}},
		{
			primitive.E{
				Key: "$group",
				Value: bson.D{
					primitive.E{
						Key:   "_id",
						Value: groupByDatesClause,
					},
					primitive.E{
						Key:   "price",
						Value: bson.M{"$last": "$c"},
					},
				},
			},
		},
	}

	assetsPrices, err := a.repo.Aggregate(pipelineOptions)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	json.NewEncoder(w).Encode(*assetsPrices)
}
