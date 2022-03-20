package app

import (
	uuid "github.com/satori/go.uuid"
)

type OrderPositionID uuid.UUID
type OrderID uuid.UUID
type PositionID uuid.UUID

type OrderPosition interface {
	ID() OrderPositionID
	OrderID() OrderID
	PositionID() PositionID
	Count() int
}

func NewOrderPosition(id OrderPositionID, orderID OrderID, positionID PositionID, count int) OrderPosition {
	return &orderPosition{
		id:         id,
		orderID:    orderID,
		positionID: positionID,
		count:      count,
	}
}

type orderPosition struct {
	id         OrderPositionID
	orderID    OrderID
	positionID PositionID
	count      int
}

func (u *orderPosition) ID() OrderPositionID {
	return u.id
}

func (u *orderPosition) OrderID() OrderID {
	return u.orderID
}

func (u *orderPosition) PositionID() PositionID {
	return u.positionID
}

func (u *orderPosition) Count() int {
	return u.count
}
