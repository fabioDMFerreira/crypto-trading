package domain

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BenchmarkInput needed to run benchmark
type BenchmarkInput struct {
	DecisionMakerOptions DecisionMakerOptions `json:"decisionMakerOptions"`
	StatisticsOptions    StatisticsOptions    `json:"statisticsOptions"`
	CollectorOptions     CollectorOptions     `json:"collectorOptions"`
	AccountInitialAmount float64              `json:"accountInitialAmount"`
	DataSourceFilePath   string               `json:"dataSourceFilePath"`
	Asset                string               `json:"asset"`
}

// BenchmarkOutput is the output of the benchmark
type BenchmarkOutput struct {
	Buys         [][]float32 `json:"buys"`
	Sells        [][]float32 `json:"sells"`
	SellsPending int         `json:"sellsPending"`
	FinalAmount  float32     `json:"finalAmount"`
	Assets       *[]Asset    `json:"assets"`
	Balances     [][]float32 `json:"balances"`
}

// String displays Output formatted
func (o *BenchmarkOutput) String() {
	fmt.Println("======")
	fmt.Printf("Buys %v\n", o.Buys)
	fmt.Printf("Sells %v\n", o.Sells)
	fmt.Printf("Sells Pending %v\n", o.SellsPending)
	fmt.Printf("Final amount %v\n", o.FinalAmount)
	fmt.Println("=======")
}

// Benchmark stores inputs and outputs of a benchmark
type Benchmark struct {
	ID        primitive.ObjectID `bson:"_id" json:"_id"`
	Input     BenchmarkInput     `json:"input"`
	Output    BenchmarkOutput    `json:"output"`
	Status    string             `json:"status"`
	CreatedAt time.Time          `json:"createdAt"`
}

// BenchmarksRepository stores and fetches benchmarks
type BenchmarksRepository interface {
	FindAll() (*[]Benchmark, error)
	InsertOne(benchmark *Benchmark) error
	DeleteByID(id string) error
	UpdateBenchmarkCompleted(id string, output *BenchmarkOutput) error
}
