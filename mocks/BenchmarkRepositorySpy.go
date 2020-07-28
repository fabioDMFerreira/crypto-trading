package mocks

import "github.com/fabiodmferreira/crypto-trading/domain"

type UpdateBenchmarkArgs struct {
	ID     string
	Output *domain.BenchmarkOutput
}

type BenchmarkRepositorySpy struct {
	FindAllCalls                  int
	InsertOneCalls                []domain.Benchmark
	DeleteByIdCalls               []string
	UpdateBenchmarkCompletedCalls []UpdateBenchmarkArgs
}

func (r *BenchmarkRepositorySpy) FindAll() (*[]domain.Benchmark, error) {
	r.FindAllCalls++
	return &[]domain.Benchmark{}, nil
}

func (r *BenchmarkRepositorySpy) InsertOne(benchmark *domain.Benchmark) error {
	r.InsertOneCalls = append(r.InsertOneCalls, *benchmark)
	return nil
}

func (r *BenchmarkRepositorySpy) DeleteByID(id string) error {
	r.DeleteByIdCalls = append(r.DeleteByIdCalls, id)
	return nil
}

func (r *BenchmarkRepositorySpy) UpdateBenchmarkCompleted(id string, output *domain.BenchmarkOutput) error {
	r.UpdateBenchmarkCompletedCalls = append(r.UpdateBenchmarkCompletedCalls, UpdateBenchmarkArgs{id, output})
	return nil
}
