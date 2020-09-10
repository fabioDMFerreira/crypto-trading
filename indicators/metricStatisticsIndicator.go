package indicators

import "github.com/fabiodmferreira/crypto-trading/domain"

type MetricStatisticsIndicatorState struct {
	HasRequiredPoints bool

	Value             float64
	Average           float64
	StandardDeviation float64

	Change                  float64
	ChangeAverage           float64
	ChangeStandardDeviation float64

	Acceleration                  float64
	AccelerationAverage           float64
	AccelerationStandardDeviation float64
}

type MetricStatisticsIndicator struct {
	stats             *Statistics
	changeStats       *Statistics
	accelerationStats *Statistics

	currentValue  float32
	previousValue float32

	currentChange  float32
	previousChange float32
}

func NewMetricStatisticsIndicator(options domain.StatisticsOptions) *MetricStatisticsIndicator {
	return &MetricStatisticsIndicator{
		stats:             NewStatistics(options),
		changeStats:       NewStatistics(options),
		accelerationStats: NewStatistics(options),
	}
}

// AddValue appends close price to statistics
func (m *MetricStatisticsIndicator) AddMetricValue(value float32) {
	if m.previousValue > 0 {
		change := value - m.previousValue

		m.previousChange = m.currentChange
		m.currentChange = change

		m.changeStats.AddPoint(float64(change))
	}

	if m.previousChange != 0 {
		acceleration := m.currentChange - m.previousChange

		m.accelerationStats.AddPoint(float64(acceleration))
	}

	m.previousValue = m.currentValue
	m.currentValue = value

	m.stats.AddPoint(float64(value))
}

// GetState returns the state value
func (m *MetricStatisticsIndicator) GetState() interface{} {
	return &MetricStatisticsIndicatorState{
		HasRequiredPoints: m.stats.HasRequiredNumberOfPoints(),

		Value:             m.stats.GetLastValue(),
		Average:           m.stats.GetAverage(),
		StandardDeviation: m.stats.GetStandardDeviation(),

		Change:                  m.changeStats.GetLastValue(),
		ChangeAverage:           m.changeStats.GetAverage(),
		ChangeStandardDeviation: m.changeStats.GetStandardDeviation(),

		Acceleration:                  m.accelerationStats.GetLastValue(),
		AccelerationAverage:           m.accelerationStats.GetAverage(),
		AccelerationStandardDeviation: m.accelerationStats.GetStandardDeviation(),
	}
}

type PriceIndicator struct {
	*MetricStatisticsIndicator
}

func NewPriceIndicator(metricIndicator *MetricStatisticsIndicator) *PriceIndicator {
	return &PriceIndicator{metricIndicator}
}

func (pi *PriceIndicator) AddValue(ohlc *domain.OHLC) {
	pi.AddMetricValue((ohlc.Close + ohlc.High + ohlc.Low + ohlc.Open) / 4)
}

type VolumeIndicator struct {
	*MetricStatisticsIndicator
}

func NewVolumeIndicator(metricIndicator *MetricStatisticsIndicator) *VolumeIndicator {
	return &VolumeIndicator{metricIndicator}
}

func (vi *VolumeIndicator) AddValue(ohlc *domain.OHLC) {
	vi.AddMetricValue(ohlc.Volume)
}
