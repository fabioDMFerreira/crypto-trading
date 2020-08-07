package notifications_test

import (
	"net/smtp"
	"testing"

	"github.com/fabiodmferreira/crypto-trading/domain"
	"github.com/fabiodmferreira/crypto-trading/mocks"
	"github.com/fabiodmferreira/crypto-trading/notifications"
	"github.com/golang/mock/gomock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestSendEmail(t *testing.T) {
	service, _, emailService := setupNotificationsService(t)

	emailService.EXPECT().SendEmail("subject", "message").Times(1)

	service.SendEmail("subject", "message")
}

func TestCreateEmailNotification(t *testing.T) {

	service, repository, _ := setupNotificationsService(t)

	repository.EXPECT().Sent(gomock.Any()).Times(1)

	repository.EXPECT().Create(gomock.Any()).Times(1)

	err := service.CreateEmailNotification("subject", "message", "type")

	if err != nil {
		t.Errorf("not expected to receive an error")
	}
}

func TestFindLastNotificationDate(t *testing.T) {
	service, repository, _ := setupNotificationsService(t)

	repository.EXPECT().FindLastEventLogsNotificationDate().Times(1)

	service.FindLastEventLogsNotificationDate()
}

func TestShouldSendNotification(t *testing.T) {
	service, repository, _ := setupNotificationsService(t)

	repository.EXPECT().FindLastEventLogsNotificationDate().Times(1)

	service.ShouldSendNotification()
}

func setupNotificationsService(t *testing.T) (*notifications.Service, *mocks.MockNotificationsRepository, *mocks.MockNotificationsService) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	repository := mocks.NewMockNotificationsRepository(ctrl)
	emailService := mocks.NewMockNotificationsService(ctrl)

	return notifications.NewService(
		repository,
		domain.NotificationOptions{Sender: "a", Receiver: "a", SenderPassword: "a"},
		func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
			return nil
		}, primitive.NewObjectID()), repository, emailService
}
