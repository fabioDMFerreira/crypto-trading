package assetsprices_test

import (
	"testing"
	"time"

	"github.com/fabiodmferreira/crypto-trading/assetsprices"
	"github.com/fabiodmferreira/crypto-trading/mocks"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestRepositoryFindAll(t *testing.T) {
	assetspricesRepository, repository := setupAssetsPricesRepository()

	assetspricesRepository.FindAll("all")

	got := len(repository.FindAllCalls)
	want := 1

	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestRepositoryAggregate(t *testing.T) {
	assetspricesRepository, repository := setupAssetsPricesRepository()

	assetspricesRepository.Aggregate(mongo.Pipeline{{}})

	got := len(repository.AggregateCalls)
	want := 1

	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestRepositoryFindOne(t *testing.T) {
	assetspricesRepository, repository := setupAssetsPricesRepository()

	assetspricesRepository.FindOne("one")

	got := len(repository.FindOneCalls)
	want := 1

	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestRepositoryCreate(t *testing.T) {
	assetspricesRepository, repository := setupAssetsPricesRepository()

	assetspricesRepository.Create(time.Now(), 1000, "BTC")

	got := len(repository.InsertOneCalls)
	want := 1

	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestRepositoryGetLastAssetsPrices(t *testing.T) {
	assetspricesRepository, repository := setupAssetsPricesRepository()

	assetspricesRepository.GetLastAssetsPrices("BTC", 5)

	got := len(repository.FindAllCalls)
	want := 1

	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestRepositoryBulkCreate(t *testing.T) {
	assetspricesRepository, repository := setupAssetsPricesRepository()

	assetspricesRepository.BulkCreate(&[]bson.M{})

	got := len(repository.BulkCreateCalls)
	want := 1

	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func setupAssetsPricesRepository() (*assetsprices.Repository, *mocks.RepositorySpy) {
	repository := &mocks.RepositorySpy{}

	return assetsprices.NewRepository(repository), repository
}
