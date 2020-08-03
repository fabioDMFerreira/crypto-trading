package assets

import (
	"time"

	"github.com/fabiodmferreira/crypto-trading/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Asset is an alias of domain.Asset
type Asset = domain.Asset

// Repository is the DAO of orders
type Repository struct {
	repo domain.Repository
}

// NewRepository returns an instance of OrdersRepository
func NewRepository(repo domain.Repository) *Repository {

	return &Repository{
		repo,
	}
}

// FindPendingAssets returns all assets that weren't sold
func (or *Repository) FindPendingAssets(accountID primitive.ObjectID) (*[]Asset, error) {

	query := bson.M{"sold": false, "accountID": accountID}

	var results []Asset
	err := or.repo.FindAll(&results, query, nil)

	if err != nil {
		return nil, err
	}

	return &results, nil
}

// FindOne returns one asset
func (or *Repository) FindOne(filter interface{}) (*Asset, error) {
	var asset Asset

	err := or.repo.FindOne(&asset, filter, options.FindOne())

	if err != nil {
		return nil, err
	}

	return &asset, nil
}

// FindAll returns every order
func (or *Repository) FindAll(accountID primitive.ObjectID) (*[]Asset, error) {
	query := bson.M{"accountID": accountID}

	var results []Asset
	err := or.repo.FindAll(&results, query, nil)

	if err != nil {
		return nil, err
	}

	return &results, nil
}

// FindCheaperAssetPrice returns the asset with the lower buy price
func (or *Repository) FindCheaperAssetPrice(accountID primitive.ObjectID) (float32, error) {
	opts := options.FindOne().SetSort(bson.M{"buyPrice": 1})
	var foundDocument Asset
	err := or.repo.FindOne(&foundDocument, bson.M{"sold": false, "accountID": accountID}, opts)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, nil
		}
		return 0, err
	}

	return foundDocument.BuyPrice, nil
}

// Create inserts a new asset in collection
func (or *Repository) Create(asset *Asset) error {
	return or.repo.InsertOne(asset)
}

// Sell updates asset sell fields
func (or *Repository) Sell(id primitive.ObjectID, price float32, sellTime time.Time) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"sellPrice": price, "sold": true, "selltime": sellTime}}
	err := or.repo.UpdateOne(filter, update)

	return err
}

// GetBalance returns the assets balance based on buys and sells
func (ar *Repository) GetBalance(accountID primitive.ObjectID, startDate, endDate time.Time) (float32, error) {
	filter := bson.M{"sold": false, "accountID": accountID, "buytime": bson.M{"$gte": startDate, "$lte": endDate}}
	var assetsBought []Asset
	err := ar.repo.FindAll(&assetsBought, filter, nil)

	if err != nil {
		return 0, err
	}

	filter = bson.M{"sold": true, "selltime": bson.M{"$gte": startDate, "$lte": endDate}}
	var assetsSold []Asset
	err = ar.repo.FindAll(&assetsSold, filter, nil)

	if err != nil {
		return 0, err
	}

	var balance float32

	for _, asset := range assetsSold {
		balance += asset.Amount * asset.SellPrice
	}

	for _, asset := range assetsBought {
		balance -= asset.Amount * asset.BuyPrice
	}

	return balance, nil
}

// CheckAssetWithCloserPriceExists checks whether an asset that has the same price within limits defined exists
func (ar *Repository) CheckAssetWithCloserPriceExists(accountID primitive.ObjectID, price, limit float32) (bool, error) {
	lowerLimit := price - (price * limit)
	upperLimit := price + (price * limit)

	filter := bson.M{"sold": false, "accountID": accountID, "buyPrice": bson.M{"$gte": lowerLimit, "$lte": upperLimit}}

	var assets []Asset
	err := ar.repo.FindAll(&assets, filter, nil)

	if err != nil {
		return false, err
	}

	if len(assets) > 0 {
		return true, nil
	}

	return false, nil
}
