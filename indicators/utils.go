package indicators

// Average calculates the average of a set of float64 numbers
func Average(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}

	total := 0.0
	for _, p := range data {
		total += p
	}
	return total / float64(len(data))
}

// Variance calculates the variance of a set of float64 numbers
func Variance(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}

	mean := Average(data)

	total := 0.0
	for _, v := range data {
		total += (v - mean) * (v - mean)
	}

	if total == 0 {
		return total
	}

	return total / float64(len(data)-1)
}
