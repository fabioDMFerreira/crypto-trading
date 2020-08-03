package notifications

import (
	"net/smtp"
	"time"

	"github.com/fabiodmferreira/crypto-trading/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	notificationsRepository domain.NotificationsRepository
	options                 domain.NotificationOptions
	sendMail                domain.SendMail
	appID                   primitive.ObjectID
}

func NewService(
	notificationsRepository domain.NotificationsRepository,
	options domain.NotificationOptions,
	sendMail domain.SendMail,
	appID primitive.ObjectID,
) *Service {
	return &Service{notificationsRepository, options, sendMail, appID}
}

// SendEmail setup an email options and sends it
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

	err := n.sendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	return err
}

// FindLastEventLogsNotificationDate returns last notification date
func (n *Service) FindLastEventLogsNotificationDate() (time.Time, error) {
	return n.notificationsRepository.FindLastEventLogsNotificationDate()
}

// CreateEmailNotification stores notification in repository and send email to Receiver
func (n *Service) CreateEmailNotification(subject, message, notificationType string) error {
	notification := &domain.Notification{
		ID:                  primitive.NewObjectID(),
		To:                  n.options.Receiver,
		Title:               subject,
		Message:             message,
		CreatedAt:           time.Now(),
		NotificationType:    notificationType,
		NotificationChannel: "email",
		ApplicationID:       n.appID,
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

// ShouldSendNotification verifies wheter last notification was sent more than 12 hours ago
func (a *Service) ShouldSendNotification() bool {
	lastNotificationTime, err := a.FindLastEventLogsNotificationDate()

	return err != nil || time.Now().Sub(lastNotificationTime).Hours() > 12
}
