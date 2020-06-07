package statistics

// SubSlices subtracts two slices.
func SubSlices(slice1, slice2 []float64) []float64 {

	var result []float64

	for i := 0; i < len(slice1); i++ {
		result = append(result, slice1[i]-slice2[i])
	}

	return result
}

// TransverseEMA calculates exponential moving average of a slice for a certain period
func TransverseEMA(slice []float64, period int) []float64 {

	var emaSlice []float64

	ak := period + 1
	k := float64(2) / float64(ak)

	emaSlice = append(emaSlice, slice[0])

	for i := 1; i < len(slice); i++ {
		emaSlice = append(emaSlice, (slice[i]*float64(k))+(emaSlice[i-1]*float64(1-k)))
	}

	return emaSlice
}

// TransverseMACD stands for moving average convergence divergence.
func TransverseMACD(points []float64, params ...int) ([]float64, []float64) {

	if len(params) < 3 {
		params = []int{12, 26, 9}
	}

	ema12 := TransverseEMA(points, params[0])
	ema26 := TransverseEMA(points, params[1])
	macd := SubSlices(ema12, ema26)
	ema9 := TransverseEMA(macd, params[2])

	return macd, ema9
}
