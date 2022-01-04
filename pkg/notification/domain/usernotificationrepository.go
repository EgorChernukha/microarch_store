package domain

type UserNotificationRepository interface {
	Store(userNotification UserNotification) error
}
