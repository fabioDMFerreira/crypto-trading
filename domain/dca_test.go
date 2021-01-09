package domain_test

import (
	"reflect"
	"testing"

	"github.com/fabiodmferreira/crypto-trading/domain"
)

func TestDCAJob(t *testing.T) {

	t.Run("GetFiatCoinsAmount should return fiat amounts that match the options proportion", func(t *testing.T) {
		dcaJob := domain.DCAJob{
			Options: domain.DCAJobOptions{
				TotalFIATAmount: 200,
				CoinsProportion: map[string]float32{
					"BTC": 50,
					"ETH": 30,
					"ADA": 20,
				},
			},
		}

		got := dcaJob.GetFiatCoinsAmount()
		want := map[string]float32{
			"BTC": 100,
			"ETH": 60,
			"ADA": 40,
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v,want %v", got, want)
		}

		checkMapValuesSumIsLessThan(t, got, dcaJob.Options.TotalFIATAmount)
	})

	t.Run("GetFiatCoinsAmount should return fiat amounts with maximum 2 decimal points", func(t *testing.T) {
		dcaJob := domain.DCAJob{
			Options: domain.DCAJobOptions{
				TotalFIATAmount: 200,
				CoinsProportion: map[string]float32{
					"BTC": 5340,
					"ETH": 3019,
					"ADA": 2029,
				},
			},
		}

		got := dcaJob.GetFiatCoinsAmount()
		want := map[string]float32{
			"BTC": 102.81,
			"ETH": 58.12,
			"ADA": 39.06,
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v,want %v", got, want)
		}

		checkMapValuesSumIsLessThan(t, got, dcaJob.Options.TotalFIATAmount)
	})

}

func checkMapValuesSumIsLessThan(t *testing.T, amounts map[string]float32, total float32) {
	var amountsSum float32

	for _, value := range amounts {
		amountsSum += value
	}

	if amountsSum > total {
		t.Errorf("Map values should not be greater than total specified: amountsSum:%f > totalValue:%f", amountsSum, total)
	}

}
