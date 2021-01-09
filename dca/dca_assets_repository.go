package dca

import (
	"github.com/fabiodmferreira/crypto-trading/domain"
	"go.mongodb.org/mongo-driver/bson"
)

// AssetsRepository fetches and stores dca assets bought
type AssetsRepository struct {
	repo domain.Repository
}

// NewAssetsRepository returns an instance of dca AssetsRepository
func NewAssetsRepository(repo domain.Repository) *AssetsRepository {
	return &AssetsRepository{repo}
}

// Save persists dca job in database
func (r *AssetsRepository) Save(dcaAsset *domain.DCAAsset) error {
	return r.repo.InsertOne(dcaAsset)
}

// FindAll fetches existing dca jobs
func (r *AssetsRepository) FindAll() (*[]domain.DCAAsset, error) {
	var results []domain.DCAAsset
	err := r.repo.FindAll(&results, bson.M{}, nil)

	if err != nil {
		return nil, err
	}

	return &results, nil
}
