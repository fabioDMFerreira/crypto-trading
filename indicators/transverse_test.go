package indicators_test

import (
	"reflect"
	"testing"

	"github.com/fabiodmferreira/crypto-trading/indicators"
)

func TestSubSlices(t *testing.T) {
	a := []float64{5, 4, 3, 2, 1}
	b := []float64{1, 1, 1, 1, 1}

	got := indicators.SubSlices(a, b)
	want := []float64{4, 3, 2, 1, 0}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestTransverseEMA(t *testing.T) {
	slice := []float64{1, 6, 2, 3, 4, 6, 1}
	period := 5

	got := indicators.TransverseEMA(slice, period)
	want := []float64{1, 2.666666666666667, 2.4444444444444446, 2.6296296296296298, 3.08641975308642, 4.057613168724281, 3.038408779149521}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestTransverseMACD(t *testing.T) {
	slice := []float64{1, 2, 3, 4, 8, 1, 2, 43, 7}

	got1, got2 := indicators.TransverseMACD(slice, 12, 26, 9)
	want1, want2 := []float64{0, 0.07977207977207978, 0.22113456871291648, 0.40914068319690045, 0.870864668240513, 0.6642852870839073, 0.5746372216525617, 3.768506831680396, 3.356085508121801}, []float64{0, 0.01595441595441596, 0.056990446506116066, 0.12742049384427295, 0.27610932872352095, 0.3537445203955982, 0.3979230606469909, 1.0720398148536718, 1.5288489535072978}

	if !reflect.DeepEqual(got1, want1) {
		t.Errorf("macd: got %v want %v", got1, want1)
	}

	if !reflect.DeepEqual(got2, want2) {
		t.Errorf("ema9: got %v want %v", got2, want2)
	}
}
