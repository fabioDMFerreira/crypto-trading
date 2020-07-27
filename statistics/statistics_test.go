package statistics

import (
	"testing"

	"github.com/fabiodmferreira/crypto-trading/domain"
)

func NewMACDMock() *MACDContainer {
	return NewMACDContainer(MACDParams{12, 26, 9}, []float64{})
}

func TestStatistics(t *testing.T) {
	t.Run("RecalculateAverage should recalculate the average", func(t *testing.T) {
		macdMock := NewMACDMock()

		tests := []struct {
			statistics         Statistics
			valueToAdd         float64
			wantAverage        float64
			wantNumberOfPoints int
		}{
			{Statistics{domain.StatisticsOptions{NumberOfPointsHold: 10}, []float64{}, 0, 0, 0, macdMock}, 2, 2, 1},
			{Statistics{domain.StatisticsOptions{NumberOfPointsHold: 10}, []float64{10}, 10, 1, 0, macdMock}, 5, 7.5, 2},
			{Statistics{domain.StatisticsOptions{NumberOfPointsHold: 10}, []float64{2, 20, 80}, 34, 3, 0, macdMock}, 4, 26.5, 4},
			{Statistics{domain.StatisticsOptions{NumberOfPointsHold: 2}, []float64{20, 80}, 50, 2, 0, macdMock}, 40, 60, 2},
			{Statistics{domain.StatisticsOptions{NumberOfPointsHold: 3}, []float64{20, 80}, 50, 2, 0, macdMock}, 50, 50, 3},
		}

		for index, tt := range tests {
			tt.statistics.AddPoint(tt.valueToAdd)

			gotAverage := tt.statistics.average
			wantAverage := tt.wantAverage

			gotNumberOfPoints := tt.statistics.numberOfPoints
			wantNumberOfPoints := tt.wantNumberOfPoints

			if gotAverage != wantAverage {
				t.Errorf("#%v average: got %v want %v", index+1, gotAverage, wantAverage)
			}

			if gotNumberOfPoints != wantNumberOfPoints {
				t.Errorf("#%v number of points: got %v want %v", index+1, gotNumberOfPoints, wantNumberOfPoints)
			}
		}
	})

	t.Run("RecalculateVariance should recalculate the variance", func(t *testing.T) {
		macdMock := NewMACDMock()

		stats := Statistics{domain.StatisticsOptions{NumberOfPointsHold: 10}, []float64{}, 0, 0, 0, macdMock}

		stats.AddPoint(5)
		stats.AddPoint(5)
		stats.AddPoint(10)
		stats.AddPoint(15)
		stats.AddPoint(15)

		got := stats.variance
		want := 25.0

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("RecalculateVariance should recalculate the variance of last points", func(t *testing.T) {
		macdMock := NewMACDMock()

		stats := Statistics{domain.StatisticsOptions{NumberOfPointsHold: 4}, []float64{}, 0, 0, 0, macdMock}

		stats.AddPoint(5)
		stats.AddPoint(5)
		stats.AddPoint(10)
		stats.AddPoint(15)
		stats.AddPoint(15)

		got := float32(stats.variance)
		want := float32(22.916666)

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})
}
