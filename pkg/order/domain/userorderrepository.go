package domain

import "errors"

type UserOrderRepository interface {
	NewID() ID
	Store(userOrder UserOrder) error
	FindOneByOrderID(orderID OrderID) (UserOrder, error)
}

var ErrUserOrderNotFound = errors.New("user order not found")
