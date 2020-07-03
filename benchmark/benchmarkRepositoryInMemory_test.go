package benchmark

import (
	"reflect"
	"testing"
	"time"

	"github.com/fabiodmferreira/crypto-trading/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestBenchmarkRepositoryInMemory(t *testing.T) {
	t.Run("DeleteByID should remove an element", func(t *testing.T) {
		br := NewRepositoryInMemory()

		b1 := domain.Benchmark{ID: primitive.NewObjectID(), CreatedAt: time.Now()}
		b2 := domain.Benchmark{ID: primitive.NewObjectID(), CreatedAt: time.Now()}

		br.Benchmarks = []domain.Benchmark{
			b1,
			b2,
		}

		br.DeleteByID(b1.ID.String())

		want := []domain.Benchmark{
			b2,
		}
		got := br.Benchmarks

		if reflect.DeepEqual(got, want) != true {
			t.Errorf("got %v want %v", got, want)
		}

	})
}
