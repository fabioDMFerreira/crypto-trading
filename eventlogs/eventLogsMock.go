package eventlogs

import (
	"github.com/fabiodmferreira/crypto-trading/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// EventLogsServiceMock mocks EventLogsService
type EventLogsServiceMock struct {
	logs [][]string
}

// Create saves log event internally
func (l *EventLogsServiceMock) Create(logType, message string) error {
	l.logs = append(l.logs, []string{logType, message})
	return nil
}

// FindAllToNotify returns an empty slice of log events
func (l *EventLogsServiceMock) FindAllToNotify() (*[]domain.EventLog, error) {
	return &[]domain.EventLog{}, nil
}

// MarkNotified mock the update of log events
func (l *EventLogsServiceMock) MarkNotified(ids []primitive.ObjectID) error {
	return nil
}
