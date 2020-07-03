package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

// Account has details about an exchange account
type Account struct {
	ID     primitive.ObjectID `bson:"_id" json:"_id"`
	Amount float32            `bson:"amount,truncate" json:"amount"`
	Broker string             `json:"broker"`
}

// AccountsRepository stores and fetches accounts
type AccountsRepository interface {
	FindById(id primitive.ObjectID) (*Account, error)
	FindByBroker(broker string) (*Account, error)
	Create(broker string, amount float32) (*Account, error)
	Withdraw(id primitive.ObjectID, amount float32) error
	Deposit(id primitive.ObjectID, amount float32) error
}

// AccountServiceReader reads information about one account
type AccountServiceReader interface {
	GetAmount() (float32, error)
}

// AccountService interacts with one account
type AccountService interface {
	AccountServiceReader
	Withdraw(amount float32) error
	Deposit(amount float32) error
}
