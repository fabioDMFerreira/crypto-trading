package accounts

import (
	"github.com/fabiodmferreira/crypto-trading/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Account domain.Account

type Repository struct {
	repo domain.Repository
}

func NewRepository(repo domain.Repository) *Repository {
	return &Repository{repo}
}

func (r *Repository) FindById(id primitive.ObjectID) (*Account, error) {
	var foundDocument Account

	err := r.repo.FindOne(&foundDocument, bson.M{"_id": id}, options.FindOne())

	if err != nil {
		return nil, err
	}

	return &foundDocument, nil
}

func (r *Repository) FindByBroker(broker string) (*Account, error) {
	var foundDocument Account

	err := r.repo.FindOne(&foundDocument, bson.M{"broker": broker}, options.FindOne())

	if err != nil {
		return nil, err
	}
	return &foundDocument, nil
}

// Create inserts a new account in collection
func (r *Repository) Create(broker string, amount float32) (*Account, error) {

	account := &Account{ID: primitive.NewObjectID(), Amount: amount, Broker: broker}
	err := r.repo.InsertOne(account)

	return account, err
}

func (r *Repository) Withdraw(id primitive.ObjectID, amount float32) error {

	filter := bson.M{"_id": id, "amount": bson.M{"$gte": amount}}
	update := bson.M{"$inc": bson.M{"amount": amount * -1}}
	return r.repo.UpdateOne(filter, update)
}

func (r *Repository) Deposit(id primitive.ObjectID, amount float32) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$inc": bson.M{"amount": amount}}

	return r.repo.UpdateOne(filter, update)
}
