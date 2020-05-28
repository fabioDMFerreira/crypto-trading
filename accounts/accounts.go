package accounts

import (
	"github.com/fabiodmferreira/crypto-trading/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Account struct {
	ID     primitive.ObjectID `bson:"_id" json:"_id"`
	Amount float32            `bson:"amount,truncate" json:"amount"`
	Broker string             `json:"broker"`
}

type AccountsRepository struct {
	collection *mongo.Collection
}

func NewAccountsRepository(collection *mongo.Collection) *AccountsRepository {
	return &AccountsRepository{collection}
}

func (r *AccountsRepository) FindByBroker(broker string) (*Account, error) {
	ctx := db.NewMongoQueryContext()

	var foundDocument Account
	err := r.collection.FindOne(ctx, bson.D{{"broker", broker}}).Decode(&foundDocument)

	if err != nil {
		return nil, err
	}

	return &foundDocument, nil
}

// Create inserts a new account in collection
func (r *AccountsRepository) Create(broker string, amount float32) (*Account, error) {
	ctx := db.NewMongoQueryContext()

	account := &Account{ID: primitive.NewObjectID(), Amount: amount, Broker: broker}
	_, err := r.collection.InsertOne(ctx, account)

	return account, err
}

func (r *AccountsRepository) Withdraw(id primitive.ObjectID, amount float32) error {
	ctx := db.NewMongoQueryContext()

	filter := bson.M{"_id": id, "amount": bson.M{"$gte": amount}}
	update := bson.D{{"$inc", bson.D{{"amount", amount * -1}}}}
	_, err := r.collection.UpdateOne(ctx, filter, update)

	return err
}

func (r *AccountsRepository) Deposit(id primitive.ObjectID, amount float32) error {
	ctx := db.NewMongoQueryContext()

	filter := bson.D{{"_id", id}}
	update := bson.D{{"$inc", bson.D{{"amount", amount}}}}
	_, err := r.collection.UpdateOne(ctx, filter, update)

	return err
}

type AccountService struct {
	ID         primitive.ObjectID
	repository *AccountsRepository
}

func NewAccountService(ID primitive.ObjectID, repository *AccountsRepository) *AccountService {
	return &AccountService{ID, repository}
}

func (a *AccountService) Withdraw(amount float32) error {
	return a.repository.Withdraw(a.ID, amount)
}

func (a *AccountService) Deposit(amount float32) error {
	return a.repository.Deposit(a.ID, amount)
}
