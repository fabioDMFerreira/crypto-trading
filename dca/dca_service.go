package dca

import (
	"errors"
	"fmt"
	"strconv"
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
	notifications domain.NotificationsService
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

// SetNotificationsService initializes service that will send notifications
func (s *Service) SetNotificationsService(service domain.NotificationsService) {
	s.notifications = service
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
			s.notify("DCA: Error", err.Error())
		}

		dcaJob.SetNextExecution()

		err = s.dcaJobsRepo.Save(dcaJob)
		if err != nil {
			s.notify("DCA: Error", "Assets bought successfully, but next time execution was not set.\n\n"+err.Error())
			return err
		}

		s.notify("DCA: success", "New crypto assets bought!")
	}

	return nil
}

// execute calls operations to buy crypto assets specified
func (s *Service) execute(dca *domain.DCAJob) error {
	coinsAmounts := dca.GetFiatCoinsAmount()

	var errorsContainer []error

	for coinSymbol, amount := range coinsAmounts {
		ticker := strings.ToUpper(coinSymbol) + "EUR"

		price, err := s.collector.GetTicker(ticker)
		if err != nil {
			errorsContainer = append(errorsContainer, fmt.Errorf("failed getting ticker \"%s\" :%s", ticker, err))
			continue
		}

		s.broker.SetTicker(coinSymbol)
		err = s.broker.AddBuyOrder(amount/price, price)
		if err != nil {
			errorsContainer = append(errorsContainer, fmt.Errorf("failed buiyng %f of %s: %s", amount/price, coinSymbol, err))
			continue
		}

		asset := &domain.DCAAsset{
			Coin:       coinSymbol,
			Amount:     amount / price,
			Price:      price,
			FiatAmount: amount,
			CreatedAt:  time.Now(),
		}
		err = s.dcaAssetsRepo.Save(asset)
		if err != nil {
			errorsContainer = append(errorsContainer, fmt.Errorf("failed saving asset: %s\n%+v", err, asset))
		}
	}

	if len(errorsContainer) > 0 {
		errorMessage := ""

		for i, err := range errorsContainer {
			errorMessage = "Error #" + strconv.Itoa(i) + ":\n" + err.Error() + "\n\n"
		}

		return errors.New(errorMessage)
	}

	return nil
}

func (s *Service) notify(subject, message string) {
	if s.notifications != nil {
		s.notifications.SendEmail(subject, message)
	}
}
