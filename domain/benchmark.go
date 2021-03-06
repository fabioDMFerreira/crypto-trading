package domain

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
	Buys                [][]float32 `json:"buys"`
	Sells               [][]float32 `json:"sells"`
	SellsPending        int         `json:"sellsPending"`
	FinalAmount         float32     `json:"finalAmount"`
	Assets              *[]Asset    `json:"assets"`
	AssetsAmountPending float32     `json:"assetsAmountPending"`
	AssetsValuePending  float32     `json:"assetsValuePending"`
	LastPrice           float32     `json:"lastPrice"`
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
	ID          primitive.ObjectID `bson:"_id" json:"_id"`
	Input       BenchmarkInput     `json:"input"`
	Output      BenchmarkOutput    `json:"output"`
	Status      string             `json:"status"`
	CreatedAt   time.Time          `json:"createdAt"`
	CompletedAt time.Time          `json:"completedAt"`
}

// BenchmarkResult stores benchmark returned value and possible error
type BenchmarkResult struct {
	Input  *BenchmarkInput
	Output *BenchmarkOutput
	Err    error
}

// BenchmarksRepository stores and fetches benchmarks
type BenchmarksRepository interface {
	FindAll() (*[]Benchmark, error)
	InsertOne(benchmark *Benchmark) error
	DeleteByID(id string) error
	UpdateBenchmarkCompleted(id string, output *BenchmarkOutput) error
}

type BenchmarkService interface {
	Create(input BenchmarkInput) (*Benchmark, error)
	DeleteByID(id string) error
	FindAll() (*[]Benchmark, error)
	BulkRun(inputs []BenchmarkInput, c chan BenchmarkResult)
	Run(input BenchmarkInput, benchmarkID *primitive.ObjectID) (*BenchmarkOutput, error)
	HandleBenchmark(benchmark *Benchmark) error
	GetDataSources() map[string]map[string]string
	AggregateApplicationState(pipeline mongo.Pipeline) (*[]bson.M, error)
}
