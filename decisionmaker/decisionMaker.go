package decisionmaker

import (
	"github.com/fabiodmferreira/crypto-trading/domain"
)

// DecisionMaker decides to buy or sell
type DecisionMaker struct {
	buyStrategy  domain.Strategy
	sellStrategy domain.Strategy
}

// NewDecisionMaker returns a new instance of DecisionMaker
func NewDecisionMaker(
	buyStrategy domain.Strategy,
	sellStrategy domain.Strategy,
) *DecisionMaker {
	return &DecisionMaker{buyStrategy, sellStrategy}
}

// ShouldBuy returns true or false if it is a good time to buy
func (dm *DecisionMaker) ShouldBuy() (bool, float32, error) {
	return dm.buyStrategy.Execute()
}

// ShouldSell returns true or false if it is a good time to sell
func (dm *DecisionMaker) ShouldSell() (bool, float32, error) {
	return dm.sellStrategy.Execute()
}
