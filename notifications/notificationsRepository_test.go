package notifications_test

import (
	"testing"
	"time"

	"github.com/fabiodmferreira/crypto-trading/domain"
	"github.com/fabiodmferreira/crypto-trading/mocks"
	"github.com/fabiodmferreira/crypto-trading/notifications"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestRepositoryCreate(t *testing.T) {
	notificationsRepository, repository := setupNotificationsRepository()

	notification := &domain.Notification{
		ID:                  primitive.NewObjectID(),
		To:                  "receiver",
		Title:               "subject",
		Message:             "message",
		CreatedAt:           time.Now(),
		NotificationType:    "notificationType",
		NotificationChannel: "email",
	}

	notificationsRepository.Create(notification)

	got := len(repository.InsertOneCalls)
	want := 1

	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestRepositoryFindLastEventLogsNotificationDate(t *testing.T) {
	notificationsRepository, repository := setupNotificationsRepository()

	notificationsRepository.FindLastEventLogsNotificationDate()

	got := len(repository.FindOneCalls)
	want := 1

	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestRepositorySent(t *testing.T) {
	notificationsRepository, repository := setupNotificationsRepository()

	notificationsRepository.Sent(primitive.NewObjectID())

	got := len(repository.UpdateOneCalls)
	want := 1

	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func setupNotificationsRepository() (*notifications.Repository, *mocks.RepositorySpy) {
	repository := &mocks.RepositorySpy{}

	return notifications.NewRepository(repository), repository
}
