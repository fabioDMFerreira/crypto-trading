package notifications

import (
	"fmt"
	"net/smtp"
	"time"

	"github.com/fabiodmferreira/crypto-trading/db"
	"github.com/fabiodmferreira/crypto-trading/domain"
	"github.com/fabiodmferreira/crypto-trading/eventlogs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Notification
type Notification struct {
	ID                  primitive.ObjectID `bson:"_id" json:"_id"`
	Title               string             `json:"title"`
	Message             string             `json:"message"`
	To                  string             `json:"to"`
	DateCreated         time.Time          `json:"dateCreated"`
	NotificationType    string             `json:"notificationType"`
	NotificationChannel string             `json:"notificationChannel"`
}

type NotificationsService struct {
	NotificationsRepository *NotificationsRepository
	EventLogsRepository     *eventlogs.EventLogsRepository
	accountService          domain.AccountServiceReader
	Receiver                string
	Sender                  string
	SenderPassword          string
}

func NewNotificationsService(
	notificationsRepository *NotificationsRepository,
	eventLogsRepository *eventlogs.EventLogsRepository,
	accountService domain.AccountServiceReader,
	receiver string,
	sender string,
	senderPassword string,
) *NotificationsService {
	return &NotificationsService{notificationsRepository, eventLogsRepository, accountService, receiver, sender, senderPassword}
}

func (n *NotificationsService) SendEmail(subject, body string) error {
	from := n.Sender
	pass := n.SenderPassword
	to := n.Receiver

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n" +
		"MIME-Version: 1.0;\n" +
		"Content-Type: text/html;\n\n" +

		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	n.EventLogsRepository.Create("email", fmt.Sprintf("\"%v\" sent to %v", subject, n.Receiver))

	return err
}

func (n *NotificationsService) CheckEventLogs() error {
	lastNotificationTime, err := n.NotificationsRepository.FindLastEventLogsNotificationDate()

	if err != nil || time.Now().Sub(lastNotificationTime).Hours() > 12 {
		eventLogs, err := n.EventLogsRepository.FindAllToNotify()

		if err != nil {
			return err
		}

		pendingAssets, err := n.accountService.GetPendingAssets()

		if err != nil {
			return err
		}

		accountAmount, err := n.accountService.GetAmount()

		if err != nil {
			return err
		}

		startDate, endDate := lastNotificationTime, time.Now()
		balance, err := n.accountService.GetBalance(startDate, endDate)

		if err != nil {
			return err
		}

		subject := "Crypto-Trading: Report"
		var eventLogsIds []primitive.ObjectID

		for _, event := range *eventLogs {
			eventLogsIds = append(eventLogsIds, event.ID)
		}

		message, err := GenerateEventlogReportEmail(accountAmount, len(*pendingAssets), balance, startDate, endDate, eventLogs, pendingAssets)

		if err != nil {
			return err
		}

		err = n.SendEmail(subject, message.String())
		if err != nil {
			return err
		}

		notification := &Notification{
			ID:                  primitive.NewObjectID(),
			To:                  n.Receiver,
			Title:               subject,
			Message:             message.String(),
			DateCreated:         time.Now(),
			NotificationType:    "eventlogs",
			NotificationChannel: "email",
		}
		err = n.NotificationsRepository.Create(notification)
		if err != nil {
			return err
		}

		err = n.EventLogsRepository.MarkNotified(eventLogsIds)
		if err != nil {
			return err
		}
	}

	return nil
}

type NotificationsRepository struct {
	collection *mongo.Collection
}

// NewAssetsRepository returns an instance of OrdersRepository
func NewNotificationsRepository(collection *mongo.Collection) *NotificationsRepository {
	return &NotificationsRepository{
		collection,
	}
}

// Create inserts a new notification in collection
func (or *NotificationsRepository) Create(notification *Notification) error {
	ctx := db.NewMongoQueryContext()
	_, err := or.collection.InsertOne(ctx, notification)
	return err
}

func (or *NotificationsRepository) FindLastEventLogsNotificationDate() (time.Time, error) {
	ctx := db.NewMongoQueryContext()

	opts := options.FindOne().SetSort(bson.D{{"datecreated", -1}})
	var foundDocument Notification
	err := or.collection.FindOne(ctx, bson.D{{"notificationtype", "eventlogs"}}, opts).Decode(&foundDocument)

	if err != nil {
		return time.Now().AddDate(-1, 0, 0), err
	}

	return foundDocument.DateCreated, nil
}
