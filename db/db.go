package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	ASSETS_COLLECTION        = "assets"
	EVENT_LOGS_COLLECTION    = "eventlogs"
	NOTIFICATIONS_COLLECTION = "notifications"
)

func NewMongoQueryContext() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	return ctx
}

func ConnectDB(mongoUrl string) (*mongo.Client, error) {

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUrl))

	if err != nil {
		return nil, err
	}

	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	err = client.Ping(ctx, readpref.Primary())

	if err != nil {
		return nil, err
	}

	return client, nil
}
