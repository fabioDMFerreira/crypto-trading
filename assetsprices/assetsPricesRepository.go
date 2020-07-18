package assetsprices

import (
	"time"

	"github.com/fabiodmferreira/crypto-trading/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AssetPrice = domain.AssetPrice

// Repository stores and gets assets prices
type Repository struct {
	repo domain.Repository
}

// NewRepository returns an instance of AssetsPricesRepository
func NewRepository(repo domain.Repository) *Repository {

	return &Repository{
		repo,
	}
}

// FindAll returns assets prices
func (r *Repository) FindAll(filter interface{}) (*[]domain.AssetPrice, error) {
	var results []domain.AssetPrice
	err := r.repo.FindAll(&results, filter, nil)

	if err != nil {
		return nil, err
	}

	return &results, nil
}

// Aggregate returns assets prices aggregated
func (r *Repository) Aggregate(pipeline mongo.Pipeline) (*[]bson.M, error) {
	var results []bson.M
	err := r.repo.Aggregate(&results, pipeline)

	if err != nil {
		return nil, err
	}

	return &results, nil
}

// FindOne returns an asset price
func (r *Repository) FindOne(filter interface{}) (interface{}, error) {
	var assetPrice interface{}

	err := r.repo.FindOne(&assetPrice, filter, options.FindOne())

	if err != nil {
		return nil, err
	}

	return assetPrice, nil
}

// Create stores an asset price
func (r *Repository) Create(date time.Time, value float32, asset string) error {
	filter := bson.D{{"date", date}, {"value", value}, {"asset", asset}}

	assets, err := r.FindAll(filter)

	if len(*assets) != 0 {
		return err
	}

	assetPrice := domain.AssetPrice{ID: primitive.NewObjectID(), Date: date, Value: value, Asset: asset}

	return r.repo.InsertOne(assetPrice)
}

// GetLastAssetsPrices return the last asset price stored in DB
func (r *Repository) GetLastAssetsPrices(asset string, limit int) (*[]domain.AssetPrice, error) {
	opts := options.Find().SetSort(bson.D{{"date", -1}}).SetLimit(int64(limit))
	var foundDocument []AssetPrice
	err := r.repo.FindAll(&foundDocument, bson.M{"asset": asset}, opts)

	if err != nil {
		return nil, err
	}

	return &foundDocument, nil
}

// BulkCreate creates multiple assets prices
func (r *Repository) BulkCreate(documents *[]bson.M) error {
	return r.repo.BulkCreate(documents)
}
