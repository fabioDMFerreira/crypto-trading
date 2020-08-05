package app

import (
	"time"

	"github.com/fabiodmferreira/crypto-trading/domain"
	"github.com/fabiodmferreira/crypto-trading/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Service interacts with applications, log events and application state repositories
type Service struct {
	repo          *Repository
	stateRepo     domain.ApplicationExecutionStateRepository
	logEventsRepo domain.EventsLog
}

// NewService returns an instance of applications service
func NewService(repo *Repository, stateRepo domain.ApplicationExecutionStateRepository, eventsLogRepo domain.EventsLog) *Service {
	return &Service{repo, stateRepo, eventsLogRepo}
}

// GetLastState returns the last application state
func (a *Service) GetLastState(appID primitive.ObjectID) (*domain.ApplicationExecutionState, error) {
	return a.stateRepo.FindLast(bson.M{"executionId": appID})
}

// FindAll returns all applications in the repository
func (a *Service) FindAll() (*[]domain.Application, error) {
	return a.repo.FindAll()
}

// DeleteByID deletes the application with the id passed by argument
func (a *Service) DeleteByID(id string) error {
	return a.repo.repo.DeleteByID(id)
}

// GetLogEvents returns all log events generated by a application
func (a *Service) GetLogEvents(appID primitive.ObjectID) (*[]domain.EventLog, error) {
	return a.logEventsRepo.FindAll(bson.M{"applicationID": appID})
}

func (a *Service) GetStateAggregated(appID string, startDate, endDate time.Time) (*[]bson.M, error) {
	groupByDatesClause := utils.GetGroupByDatesIDClause(startDate, endDate)

	oid, err := primitive.ObjectIDFromHex(appID)

	if err != nil {
		return nil, err
	}

	pipelineOptions := mongo.Pipeline{
		{
			primitive.E{
				Key: "$match",
				Value: bson.M{
					"executionId": oid,
					"date":        bson.M{"$gte": startDate, "$lte": endDate}},
			},
		},
		{
			primitive.E{
				Key: "$group",
				Value: bson.M{
					"_id":                 groupByDatesClause,
					"average":             bson.M{"$avg": "$state.average"},
					"standardDeviation":   bson.M{"$avg": "$state.standardDeviation"},
					"higherBollingerBand": bson.M{"$avg": "$state.higherBollingerBand"},
					"lowerBollingerBand":  bson.M{"$avg": "$state.lowerBollingerBand"},
					"currentChange":       bson.M{"$avg": "$state.currentChange"},
					"accountAmount":       bson.M{"$avg": "$state.accountAmount"},
				},
			}},
	}
	return a.stateRepo.Aggregate(pipelineOptions)
}
