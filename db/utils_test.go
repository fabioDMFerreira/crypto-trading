package db_test

import (
	"testing"

	"github.com/fabiodmferreira/crypto-trading/db"
	"go.mongodb.org/mongo-driver/bson"
)

type RepoStub struct {
	calls int
}

func (r *RepoStub) BulkCreate(documents *[]bson.M) error {
	r.calls++
	return nil
}

func TestBatchBulkCreate(t *testing.T) {
	t.Run("should call bulk create as many times as needed", func(t *testing.T) {
		repo := &RepoStub{}
		documents := []bson.M{}

		for i := 0; i < 100; i++ {
			documents = append(documents, bson.M{"index": i})
		}

		err := db.BatchBulkCreate(repo.BulkCreate, &documents, 10)

		if err != nil {
			t.Errorf("%v", err)
		}

		got := repo.calls
		want := 10

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}

	})
}
