package benchmark

import (
	"github.com/fabiodmferreira/crypto-trading/domain"
)

// RepositoryInMemory stores benchmarks in memory
type RepositoryInMemory struct {
	Benchmarks []domain.Benchmark
}

// NewRepositoryInMemory returns an instance of RepositoryInMemory
func NewRepositoryInMemory() *RepositoryInMemory {
	return &RepositoryInMemory{[]domain.Benchmark{}}
}

// FindAll returns all benchmarks stored
func (r *RepositoryInMemory) FindAll() (*[]domain.Benchmark, error) {
	return &r.Benchmarks, nil
}

// InsertOne creates a benchmark and stores it in a data structure
func (r *RepositoryInMemory) InsertOne(benchmark *domain.Benchmark) error {
	r.Benchmarks = append(r.Benchmarks, *benchmark)
	return nil
}

// DeleteByID removes a benchmark from store
func (r *RepositoryInMemory) DeleteByID(id string) error {
	for index, b := range r.Benchmarks {
		if b.ID.String() == id {
			r.Benchmarks[index], r.Benchmarks[len(r.Benchmarks)-1] = r.Benchmarks[len(r.Benchmarks)-1], r.Benchmarks[index]
			r.Benchmarks = r.Benchmarks[:len(r.Benchmarks)-1]
			break
		}
	}

	return nil
}

// UpdateBenchmarkCompleted updates benchmark status
func (r *RepositoryInMemory) UpdateBenchmarkCompleted(id string, output *domain.BenchmarkOutput) error {
	for _, b := range r.Benchmarks {
		if b.ID.String() == id {
			b.Status = "Completed"
			break
		}
	}
	return nil
}
