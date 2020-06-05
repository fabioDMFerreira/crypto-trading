package statistics

type StatisticsOptions struct {
	NumberOfPointsHold int
}

type Statistics struct {
	options        StatisticsOptions
	Points         []float64
	Average        float64
	NumberOfPoints int
	MACD           *MACDContainer
}

func NewStatistics(options StatisticsOptions, macd *MACDContainer) *Statistics {
	return &Statistics{options, []float64{}, 0, 0, macd}
}

func (s *Statistics) AddPoint(p float64) {
	s.RecalculateAverage(p)
	s.MACD.AddPoint(p)

	if s.options.NumberOfPointsHold > s.NumberOfPoints {
		s.NumberOfPoints++
		s.Points = append(s.Points, p)
	} else {
		s.Points = append(s.Points, p)
		s.Points = s.Points[1:]
	}
}

func (s *Statistics) RecalculateAverage(p float64) {
	if s.options.NumberOfPointsHold > s.NumberOfPoints {
		s.Average += (p - s.Average) / (float64(s.NumberOfPoints) + 1)
	} else {
		s.Average = s.Average + (1/float64(s.options.NumberOfPointsHold))*(p-s.Points[0])
	}
}

func Average(xs []float64) float64 {
	total := 0.0
	for _, v := range xs {
		total += v
	}
	return total / float64(len(xs))
}

type MACDParams struct {
	Fast int
	Slow int
	Lag  int
}

type MACDContainer struct {
	params     MACDParams
	FastEMA    []float64
	SlowEMA    []float64
	MACD       []float64
	LagEMA     []float64
	Histogram  []float64
	holdPoints []float64
}

func NewMACDContainer(macdParams MACDParams, holdPoints []float64, params ...[]float64) *MACDContainer {

	var fastEMA, slowEMA, macd, lagEMA, histogram []float64

	switch len(params) {
	case 1:
		fastEMA = params[0]
		break
	case 2:
		fastEMA = params[0]
		slowEMA = params[1]
		break
	case 3:
		fastEMA = params[0]
		slowEMA = params[1]
		macd = params[2]
		break
	case 4:
		fastEMA = params[0]
		slowEMA = params[1]
		macd = params[2]
		lagEMA = params[3]
	case 5:
		fastEMA = params[0]
		slowEMA = params[1]
		macd = params[2]
		lagEMA = params[3]
		histogram = params[4]
	}

	return &MACDContainer{macdParams, fastEMA, slowEMA, macd, lagEMA, histogram, holdPoints}
}

func (mc *MACDContainer) AddPoint(p float64) {
	if len(mc.holdPoints) <= mc.params.Slow {
		mc.holdPoints = append(mc.holdPoints, p)
	}

	if len(mc.holdPoints) > mc.params.Slow {
		slow := EMA(p, mc.params.Slow, mc.SlowEMA[len(mc.SlowEMA)-1])
		mc.SlowEMA = append(mc.SlowEMA, slow)
	} else if len(mc.holdPoints) == mc.params.Slow {
		slow := Average(mc.holdPoints)
		mc.SlowEMA = []float64{slow}
	}

	if len(mc.holdPoints) > mc.params.Fast {
		fast := EMA(p, mc.params.Fast, mc.FastEMA[len(mc.FastEMA)-1])
		mc.FastEMA = append(mc.FastEMA, fast)
	} else if len(mc.holdPoints) == mc.params.Fast {
		fast := Average(mc.holdPoints)
		mc.FastEMA = []float64{fast}
	}

	if len(mc.SlowEMA) > 0 {
		macd := mc.FastEMA[len(mc.FastEMA)-1] - mc.SlowEMA[len(mc.SlowEMA)-1]
		mc.MACD = append(mc.MACD, macd)
	}

	if len(mc.MACD) > mc.params.Lag {
		lag := EMA(mc.MACD[len(mc.MACD)-1], mc.params.Lag, mc.LagEMA[len(mc.LagEMA)-1])
		mc.LagEMA = append(mc.LagEMA, lag)
		histogram := mc.MACD[len(mc.MACD)-1] - lag
		mc.Histogram = []float64{histogram}
	} else if len(mc.MACD) == mc.params.Lag {
		lag := Average(mc.MACD)
		mc.LagEMA = append(mc.LagEMA, lag)
		histogram := mc.MACD[len(mc.MACD)-1] - lag
		mc.Histogram = []float64{histogram}
	}
}

func (mc *MACDContainer) GetLastMacdAndSignal() (float64, float64) {
	if len(mc.MACD) == 0 {
		return 0, 0
	} else if len(mc.LagEMA) == 0 {
		return mc.MACD[len(mc.MACD)-1], 0
	}

	return mc.MACD[len(mc.MACD)-1], mc.LagEMA[len(mc.LagEMA)-1]
}

func (mc *MACDContainer) GetLastHistogramPoint() float64 {
	if len(mc.Histogram) == 0 {
		return 0
	}

	return mc.Histogram[len(mc.Histogram)-1]
}

func EMA(current float64, period int, previousEMA float64) float64 {

	k := 2 / float64(period+1)
	return (current * k) + (previousEMA * (1 - k))
}
