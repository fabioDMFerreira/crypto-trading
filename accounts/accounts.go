package accounts

import (
	"github.com/fabiodmferreira/crypto-trading/db"
	"github.com/fabiodmferreira/crypto-trading/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Account domain.Account

type Repository struct {
	collection *mongo.Collection
}

func NewRepository(collection *mongo.Collection) *Repository {
	return &Repository{collection}
}

func (r *Repository) FindById(id primitive.ObjectID) (*Account, error) {
	ctx := db.NewMongoQueryContext()

	var foundDocument Account
	err := r.collection.FindOne(ctx, bson.D{{"_id", id}}).Decode(&foundDocument)

	if err != nil {
		return nil, err
	}

	return &foundDocument, nil
}

func (r *Repository) FindByBroker(broker string) (*Account, error) {
	ctx := db.NewMongoQueryContext()

	var foundDocument Account
	err := r.collection.FindOne(ctx, bson.D{{"broker", broker}}).Decode(&foundDocument)

	if err != nil {
		return nil, err
	}

	return &foundDocument, nil
}

// Create inserts a new account in collection
func (r *Repository) Create(broker string, amount float32) (*Account, error) {
	ctx := db.NewMongoQueryContext()

	account := &Account{ID: primitive.NewObjectID(), Amount: amount, Broker: broker}
	_, err := r.collection.InsertOne(ctx, account)

	return account, err
}

func (r *Repository) Withdraw(id primitive.ObjectID, amount float32) error {
	ctx := db.NewMongoQueryContext()

	filter := bson.M{"_id": id, "amount": bson.M{"$gte": amount}}
	update := bson.D{{"$inc", bson.D{{"amount", amount * -1}}}}
	_, err := r.collection.UpdateOne(ctx, filter, update)

	return err
}

func (r *Repository) Deposit(id primitive.ObjectID, amount float32) error {
	ctx := db.NewMongoQueryContext()

	filter := bson.D{{"_id", id}}
	update := bson.D{{"$inc", bson.D{{"amount", amount}}}}
	_, err := r.collection.UpdateOne(ctx, filter, update)

	return err
}
