package mocks

import "time"

// NotificationsServiceSpy mocks Notifications service
type NotificationsServiceSpy struct {
	emailsNotifications int
}

// CreateEmailNotification increases the counter of emails notifications
func (n *NotificationsServiceSpy) CreateEmailNotification(subject, message, notificationType string) error {
	n.emailsNotifications++
	return nil
}

// FindLastEventLogsNotificationDate returns current date
func (n *NotificationsServiceSpy) FindLastEventLogsNotificationDate() (time.Time, error) {
	return time.Now(), nil
}
