package domain

// StatisticsOptions used in Statistics
type StatisticsOptions struct {
	NumberOfPointsHold int `json:"numberOfPointsHold"`
}

// Statistics receives points and do statitics calculations
type Statistics interface {
	AddPoint(p float64)
	GetStandardDeviation() float64
	GetAverage() float64
	HasRequiredNumberOfPoints() bool
}
