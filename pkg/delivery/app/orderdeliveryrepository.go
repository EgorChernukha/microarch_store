package app

import "errors"

type OrderDeliveryRepository interface {
	NewID() ID
	Store(orderDelivery OrderDelivery) error
	FindByOrderID(orderID OrderID) (OrderDelivery, error)
	FindByUserID(userID UserID) ([]OrderDelivery, error)
}

var ErrOrderDeliveryNotFound = errors.New("order delivery not found")
