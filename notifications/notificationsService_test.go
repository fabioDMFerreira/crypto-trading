package notifications_test

import (
	"strings"
	"testing"

	"github.com/fabiodmferreira/crypto-trading/domain"
	"github.com/fabiodmferreira/crypto-trading/mocks"
	"github.com/fabiodmferreira/crypto-trading/notifications"
)

func TestSendEmail(t *testing.T) {
	service, _ := setupNotificationsService()

	got := service.SendEmail("test", "test")
	want := "Username and Password not accepted."

	if strings.Contains(got.Error(), want) != true {
		t.Errorf("%v should contain %v", got, want)
	}
}

func TestCreateEmailNotification(t *testing.T) {
	service, repository := setupNotificationsService()

	err := service.CreateEmailNotification("subject", "message", "type")

	if err == nil {
		t.Errorf("expected to receive error")
	}

	got := len(repository.SentCalls)
	want := 0
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
	service, repository := setupNotificationsService()

	service.FindLastEventLogsNotificationDate()

	got := repository.FindLastEventLogsNotificationDateCalls
	want := 1

	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestShouldSendNotification(t *testing.T) {
	service, repository := setupNotificationsService()

	service.ShouldSendNotification()

	got := repository.FindLastEventLogsNotificationDateCalls
	want := 1

	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func setupNotificationsService() (*notifications.Service, *mocks.NotificaitonsRepositorySpy) {
	repository := &mocks.NotificaitonsRepositorySpy{}

	return notifications.NewService(repository, domain.NotificationOptions{Sender: "a", Receiver: "a", SenderPassword: "a"}), repository
}
