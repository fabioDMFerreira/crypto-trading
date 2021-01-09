package domain

import (
	"time"

	"github.com/fabiodmferreira/crypto-trading/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DCAJob stores details about a dca job like the time it has to be executed and strategy options
type DCAJob struct {
	ID            primitive.ObjectID `bson:"_id" json:"_id"`
	Period        int64
	NextExecution int64
	Options       DCAJobOptions
}

// SetNextExecution updates dca job with the next timestamp to execute the job
func (d *DCAJob) SetNextExecution() {
	d.NextExecution += d.Period
}

// GetFiatCoinsAmount use coins values proportion and total FIAT amount that must be invested and returns amounts to invest in each coin
func (d *DCAJob) GetFiatCoinsAmount() map[string]float32 {
	coinsToBuy := map[string]float32{}
	var totalCoinsValue float32

	for _, coinProportion := range d.Options.CoinsProportion {
		totalCoinsValue += coinProportion
	}

	for coinSymbol, coinProportion := range d.Options.CoinsProportion {
		fiatAmount := float64((d.Options.TotalFIATAmount * coinProportion) / totalCoinsValue)
		coinsToBuy[coinSymbol] = utils.RoundFloorDecimals(float32(fiatAmount), 2)
	}

	return coinsToBuy
}

// DCAJobOptions stores dca jobs strategy options
type DCAJobOptions struct {
	CoinsProportion map[string]float32
	TotalFIATAmount float32
}

// DCAJobsRepository stores and gets dca jobs
type DCAJobsRepository interface {
	Save(dcaJob *DCAJob) error
	FindAll() (*[]DCAJob, error)
}

// DCAAsset represents an asset bought on a dca operation
type DCAAsset struct {
	Coin       string
	Amount     float32
	Price      float32
	FiatAmount float32
	CreatedAt  time.Time `bson:"createdAt" json:"createdAt"`
}

// DCAAssetsRepository stores and gets dca assets
type DCAAssetsRepository interface {
	Save(dcaJob *DCAAsset) error
	FindAll() (*[]DCAAsset, error)
}
