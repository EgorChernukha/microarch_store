package app

import "errors"

type OrderPositionRepository interface {
	NewID() OrderPositionID
	Store(orderPosition OrderPosition) error
	FindByOrderID(orderID OrderID) ([]OrderPosition, error)
	FindByPositionID(positionID PositionID) ([]OrderPosition, error)
	FindByOrderIDAndPositionID(orderID OrderID, positionID PositionID) (OrderPosition, error)
}

var ErrOrderPositionNotFound = errors.New("order position not found")
