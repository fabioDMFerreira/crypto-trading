package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

// ApplicationOptions aggregates options of every service options needed to run an application
type ApplicationOptions struct {
	NotificationOptions  `bson:"notificationOptions" json:"notificationOptions"`
	StatisticsOptions    `bson:"statisticsOptions" json:"statisticsOptions"`
	DecisionMakerOptions `bson:"decisionMakerOptions" json:"decisionMakerOptions"`
	CollectorOptions     `bson:"collectorOptions" json:"collectorOptions"`
}

// Application stores all options and required relations ids for a running application
type Application struct {
	ID        primitive.ObjectID `bson:"_id" json:"_id"`
	Asset     string             `json:"asset"`
	AccountID primitive.ObjectID `bson:"accountID" json:"accountID"`
	Options   ApplicationOptions `bson:"options" json:"options"`
}

// ApplicationRepository stores and gets applications from db
type ApplicationRepository interface {
	FindByID(id string) (*Application, error)
	Create(asset string, options ApplicationOptions, acountID primitive.ObjectID) (*Application, error)
}
