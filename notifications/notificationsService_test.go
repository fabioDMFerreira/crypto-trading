package notifications_test

import (
	"testing"

	"github.com/fabiodmferreira/crypto-trading/domain"
	"github.com/fabiodmferreira/crypto-trading/mocks"
	"github.com/fabiodmferreira/crypto-trading/notifications"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestSendEmail(t *testing.T) {
	service, _, emailService := setupNotificationsService()

	service.SendEmail("subject", "message")

	got := len(emailService.SendMailCalls)
	want := 1

	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestCreateEmailNotification(t *testing.T) {
	service, repository, _ := setupNotificationsService()

	err := service.CreateEmailNotification("subject", "message", "type")

	if err != nil {
		t.Errorf("not expected to receive an error")
	}

	got := len(repository.SentCalls)
	want := 1
	if got != want {
		t.Errorf("expected repository.Sent to be called %v, but received %v", want, got)
	}

	got = len(repository.CreateCalls)
	want = 1
	if got != want {
		t.Errorf("expected repository.Create to be called %v, but received %v", want, got)
	}
}

func TestFindLastNotificationDate(t *testing.T) {
	service, repository, _ := setupNotificationsService()

	service.FindLastEventLogsNotificationDate()

	got := repository.FindLastEventLogsNotificationDateCalls
	want := 1

	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestShouldSendNotification(t *testing.T) {
	service, repository, _ := setupNotificationsService()

	service.ShouldSendNotification()

	got := repository.FindLastEventLogsNotificationDateCalls
	want := 1

	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func setupNotificationsService() (*notifications.Service, *mocks.NotificaitonsRepositorySpy, *mocks.EmailServiceSpy) {
	repository := &mocks.NotificaitonsRepositorySpy{}
	emailService := &mocks.EmailServiceSpy{}

	return notifications.NewService(repository, domain.NotificationOptions{Sender: "a", Receiver: "a", SenderPassword: "a"}, emailService.SendMail, primitive.NewObjectID()), repository, emailService
}
