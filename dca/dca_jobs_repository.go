package dca

import (
	"github.com/fabiodmferreira/crypto-trading/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// JobsRepository fetches and stores dca jobs
type JobsRepository struct {
	repo domain.Repository
}

// NewJobsRepository returns an instance of dca Jobsrepository
func NewJobsRepository(repo domain.Repository) *JobsRepository {
	return &JobsRepository{repo}
}

// Save persists dca job in database
func (r *JobsRepository) Save(dcaJob *domain.DCAJob) error {
	if dcaJob.ID.IsZero() {
		dcaJob.ID = primitive.NewObjectID()
		return r.repo.InsertOne(dcaJob)
	}

	return r.repo.UpdateOne(bson.M{"_id": dcaJob.ID}, bson.M{"$set": bson.M{"nextexecution": dcaJob.NextExecution}})
}

// FindAll fetches existing dca jobs
func (r *JobsRepository) FindAll() (*[]domain.DCAJob, error) {
	var results []domain.DCAJob
	err := r.repo.FindAll(&results, bson.M{}, nil)

	if err != nil {
		return nil, err
	}

	return &results, nil
}
