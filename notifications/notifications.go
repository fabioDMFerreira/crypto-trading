package notifications

import (
	"net/smtp"
	"time"

	"github.com/fabiodmferreira/crypto-trading/db"
	"github.com/fabiodmferreira/crypto-trading/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Notification domain.Notification

type NotificationsService struct {
	NotificationsRepository *NotificationsRepository
	Receiver                string
	Sender                  string
	SenderPassword          string
}

func NewNotificationsService(
	notificationsRepository *NotificationsRepository,
	receiver string,
	sender string,
	senderPassword string,
) *NotificationsService {
	return &NotificationsService{notificationsRepository, receiver, sender, senderPassword}
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

	return err
}

func (n *NotificationsService) FindLastEventLogsNotificationDate() (time.Time, error) {
	return n.NotificationsRepository.FindLastEventLogsNotificationDate()
}

func (n *NotificationsService) CreateEmailNotification(subject, message, notificationType string) error {
	notification := &Notification{
		ID:                  primitive.NewObjectID(),
		To:                  n.Receiver,
		Title:               subject,
		Message:             message,
		CreatedAt:           time.Now(),
		NotificationType:    notificationType,
		NotificationChannel: "email",
	}

	err := n.NotificationsRepository.Create(notification)

	if err != nil {
		return err
	}

	err = n.SendEmail(subject, message)
	if err != nil {
		return err
	}

	err = n.NotificationsRepository.Sent(notification.ID)

	if err != nil {
		return err
	}

	return err
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

	opts := options.FindOne().SetSort(bson.D{{"createdat", -1}})
	var foundDocument Notification
	err := or.collection.FindOne(ctx, bson.D{{"notificationtype", "eventlogs"}}, opts).Decode(&foundDocument)

	if err != nil {
		return time.Now().AddDate(-1, 0, 0), err
	}

	return foundDocument.CreatedAt, nil
}

func (or *NotificationsRepository) Sent(id primitive.ObjectID) error {
	ctx := db.NewMongoQueryContext()

	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"send", true}}}}
	_, err := or.collection.UpdateOne(ctx, filter, update)

	return err
}
