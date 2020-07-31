package mocks

import (
	"time"

	"github.com/fabiodmferreira/crypto-trading/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotificaitonsRepositorySpy struct {
	CreateCalls                            []domain.Notification
	FindLastEventLogsNotificationDateCalls int
	SentCalls                              []interface{}
}

func (n *NotificaitonsRepositorySpy) Create(notification *domain.Notification) error {
	n.CreateCalls = append(n.CreateCalls, *notification)
	return nil
}

func (n *NotificaitonsRepositorySpy) FindLastEventLogsNotificationDate() (time.Time, error) {
	n.FindLastEventLogsNotificationDateCalls++
	return time.Now(), nil
}

func (n *NotificaitonsRepositorySpy) Sent(id primitive.ObjectID) error {
	n.SentCalls = append(n.SentCalls, id)
	return nil
}
