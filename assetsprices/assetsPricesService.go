package assetsprices

import (
	"strings"
	"time"

	"github.com/fabiodmferreira/crypto-trading/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type FetchRemotePricesType = func(startDate, endDate time.Time, asset string) (*[]bson.M, error)

// Service provides assets prices methods
type Service struct {
	repo              domain.AssetPriceRepository
	FetchRemotePrices FetchRemotePricesType
}

// NewService returns an instance of AssetsPricesService
func NewService(repo domain.AssetPriceRepository, fetchRemotePrices FetchRemotePricesType) *Service {
	return &Service{repo, fetchRemotePrices}
}

// FetchAndStoreAssetPrices fetches asset prices remotely and save it in repository
func (s *Service) FetchAndStoreAssetPrices(asset string, endDate time.Time) error {
	lastAssetPrice, err := s.GetLastAssetsPrices(asset, 1)

	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}

	var startDate time.Time

	if len(*lastAssetPrice) == 0 {
		startDate = endDate.AddDate(0, 0, -180)
	} else {
		startDate = (*lastAssetPrice)[0].Date
	}

	counter := 0

	TransverseDatesRange(startDate, endDate, func(startDate, endDate time.Time) error {
		assetsPrices, err := s.FetchRemotePrices(startDate, endDate, asset)

		if err != nil {
			return err
		}

		// Tech Debt: Condition to avoid create the same asset price multiple time.
		// This happens when `FetchAndStoreAssetPrices` is executed multiple times in a short time
		if len(*assetsPrices) > 1 {
			err = s.repo.BulkCreate(assetsPrices)

			counter += len(*assetsPrices)
			// fmt.Printf("\rAsset: %v Created: %d", asset, counter)
		}

		return nil
	})

	return nil
}

// Create creates an asset price in repository
func (s *Service) Create(date time.Time, value float32, asset string) error {
	return s.repo.Create(date, value, asset)
}

// GetLastAssetsPrices returns the last price of an asset in the repository
func (s *Service) GetLastAssetsPrices(asset string, limit int) (*[]domain.AssetPrice, error) {
	return s.repo.GetLastAssetsPrices(asset, limit)
}

// TransverseDatesRange iterate every day in the dates range passed and call the callback function
func TransverseDatesRange(startDate, endDate time.Time, handle func(time.Time, time.Time) error) {
	startDateCursor := startDate
	endDateCursor := startDateCursor.Add(
		23*time.Hour +
			59*time.Minute +
			59*time.Second)

	for startDateCursor.Before(endDate) {
		handle(startDateCursor, endDateCursor)
		startDateCursor, endDateCursor = GetDatesPlusOneDay(startDateCursor, endDateCursor)
	}
}

// GetDatesPlusOneDay returns the two dates passed by parameter with one more day
func GetDatesPlusOneDay(startDate, endDate time.Time) (time.Time, time.Time) {
	return startDate.AddDate(0, 0, 1), endDate.AddDate(0, 0, 1)
}

// SerializeDate convert time to string formatted to be accepted on coindesk url as parameter
func SerializeDate(date time.Time) string {
	return strings.Replace(date.Format("2006-01-02 15:04"), " ", "T", 1)
}
