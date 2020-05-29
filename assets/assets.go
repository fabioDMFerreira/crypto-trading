package assets

import (
	"time"

	"github.com/fabiodmferreira/crypto-trading/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Asset
type Asset struct {
	ID        primitive.ObjectID `bson:"_id" json:"_id"`
	Amount    float32            `bson:"amount,truncate" json:"amount"`
	BuyTime   time.Time          `json:"buyTime"`
	SellTime  time.Time          `json:"sellTime"`
	BuyPrice  float32            `bson:"buyPrice,truncate" json:"buyPrice"`
	SellPrice float32            `bson:"sellPrice,truncate" json:"sellPrice"`
	Sold      bool               `json:"sold"`
}

// AssetsRepository is the DAO of orders
type AssetsRepository struct {
	collection *mongo.Collection
}

// NewAssetsRepository returns an instance of OrdersRepository
func NewAssetsRepository(collection *mongo.Collection) *AssetsRepository {

	return &AssetsRepository{
		collection,
	}
}

// FindAll returns every order
func (or *AssetsRepository) FindAll() (*[]Asset, error) {
	ctx := db.NewMongoQueryContext()
	cur, err := or.collection.Find(ctx, bson.D{{"sold", false}})

	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)
	var results []Asset
	if err = cur.All(ctx, &results); err != nil {
		return nil, err
	}

	return &results, nil
}

// FindCheaperAssetPrice returns the asset with the lower buy price
func (or *AssetsRepository) FindCheaperAssetPrice() (float32, error) {
	ctx := db.NewMongoQueryContext()

	opts := options.FindOne().SetSort(bson.D{{"buyprice", 1}})
	var foundDocument Asset
	err := or.collection.FindOne(ctx, bson.D{{"sold", false}}, opts).Decode(&foundDocument)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, nil
		}
		return 0, err
	}

	return foundDocument.BuyPrice, nil
}

// Create inserts a new asset in collection
func (or *AssetsRepository) Create(asset *Asset) error {
	ctx := db.NewMongoQueryContext()
	_, err := or.collection.InsertOne(ctx, asset)
	return err
}

func (or *AssetsRepository) Sell(id primitive.ObjectID, price float32) error {
	ctx := db.NewMongoQueryContext()

	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"sellprice", price}, {"sold", true}, {"selltime", time.Now()}}}}
	_, err := or.collection.UpdateOne(ctx, filter, update)

	return err
}