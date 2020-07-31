package assetsprices_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/fabiodmferreira/crypto-trading/assetsprices"
	"github.com/fabiodmferreira/crypto-trading/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type HTTPSpy struct {
	GetCalls []string
}

func (h *HTTPSpy) Get(url string) (resp *http.Response, err error) {
	h.GetCalls = append(h.GetCalls, url)

	response := "{\"statusCode\":200,\"message\":\"OK\",\"data\":{\"iso\":\"BTC\",\"name\":\"Bitcoin\",\"slug\":\"bitcoin\",\"interval\":\"1-minute\",\"entries\":[[1596183599999,11201.3598739665],[1596183659999,11209.1649879131]]}}"

	resp = &http.Response{
		Body: ioutil.NopCloser(bytes.NewBufferString(response)),
	}

	return resp, nil
}

func TestCoindeskRemoteSource(t *testing.T) {
	t.Run("should return assets prices", func(t *testing.T) {
		coindeskRemoteSource, httpSpy := setupCoindeskRemoteSource()

		got, err := coindeskRemoteSource.FetchRemoteAssetsPrices(time.Now(), time.Now(), "BTC")
		want := []bson.M{
			{
				"asset": "BTC",
				"date":  time.Unix(1596183599999/1000, 0),
				"value": float32(11201.3598739665) * utils.DollarEuroRate,
			}, {
				"asset": "BTC",
				"date":  time.Unix(1596183659999/1000, 0),
				"value": float32(11209.1649879131) * utils.DollarEuroRate,
			},
		}

		if err != nil {
			t.Errorf("not expected next error %v", err)
		}

		if reflect.DeepEqual(*got, want) != true {
			t.Errorf("got %#v want %#v", got, want)
		}

		if len(httpSpy.GetCalls) != 1 {
			t.Errorf("Expected fetchCoindeskData to be called 1 time")
		}

	})
}

func setupCoindeskRemoteSource() (*assetsprices.CoindeskRemoteSource, *HTTPSpy) {
	httpSpy := &HTTPSpy{}

	return assetsprices.NewCoindeskRemoteSource(httpSpy.Get), httpSpy
}
