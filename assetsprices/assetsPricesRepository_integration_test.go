package assetsprices

import (
	"testing"
)

func TestAssetsPricesRepository(t *testing.T) {
	// load environment variables
	// err := godotenv.Load("../.env")
	// if err != nil {
	// 	fmt.Println(".env file does not exist")
	// }

	// mongoURL := os.Getenv("MONGO_URL")
	// mongoDB := os.Getenv("MONGO_DB")

	// dbClient, err := db.ConnectDB(mongoURL)

	// if err != nil {
	// 	log.Fatal("connecting db", err)
	// }

	// mongoDatabase := dbClient.Database(mongoDB)
	// assetspricesCollection := mongoDatabase.Collection(db.ASSETS_PRICES_COLLECTION)
	// assetspricesRepository := NewRepository(db.NewRepository(assetspricesCollection))

	// t.Run("create and find created asset price", func(t *testing.T) {
	// date := time.Now()
	// value := float32(1000.20)
	// asset := "BTC"

	// err := assetspricesRepository.Create(date, value, asset)

	// if err != nil {
	// 	t.Errorf("error creating asset price %v", err)
	// }

	// filter := bson.D{{"date", date}, {"value", value}, {"asset", asset}}
	// 	filter := bson.D{}

	// 	assetPrice, err := assetsRepository.FindOne(filter)

	// 	if err != nil {
	// 		t.Errorf("error finding asset price %v", err)
	// 	}

	// 	if assetPrice == nil {
	// 		t.Errorf("Not able to find asset price")
	// 	}

	// 	t.Errorf("%v", assetPrice)

	// })

	// t.Run("aggregate prices per date", func(t *testing.T) {
	// 	pipeline := mongo.Pipeline{
	// 		{{"$match", bson.D{{"asset", "BTC"}}}},
	// 		{{
	// 			"$group",
	// 			bson.D{
	// 				{
	// 					"_id", bson.D{
	// 						{"year", bson.D{{"$year", "$date"}}},
	// 						{"month", bson.D{{"$month", "$date"}}},
	// 						{"day", bson.D{{"$dayOfMonth", "$date"}}},
	// 					},
	// 				},
	// 				{"price", bson.D{{"$last", "$value"}}},
	// 			},
	// 		}},
	// 	}

	// 	assetsPrices, err := assetspricesRepository.Aggregate(pipeline)

	// 	if err != nil {
	// 		t.Errorf("Error on aggregating data: %v", err)
	// 	}

	// 	for _, a := range *assetsPrices {
	// 		var ap domain.AssetPrice
	// 		bsonBytes, _ := bson.Marshal(a)
	// 		bson.Unmarshal(bsonBytes, &ap)
	// 		fmt.Printf("%v \n", ap)
	// 	}
	// 	t.Errorf("Fail")
	// })

}
