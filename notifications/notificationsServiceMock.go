package notifications

import "time"

// NotificationsMock mocks Notifications service
type NotificationsMock struct {
	emailsNotifications int
}

// CreateEmailNotification increases the counter of emails notifications
func (n *NotificationsMock) CreateEmailNotification(subject, message, notificationType string) error {
	n.emailsNotifications++
	return nil
}

// FindLastEventLogsNotificationDate returns current date
func (n *NotificationsMock) FindLastEventLogsNotificationDate() (time.Time, error) {
	return time.Now(), nil
}
