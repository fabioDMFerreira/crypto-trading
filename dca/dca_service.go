package dca

import (
	"fmt"
	"strings"
	"time"

	"github.com/fabiodmferreira/crypto-trading/domain"
)

// Service executes DCA operations
type Service struct {
	broker        domain.Broker
	collector     domain.Collector
	dcaJobsRepo   domain.DCAJobsRepository
	dcaAssetsRepo domain.DCAAssetsRepository
}

// NewService returns an instance of the DCAService
func NewService(broker domain.Broker, collector domain.Collector, dcaJobsRepo domain.DCAJobsRepository, dcaAssetsRepo domain.DCAAssetsRepository) *Service {
	return &Service{
		broker:        broker,
		collector:     collector,
		dcaJobsRepo:   dcaJobsRepo,
		dcaAssetsRepo: dcaAssetsRepo,
	}
}

// DrainDCA fetches dca jobs and execute dca operations
func (s *Service) DrainDCA() error {
	dcaJobs, err := s.dcaJobsRepo.FindAll()
	if err != nil {
		return err
	}

	for _, job := range *dcaJobs {
		err = s.executeDCA(&job)
		if err != nil {
			return err
		}
	}

	return nil
}

// CreateDCA creates a dca job in repository
func (s *Service) CreateDCA(dcaJob *domain.DCAJob) error {
	return s.dcaJobsRepo.Save(dcaJob)
}

// executeDCA executes dca job operation if dca execution schedule date has past
func (s *Service) executeDCA(dcaJob *domain.DCAJob) error {
	currentTime := time.Now().Unix()

	if dcaJob.NextExecution < currentTime {
		err := s.execute(dcaJob)
		if err != nil {
			return err
		}
	}

	dcaJob.SetNextExecution()

	err := s.dcaJobsRepo.Save(dcaJob)
	if err != nil {
		return err
	}

	return nil
}

// execute calls operations to buy crypto assets specified
func (s *Service) execute(dca *domain.DCAJob) error {
	coinsAmounts := dca.GetFiatCoinsAmount()

	for coinSymbol, amount := range coinsAmounts {
		price, err := s.collector.GetTicker(strings.ToUpper(coinSymbol) + "EUR")
		if err != nil {
			return err
		}

		s.broker.SetTicker(coinSymbol)
		err = s.broker.AddBuyOrder(amount/price, price)
		if err != nil {
			return fmt.Errorf("failed buyng %f of %s: %v", amount/price, coinSymbol, err)
		}

		err = s.dcaAssetsRepo.Save(&domain.DCAAsset{
			Coin:       coinSymbol,
			Amount:     amount / price,
			Price:      price,
			FiatAmount: amount,
			CreatedAt:  time.Now(),
		})
		if err != nil {
			return err
		}
	}

	return nil
}
