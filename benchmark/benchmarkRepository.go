package benchmark

import (
	"time"

	"github.com/fabiodmferreira/crypto-trading/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Repository stores and returns benchmarks documents
type Repository struct {
	repo domain.Repository
}

// NewRepository returns an instance of benchmark repository
func NewRepository(repo domain.Repository) *Repository {
	return &Repository{repo}
}

// FindAll returns every benchmark
func (r *Repository) FindAll() (*[]domain.Benchmark, error) {
	var benchmarks []domain.Benchmark

	err := r.repo.FindAll(&benchmarks, bson.D{}, nil)

	if err != nil {
		return nil, err
	}

	return &benchmarks, nil
}

// InsertOne creates one benchmark
func (r *Repository) InsertOne(benchmark *domain.Benchmark) error {
	return r.repo.InsertOne(benchmark)
}

// DeleteByID deletes one benchmark
func (r *Repository) DeleteByID(id string) error {
	return r.repo.DeleteByID(id)
}

// UpdateBenchmarkCompleted updates one benchmark
func (r *Repository) UpdateBenchmarkCompleted(id string, output *domain.BenchmarkOutput) error {
	primitiveID, _ := primitive.ObjectIDFromHex(id)

	filter := bson.D{{"_id", primitiveID}}
	update := bson.D{{"$set", bson.D{{"status", "Completed"}, {"output", output}, {"completedat", time.Now()}}}}

	return r.repo.UpdateOne(filter, update)
}
