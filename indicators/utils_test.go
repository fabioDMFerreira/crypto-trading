package indicators_test

import (
	"testing"

	"github.com/fabiodmferreira/crypto-trading/indicators"
)

func TestAverage(t *testing.T) {
	t.Run("should return average", func(t *testing.T) {
		tests := []struct {
			sample []float64
			want   float64
		}{
			{[]float64{1, 2, 3}, 2},
			{[]float64{5, 5, 10, 15, 15}, 10},
			{[]float64{1}, 1},
			{[]float64{}, 0},
			{[]float64{1, -1}, 0},
		}

		for index, tt := range tests {
			got := indicators.Average(tt.sample)
			want := tt.want

			if got != want {
				t.Errorf("#%v: got %v want %v", index+1, got, want)
			}
		}
	})
}

func TestVariance(t *testing.T) {
	t.Run("should return variance", func(t *testing.T) {
		tests := []struct {
			sample []float64
			want   float64
		}{
			{[]float64{1, 2, 3}, 1},
			{[]float64{5, 5, 10, 15, 15}, 25},
			{[]float64{1}, 0},
			{[]float64{}, 0},
		}

		for index, tt := range tests {
			got := indicators.Variance(tt.sample)
			want := tt.want

			if got != want {
				t.Errorf("#%v: got %v want %v", index+1, got, want)
			}
		}
	})

	t.Run("should return variance using online variance algorithm", func(t *testing.T) {
		tests := []struct {
			sample []float64
			want   float64
		}{
			{[]float64{1, 2, 3}, 1},
			{[]float64{5, 5, 10, 15, 15}, 25},
			{[]float64{1}, 0},
			{[]float64{}, 0},
		}

		for index, tt := range tests {
			got := indicators.Variance(tt.sample)
			want := tt.want

			if got != want {
				t.Errorf("#%v: got %v want %v", index+1, got, want)
			}
		}
	})
}
