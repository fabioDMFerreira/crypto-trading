package notifications

import (
	"net/smtp"
	"time"

	"github.com/fabiodmferreira/crypto-trading/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	notificationsRepository *Repository
	options                 domain.NotificationOptions
}

func NewService(
	notificationsRepository *Repository,
	options domain.NotificationOptions,
) *Service {
	return &Service{notificationsRepository, options}
}

func (n *Service) SendEmail(subject, body string) error {
	from := n.options.Sender
	pass := n.options.SenderPassword
	to := n.options.Receiver

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n" +
		"MIME-Version: 1.0;\n" +
		"Content-Type: text/html;\n\n" +

		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	return err
}

func (n *Service) FindLastEventLogsNotificationDate() (time.Time, error) {
	return n.notificationsRepository.FindLastEventLogsNotificationDate()
}

func (n *Service) CreateEmailNotification(subject, message, notificationType string) error {
	notification := &Notification{
		ID:                  primitive.NewObjectID(),
		To:                  n.options.Receiver,
		Title:               subject,
		Message:             message,
		CreatedAt:           time.Now(),
		NotificationType:    notificationType,
		NotificationChannel: "email",
	}

	err := n.notificationsRepository.Create(notification)

	if err != nil {
		return err
	}

	err = n.SendEmail(subject, message)
	if err != nil {
		return err
	}

	err = n.notificationsRepository.Sent(notification.ID)

	if err != nil {
		return err
	}

	return err
}

// isThereANotificationToSend verifies wheter there are log events to notify the user
func (a *Service) ShouldSendNotification() bool {
	lastNotificationTime, err := a.FindLastEventLogsNotificationDate()

	return err != nil || time.Now().Sub(lastNotificationTime).Hours() > 12
}
