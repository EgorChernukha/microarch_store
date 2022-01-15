package app

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type UserNotificationData struct {
	UserID    uuid.UUID `json:"user_id"`
	OrderID   uuid.UUID `json:"order_id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

type UserNotificationQueryService interface {
	ListUserNotifications(userID uuid.UUID) ([]UserNotificationData, error)
}
