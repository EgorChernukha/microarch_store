package app

import (
	"errors"
	"time"

	uuid "github.com/satori/go.uuid"
)

type UserOrderData struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	OrderID   uuid.UUID `json:"order_id"`
	Price     float64   `json:"price"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserOrderQueryService interface {
	FindUserOrderByOrderID(orderID uuid.UUID) (UserOrderData, error)
	ListUserOrdersByUserIDs(userID uuid.UUID) ([]UserOrderData, error)
	ListUserOrders() ([]UserOrderData, error)
}

var ErrUserOrderNotExists = errors.New("user order does not exist")
