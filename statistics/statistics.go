package statistics

import (
	"math"

	"github.com/fabiodmferreira/crypto-trading/domain"
)

// Statistics receives points and do statitics calculations
type Statistics struct {
	options        domain.StatisticsOptions
	points         []float64
	average        float64
	numberOfPoints int
	variance       float64
	MACD           *MACDContainer
}

// NewStatistics returns a Statitics instance
func NewStatistics(options domain.StatisticsOptions, macd *MACDContainer) *Statistics {
	return &Statistics{options, []float64{}, 0, 0, 0, macd}
}

// AddPoint recalculate values of interest
func (s *Statistics) AddPoint(p float64) {
	s.RecalculateVariance(p)
	s.RecalculateAverage(p)
	s.MACD.AddPoint(p)

	if s.options.NumberOfPointsHold > s.numberOfPoints {
		s.numberOfPoints++
		s.points = append(s.points, p)
	} else {
		s.points = append(s.points, p)
		s.points = s.points[1:]
	}
}

// RecalculateAverage calculates new average using the old average and the new sample value
func (s *Statistics) RecalculateAverage(p float64) {
	if s.options.NumberOfPointsHold > s.numberOfPoints {
		s.average += (p - s.average) / (float64(s.numberOfPoints) + 1)
	} else {
		s.average = s.average + (1/float64(s.options.NumberOfPointsHold))*(p-s.points[0])
	}
}

// RecalculateVariance calculates new variance using the old variance and the new sample value
func (s *Statistics) RecalculateVariance(p float64) {
	if s.options.NumberOfPointsHold > s.numberOfPoints {
		// https://math.stackexchange.com/questions/198336/how-to-calculate-standard-deviation-with-streaming-inputs
		n := s.numberOfPoints + 1
		mean := s.average
		delta := p - mean
		mean += delta / float64(n)
		m := s.variance * float64(s.numberOfPoints-1)
		m += delta * (p - mean)
		if n >= 2 {
			s.variance = m / float64(n-1)
		}
	} else {
		// https://math.stackexchange.com/questions/3112650/formula-to-recalculate-variance-after-removing-a-value-and-adding-another-one-gi
		n := s.numberOfPoints
		oldP := s.points[0]
		newMean := s.average + ((p - oldP) / float64(n))
		s.variance = s.variance + math.Pow(newMean-s.average, 2) + ((math.Pow(p-newMean, 2) - math.Pow(oldP-newMean, 2)) / float64(n))
	}
}

// GetStandardDeviation returns the current standard deviation of data sample
func (s *Statistics) GetStandardDeviation() float64 {
	return math.Sqrt(s.variance)
}

// GetAverage returns the current standard deviation of data sample
func (s *Statistics) GetAverage() float64 {
	return s.average
}

// HasRequiredNumberOfPoints check whether the number of points hold are the same of required in options
func (s *Statistics) HasRequiredNumberOfPoints() bool {
	return s.numberOfPoints == s.options.NumberOfPointsHold
}
