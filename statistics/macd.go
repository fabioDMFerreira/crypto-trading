package statistics

// MACDParams has params to calculate MACD
type MACDParams struct {
	Fast int
	Slow int
	Lag  int
}

// MACDContainer calculates and stores MACD values
type MACDContainer struct {
	params     MACDParams
	FastEMA    []float64
	SlowEMA    []float64
	MACD       []float64
	LagEMA     []float64
	Histogram  []float64
	holdPoints []float64
}

// NewMACDContainer returns an instance of MACDContainer
func NewMACDContainer(macdParams MACDParams, params ...[]float64) *MACDContainer {

	var holdPoints, fastEMA, slowEMA, macd, lagEMA, histogram []float64

	switch len(params) {
	case 1:
		holdPoints = params[0]
		break
	case 2:
		holdPoints = params[0]
		fastEMA = params[1]
		break
	case 3:
		holdPoints = params[0]
		fastEMA = params[1]
		slowEMA = params[2]
		break
	case 4:
		holdPoints = params[0]
		fastEMA = params[1]
		slowEMA = params[2]
		macd = params[3]
		break
	case 5:
		holdPoints = params[0]
		fastEMA = params[1]
		slowEMA = params[2]
		macd = params[3]
		lagEMA = params[4]
	case 6:
		holdPoints = params[0]
		fastEMA = params[1]
		slowEMA = params[2]
		macd = params[3]
		lagEMA = params[4]
		histogram = params[5]
	}

	return &MACDContainer{macdParams, fastEMA, slowEMA, macd, lagEMA, histogram, holdPoints}
}

// AddPoint users new point to calculate new MACD value
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

	if len(mc.SlowEMA) > 0 && len(mc.FastEMA) > 0 {
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

// GetLastMacdAndSignal returns the last MACD and signal values calculated
func (mc *MACDContainer) GetLastMacdAndSignal() (float64, float64) {
	if len(mc.MACD) == 0 {
		return 0, 0
	} else if len(mc.LagEMA) == 0 {
		return mc.MACD[len(mc.MACD)-1], 0
	}

	return mc.MACD[len(mc.MACD)-1], mc.LagEMA[len(mc.LagEMA)-1]
}

// GetLastHistogramPoint returns last histogram point calculated
func (mc *MACDContainer) GetLastHistogramPoint() float64 {
	if len(mc.Histogram) == 0 {
		return 0
	}

	return mc.Histogram[len(mc.Histogram)-1]
}

// EMA calculates the Exponential Moving Average
func EMA(current float64, period int, previousEMA float64) float64 {

	k := 2 / float64(period+1)
	return (current * k) + (previousEMA * (1 - k))
}
