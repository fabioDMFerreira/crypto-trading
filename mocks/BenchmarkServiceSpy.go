package mocks

import (
	"github.com/fabiodmferreira/crypto-trading/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BenchmarkServiceSpy struct {
	CreateCalls                    []domain.BenchmarkInput
	DeleteByIDCalls                []string
	FindAllCalls                   int
	BulkRunCalls                   [][]interface{}
	RunCalls                       [][]interface{}
	HandleBenchmarkCalls           []domain.Benchmark
	GetDataSourcesCalls            int
	AggregateApplicationStateCalls []interface{}
}

func (s *BenchmarkServiceSpy) Create(input domain.BenchmarkInput) (*domain.Benchmark, error) {
	s.CreateCalls = append(s.CreateCalls, input)
	return &domain.Benchmark{}, nil
}

func (s *BenchmarkServiceSpy) DeleteByID(id string) error {
	s.DeleteByIDCalls = append(s.DeleteByIDCalls, id)
	return nil
}

func (s *BenchmarkServiceSpy) FindAll() (*[]domain.Benchmark, error) {
	s.FindAllCalls++
	return &[]domain.Benchmark{}, nil
}

func (s *BenchmarkServiceSpy) BulkRun(inputs []domain.BenchmarkInput, c chan domain.BenchmarkResult) {
	s.BulkRunCalls = append(s.BulkRunCalls, []interface{}{inputs, c})
}

func (s *BenchmarkServiceSpy) Run(input domain.BenchmarkInput, benchmarkID *primitive.ObjectID) (*domain.BenchmarkOutput, error) {
	s.RunCalls = append(s.RunCalls, []interface{}{input, benchmarkID})

	return &domain.BenchmarkOutput{}, nil
}

func (s *BenchmarkServiceSpy) HandleBenchmark(benchmark *domain.Benchmark) error {
	s.HandleBenchmarkCalls = append(s.HandleBenchmarkCalls, *benchmark)
	return nil
}

func (s *BenchmarkServiceSpy) GetDataSources() map[string]map[string]string {
	s.GetDataSourcesCalls++
	return map[string]map[string]string{}
}

func (s *BenchmarkServiceSpy) AggregateApplicationState(pipeline mongo.Pipeline) (*[]bson.M, error) {
	s.AggregateApplicationStateCalls = append(s.AggregateApplicationStateCalls, pipeline)

	return &[]bson.M{}, nil
}
