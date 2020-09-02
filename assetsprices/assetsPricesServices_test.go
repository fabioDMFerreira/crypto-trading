package assetsprices_test

import (
	"testing"
	"time"

	"github.com/fabiodmferreira/crypto-trading/assetsprices"
	"github.com/fabiodmferreira/crypto-trading/domain"
	"github.com/fabiodmferreira/crypto-trading/mocks"
	"github.com/golang/mock/gomock"
	"go.mongodb.org/mongo-driver/bson"
)

type FetchRemotePrice struct {
	Calls []interface{}
}

func (f *FetchRemotePrice) fetch(startDate, endDate time.Time, asset string) (*[]bson.M, error) {
	f.Calls = append(f.Calls, []interface{}{startDate, endDate, asset})

	return &[]bson.M{{}, {}}, nil
}

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
	service, repo, _ := NewAssetsPricesService(t)

	ohlc := domain.OHLC{Close: 30, Time: time.Now()}

	repo.EXPECT().Create(&ohlc, "BTC").Return(nil).Times(1)

	got := service.Create(&ohlc, "BTC")
	var want error

	if got != want {
		t.Errorf("Expected AssetsPricesRepository.Create to be called 1 time")
	}
}

func TestServiceGetLastAssetsPrices(t *testing.T) {
	service, repo, _ := NewAssetsPricesService(t)

	repo.EXPECT().GetLastAssetsPrices("BTC", 10).Return(&[]domain.AssetPrice{}, nil).Times(1)

	assetsPrices, err := service.GetLastAssetsPrices("BTC", 10)

	if err != nil {
		t.Errorf("Should not return error")
	}

	if len(*assetsPrices) != 0 {
		t.Errorf("Expected assetsPrices to be empty")
	}
}

func TestServiceFetchAndStore(t *testing.T) {
	service, assetsPriceRepo, fetchRemotePrices := NewAssetsPricesService(t)

	assetsPriceRepo.EXPECT().GetLastAssetsPrices("BTC", 1).Return(&[]domain.AssetPrice{{Date: time.Date(2020, 4, 1, 0, 0, 0, 0, time.UTC)}}, nil).Times(1)
	assetsPriceRepo.EXPECT().BulkCreate(gomock.Any()).Return(nil).Times(2)

	service.FetchAndStoreAssetPrices("BTC", time.Date(2020, 4, 3, 0, 0, 0, 0, time.UTC))

	got := len(fetchRemotePrices.Calls)
	want := 2

	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func NewAssetsPricesService(t *testing.T) (*assetsprices.Service, *mocks.MockAssetPriceRepository, *FetchRemotePrice) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	assetsPriceRepo := mocks.NewMockAssetPriceRepository(ctrl)
	fetchRemotePrices := &FetchRemotePrice{}
	return assetsprices.NewService(assetsPriceRepo, fetchRemotePrices.fetch), assetsPriceRepo, fetchRemotePrices
}
