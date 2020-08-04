package accounts

import (
	"github.com/fabiodmferreira/crypto-trading/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Repository stores and gets accounts from database
type Repository struct {
	repo domain.Repository
}

// NewRepository returns an instance of accounts repository
func NewRepository(repo domain.Repository) *Repository {
	return &Repository{repo}
}

// FindById returns an account with the id passed by argument
func (r *Repository) FindById(id primitive.ObjectID) (*domain.Account, error) {
	var foundDocument domain.Account

	err := r.repo.FindOne(&foundDocument, bson.M{"_id": id}, options.FindOne())

	if err != nil {
		return nil, err
	}

	return &foundDocument, nil
}

// FindByBroker returns an account with the broker passed by argument
func (r *Repository) FindByBroker(broker string) (*domain.Account, error) {
	var foundDocument domain.Account

	err := r.repo.FindOne(&foundDocument, bson.M{"broker": broker}, options.FindOne())

	if err != nil {
		return nil, err
	}
	return &foundDocument, nil
}

// Create inserts a new account in collection
func (r *Repository) Create(broker string, amount float32) (*domain.Account, error) {

	account := &domain.Account{ID: primitive.NewObjectID(), Amount: amount, Broker: broker}
	err := r.repo.InsertOne(account)

	return account, err
}

// Withdraw decrements an amount from the account
func (r *Repository) Withdraw(id primitive.ObjectID, amount float32) error {

	filter := bson.M{"_id": id, "amount": bson.M{"$gte": amount}}
	update := bson.M{"$inc": bson.M{"amount": amount * -1}}
	return r.repo.UpdateOne(filter, update)
}

// Deposit increments an amount to the account
func (r *Repository) Deposit(id primitive.ObjectID, amount float32) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$inc": bson.M{"amount": amount}}

	return r.repo.UpdateOne(filter, update)
}
