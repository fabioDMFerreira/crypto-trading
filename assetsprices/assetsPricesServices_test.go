package assetsprices_test

import (
	"testing"
	"time"

	"github.com/fabiodmferreira/crypto-trading/assetsprices"
	"github.com/fabiodmferreira/crypto-trading/mocks"
)

func TestGetNextDateParams(t *testing.T) {
	t.Run("should decrease one day to both dates", func(t *testing.T) {
		dateLayout := "2006-01-02 15:04"

		startDate, _ := time.Parse(dateLayout, "2019-05-04 00:00")
		endDate, _ := time.Parse(dateLayout, "2019-05-04 23:59")

		gotStartDate, gotEndDate := assetsprices.GetDatesPlusOneDay(startDate, endDate)
		wantStartDate, _ := time.Parse(dateLayout, "2019-05-05 00:00")
		wantEndDate, _ := time.Parse(dateLayout, "2019-05-05 23:59")

		if gotStartDate != wantStartDate {
			t.Errorf("startDate: got %v want %v", gotStartDate, wantStartDate)
		}
		if gotEndDate != wantEndDate {
			t.Errorf("endDate: got %v want %v", gotEndDate, wantEndDate)
		}
	})
}

func TestSerializeDate(t *testing.T) {
	t.Run("should return date as string to be used in an URL param", func(t *testing.T) {
		dateLayout := "2006-01-02 15:04"

		date, _ := time.Parse(dateLayout, "2019-05-04 03:00")

		got := assetsprices.SerializeDate(date)
		want := "2019-05-04T03:00"

		if got != want {
			t.Errorf("got %v want %v", got, want)

		}
	})
}

type TransverseSpy struct {
	calls [][]time.Time
}

func (ts *TransverseSpy) handle(startDate, endDate time.Time) error {
	ts.calls = append(ts.calls, []time.Time{startDate, endDate})
	return nil
}

func TestTransverseDatesRanges(t *testing.T) {
	t.Run("should transverse daily from start to end dates", func(t *testing.T) {
		dateLayout := "2006-01-02 15:04"

		startDate, _ := time.Parse(dateLayout, "2019-05-04 00:00")
		endDate, _ := time.Parse(dateLayout, "2019-05-06 00:00")

		transverseSpy := TransverseSpy{}

		assetsprices.TransverseDatesRange(startDate, endDate, transverseSpy.handle)

		got := len(transverseSpy.calls)
		want := 2

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})
}

func TestServiceCreate(t *testing.T) {
	service, repo := NewAssetsPricesService()

	service.Create(time.Now(), 30, "BTC")

	if len(repo.CreateCalls) != 1 {
		t.Errorf("Expected AssetsPricesRepository.Create to be called 1 time")
	}
}

func TestServiceGetLastAssetsPrices(t *testing.T) {
	service, repo := NewAssetsPricesService()

	service.GetLastAssetsPrices("BTC", 10)

	if len(repo.GetLastAssetsPricesCalls) != 1 {
		t.Errorf("Expected AssetsPricesRepository.GetLastAssetsPricesCalls to be called 1 time")
	}
}

func NewAssetsPricesService() (*assetsprices.Service, *mocks.AssetPriceRepositorySpy) {
	assetsPriceRepo := &mocks.AssetPriceRepositorySpy{}

	return assetsprices.NewService(assetsPriceRepo), assetsPriceRepo
}
