package statistics_test

import (
	"testing"

	"github.com/fabiodmferreira/crypto-trading/statistics"
)

func TestEMA(t *testing.T) {
	t.Run("should return right EMA", func(t *testing.T) {
		got := float32(statistics.EMA(461.14, 12, 449.823954868254))
		want := float32(451.564884888523)

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})
}

func TestMACDContainerAddPoint(t *testing.T) {
	t.Run("after adding 12 points should add first point to fast EMA", func(t *testing.T) {
		params := statistics.MACDParams{12, 26, 9}
		holdPoints := []float64{
			459.99,
			448.85,
			446.06,
			450.81,
			442.8,
			448.97,
			444.57,
			441.4,
			430.47,
			420.05,
			431.14,
		}

		mc := statistics.NewMACDContainer(params, holdPoints)

		mc.AddPoint(425.66)

		if len(mc.FastEMA) != 1 {
			t.Errorf("Fast EMA should have one point")
		}

		got := float32(mc.FastEMA[0])
		want := float32(440.8975)
		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("after adding 26 points should add first point to slow EMA and MACD", func(t *testing.T) {
		params := statistics.MACDParams{12, 26, 9}
		holdPoints := []float64{
			459.99,
			448.85,
			446.06,
			450.81,
			442.8,
			448.97,
			444.57,
			441.4,
			430.47,
			420.05,
			431.14,
			425.66,
			430.58,
			431.72,
			437.87,
			428.43,
			428.35,
			432.5,
			443.66,
			455.72,
			454.49,
			452.08,
			452.73,
			461.91,
			463.58,
		}
		fastEma := []float64{
			440.8975,
			439.310192307692,
			438.142470414201,
			438.100551888939,
			436.612774675256,
			435.34157857137,
			434.904412637313,
			436.251426077726,
			439.246591296537,
			441.59173109707,
			443.20531092829,
			444.670647708553,
			447.322855753391,
			449.823954868254,
		}
		mc := statistics.NewMACDContainer(params, holdPoints, fastEma)

		mc.AddPoint(461.14)

		if len(mc.SlowEMA) != 1 {
			t.Errorf("Slow EMA should have one point")
		}

		got := float32(mc.SlowEMA[0])
		want := float32(443.289615384615)
		if got != want {
			t.Errorf("got %v want %v", got, want)
		}

		gotLastFastEMA := float32(mc.FastEMA[len(mc.FastEMA)-1])
		wantLastFastEMA := float32(451.564884888523)
		if gotLastFastEMA != wantLastFastEMA {
			t.Errorf("got %v want %v", gotLastFastEMA, wantLastFastEMA)
		}

		gotMACD := float32(mc.MACD[0])
		wantMACD := float32(8.275269503908)

		if gotMACD != wantMACD {
			t.Errorf("got %v want %v", gotMACD, wantMACD)
		}
	})

	t.Run("after adding 35 points should add first point to lag EMA", func(t *testing.T) {
		params := statistics.MACDParams{12, 26, 9}
		holdPoints := []float64{
			459.99,
			448.85,
			446.06,
			450.81,
			442.8,
			448.97,
			444.57,
			441.4,
			430.47,
			420.05,
			431.14,
			425.66,
			430.58,
			431.72,
			437.87,
			428.43,
			428.35,
			432.5,
			443.66,
			455.72,
			454.49,
			452.08,
			452.73,
			461.91,
			463.58,
			461.14,
			452.08,
			442.66,
			428.91,
			429.79,
			431.99,
			427.72,
			423.2,
		}
		fastEma := []float64{
			440.8975,
			439.310192307692,
			438.142470414201,
			438.100551888939,
			436.612774675256,
			435.34157857137,
			434.904412637313,
			436.251426077726,
			439.246591296537,
			441.59173109707,
			443.20531092829,
			444.670647708553,
			447.322855753391,
			449.823954868254,
			451.564884888523,
			451.644133367212,
			450.261959003026,
			446.97704223333,
			444.332881889741,
			442.433976983627,
			440.170288216915,
			437.559474645082,
		}
		slowEma := []float64{
			443.289615384615,
			443.940754985755,
			443.845884246069,
			442.739522450064,
			441.780298564874,
			441.055091263772,
			440.067306725715,
			438.817876597884,
		}
		macd := []float64{
			8.275269503908,
			7.703378381457,
			6.416074756957,
			4.237519783266,
			2.552583324867,
			1.378885719855,
			0.1029814912,
			-1.258401952802,
		}
		mc := statistics.NewMACDContainer(params, holdPoints, fastEma, slowEma, macd)

		mc.AddPoint(426.21)

		gotLastLagEMA := float32(mc.LagEMA[len(mc.LagEMA)-1])
		wantLastLagEMA := float32(3.03752586873489)
		if gotLastLagEMA != wantLastLagEMA {
			t.Errorf("got %v want %v", gotLastLagEMA, wantLastLagEMA)
		}

		gotMacd, gotSignal := mc.GetLastMacdAndSignal()
		wantMacd, wantSignal := -2.070558190094, 3.03752586873489

		if float32(gotMacd) != float32(wantMacd) {
			t.Errorf("MACD: got %v want %v", gotMacd, wantMacd)
		}

		if float32(gotSignal) != float32(wantSignal) {
			t.Errorf("Signal: got %v want %v", gotSignal, wantSignal)
		}

		gotHistogram := mc.GetLastHistogramPoint()
		wantHistogram := -5.10808405882889
		if float32(gotHistogram) != float32(wantHistogram) {
			t.Errorf("Histogram: got %v want %v", gotHistogram, wantHistogram)
		}
	})
}
