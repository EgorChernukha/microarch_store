package app

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type UserID uuid.UUID
type OrderID uuid.UUID

type NotificationType int

const (
	TypeOrderConfirmed NotificationType = 1
	TypeOrderRejected  NotificationType = 2
)

type Notification struct {
	UserID       UserID
	OrderID      OrderID
	Message      string
	Type         NotificationType
	CreationDate time.Time
}

type NotificationRepository interface {
	Store(notification *Notification) error
	FindAllByUserID(userID UserID) ([]Notification, error)
}
