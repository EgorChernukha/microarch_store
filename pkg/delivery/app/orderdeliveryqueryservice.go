package app

import (
	"errors"
	"time"

	uuid "github.com/satori/go.uuid"
)

type OrderDeliveryData struct {
	ID        uuid.UUID `json:"id"`
	OrderID   uuid.UUID `json:"order_id"`
	UserID    uuid.UUID `json:"user_id"`
	status    int       `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
}

type OrderDeliveryQueryService interface {
	FindByOrderID(orderID uuid.UUID) (OrderDeliveryData, error)
	FindByUserID(userID uuid.UUID) ([]OrderDeliveryData, error)
	ListOrderDelivery() ([]OrderDeliveryData, error)
}

var ErrOrderDeliveryNotExists = errors.New("order delivery does not exist")
