package app

import (
	"fmt"

	"github.com/fabiodmferreira/crypto-trading/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Repository stores and gets applications from database
type Repository struct {
	repo domain.Repository
}

// NewRepository returns an instance of applications repository
func NewRepository(repo domain.Repository) *Repository {
	return &Repository{repo}
}

// FindByID returns an application with the id
func (r *Repository) FindByID(id string) (*domain.Application, error) {
	var app *domain.Application

	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, fmt.Errorf(fmt.Sprint("primitive.ObjectIDFromHex ERROR:", err))
	}

	err = r.repo.FindOne(&app, bson.M{"_id": oid}, &options.FindOneOptions{})

	return app, err
}

// Create creates an application object
func (r *Repository) Create(asset string, options domain.ApplicationOptions, accountID primitive.ObjectID) (*domain.Application, error) {

	app := domain.Application{ID: primitive.NewObjectID(), Options: options, AccountID: accountID, Asset: asset}

	err := r.repo.InsertOne(app)

	return &app, err
}

// FindAll returns all applications
func (r *Repository) FindAll() (*[]domain.Application, error) {
	var applications []domain.Application

	err := r.repo.FindAll(&applications, bson.M{}, &options.FindOptions{})

	return &applications, err
}

// DeleteByID deletes an application with the id passed by argument
func (r *Repository) DeleteByID(id string) error {
	return r.repo.DeleteByID(id)
}
