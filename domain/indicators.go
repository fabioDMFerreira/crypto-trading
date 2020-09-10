package domain

// Indicator collects metric values and produces outputs to help decision maker deciding
type Indicator interface {
	AddValue(ohlc *OHLC)
	GetState() interface{}
}

// StatisticsOptions used in Statistics
type StatisticsOptions struct {
	NumberOfPointsHold int `json:"numberOfPointsHold"`
}

// Statistics receives points and do statitics calculations
type Statistics interface {
	AddPoint(p float64)
	GetStandardDeviation() float64
	GetVariance() float64
	GetAverage() float64
	HasRequiredNumberOfPoints() bool
}
