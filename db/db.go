package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	ASSETS_COLLECTION                       = "assets"
	EVENT_LOGS_COLLECTION                   = "eventlogs"
	NOTIFICATIONS_COLLECTION                = "notifications"
	ACCOUNTS_COLLECTION                     = "accounts"
	BENCHMARKS_COLLECTION                   = "benchmarks"
	ASSETS_PRICES_COLLECTION                = "assetsprices"
	APPLICATION_EXECUTION_STATES_COLLECTION = "applicationExecutionStates"
	APPLICATIONS_COLLECTION                 = "applications"
	DCA_JOBS_COLLECTION                     = "dcaJobs"
	DCA_ASSETS_COLLECTION                   = "dcaAssets"
)

func NewMongoQueryContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)

	return ctx, cancel
}

func ConnectDB(mongoUrl string) (*mongo.Client, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUrl))

	if err != nil {
		cancel()
		return nil, err
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	err = client.Ping(ctx, readpref.Primary())

	defer cancel()

	if err != nil {
		return nil, err
	}

	return client, nil
}
