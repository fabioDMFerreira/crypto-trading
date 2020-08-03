package app

import (
	"fmt"

	"github.com/fabiodmferreira/crypto-trading/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	repo domain.Repository
}

func NewRepository(repo domain.Repository) *Repository {
	return &Repository{repo}
}

func (r *Repository) FindByID(id string) (*domain.Application, error) {
	var app *domain.Application

	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, fmt.Errorf(fmt.Sprint("primitive.ObjectIDFromHex ERROR:", err))
	}

	err = r.repo.FindOne(&app, bson.M{"_id": oid}, &options.FindOneOptions{})

	return app, err
}

func (r *Repository) Create(asset string, options domain.ApplicationOptions, accountID primitive.ObjectID) (*domain.Application, error) {

	app := domain.Application{ID: primitive.NewObjectID(), Options: options, AccountID: accountID, Asset: asset}

	err := r.repo.InsertOne(app)

	return &app, err
}
