package appkeeper

import (
	"context"
	"fmt"

	krakenapi "github.com/beldur/kraken-go-api-client"
	"github.com/fabiodmferreira/crypto-trading/app"
	"github.com/fabiodmferreira/crypto-trading/appfactory"
	"github.com/fabiodmferreira/crypto-trading/collectors"
	"github.com/fabiodmferreira/crypto-trading/db"
	"github.com/fabiodmferreira/crypto-trading/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoQuery struct {
	ID primitive.ObjectID `bson:"_id"  json:"_id"`
}

// AppKeeper manages algorithm applications state by starting and stopping them
type AppKeeper struct {
	applications     map[string]*app.App
	mongoDatabase    *mongo.Database
	krakenAPI        *krakenapi.KrakenAPI
	appEnv           string
	applicationsRepo domain.ApplicationRepository
}

// NewAppKeeper returns an instance of AppKeeper
func NewAppKeeper(db *mongo.Database, krakenAPI *krakenapi.KrakenAPI, applicationRepo domain.ApplicationRepository) *AppKeeper {
	return &AppKeeper{
		applications:     map[string]*app.App{},
		mongoDatabase:    db,
		krakenAPI:        krakenAPI,
		applicationsRepo: applicationRepo,
	}
}

// SetAppEnv sets the application environment variable that decides whether the broker API should be mocked
func (ak *AppKeeper) SetAppEnv(appEnv string) {
	ak.appEnv = appEnv
}

// Initialize starts events listener that helps AppKeeper to keep applications state consistent with database
func (ak *AppKeeper) Initialize() {
	applicationsCollection := ak.mongoDatabase.Collection(db.APPLICATIONS_COLLECTION)

	feeder := make(chan bson.M)

	go listenMongoCollectionChanges(applicationsCollection, feeder)

	for {
		change := <-feeder
		ak.applyCollectionChanges(change)
	}
}

// StartApplications starts applications with metadata passed by argument
func (ak *AppKeeper) StartApplications(applications *[]domain.Application) error {
	for _, metadata := range *applications {
		err := ak.StartApplication(&metadata)

		if err != nil {
			return err
		}
	}

	return nil
}

// StartApplication starts application with metadata passed by argument
func (ak *AppKeeper) StartApplication(metadata *domain.Application) error {
	brokerService := appfactory.GetBroker(ak.appEnv, ak.krakenAPI)

	collector := collectors.NewKrakenCollector(metadata.Asset, domain.CollectorOptions{NewPriceTimeRate: 1}, ak.krakenAPI, &[]domain.Indicator{})

	application, err := appfactory.SetupApplication(metadata, ak.mongoDatabase, brokerService, collector)

	if err != nil {
		return err
	}

	go application.Start()

	ak.applications[metadata.ID.Hex()] = application

	return nil
}

func (ak *AppKeeper) restartApp(metadata *domain.Application) error {
	if application, ok := ak.applications[metadata.ID.Hex()]; ok {
		application.Stop()
	}

	err := ak.StartApplication(metadata)
	if err != nil {
		return fmt.Errorf("Not able to start application with ID %v due to next error: %v", metadata.ID, err)
	}

	return nil
}

func (ak *AppKeeper) applyCollectionChanges(change bson.M) {
	switch operationType := change["operationType"]; operationType {
	case "delete":
		metadata := parseBsonToAppMetadata(change["fullDocument"])
		if application, ok := ak.applications[metadata.ID.Hex()]; ok {
			application.Stop()
			ak.applications[metadata.ID.Hex()] = nil
		}
	case "update":
		var query mongoQuery
		bsonBytes, _ := bson.Marshal(change["documentKey"])
		bson.Unmarshal(bsonBytes, &query)
		if query.ID.Hex() != "000000000000000000000000" {
			metadata, err := ak.applicationsRepo.FindByID(query.ID.Hex())
			if err != nil {
				fmt.Printf("Not able to start application with ID %v due to next error: %v", metadata.ID, err)
			}
			ak.restartApp(metadata)
		} else {
			fmt.Printf("Query should have an _id (%+v).", change["documentKey"])
		}
	case "replace", "insert":
		metadata := parseBsonToAppMetadata(change["fullDocument"])
		ak.restartApp(metadata)
	default:
		fmt.Printf("new event => %v", operationType)
	}
}

func parseBsonToAppMetadata(data interface{}) *domain.Application {
	var metadata domain.Application
	bsonBytes, _ := bson.Marshal(data)
	bson.Unmarshal(bsonBytes, &metadata)
	return &metadata
}

func listenMongoCollectionChanges(collection *mongo.Collection, ch chan bson.M) {
	applicationsStream, err := collection.Watch(context.TODO(), mongo.Pipeline{})
	if err != nil {
		panic(err)
	}

	defer applicationsStream.Close(context.TODO())

	for applicationsStream.Next(context.TODO()) {
		var data bson.M
		if err := applicationsStream.Decode(&data); err != nil {
			panic(err)
		}

		ch <- data
	}
}
