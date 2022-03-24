package app

import (
	uuid "github.com/satori/go.uuid"
)

const (
	Reserved = iota
	Confirmed
	Cancelled
)

type OrderPositionID uuid.UUID
type OrderID uuid.UUID
type PositionID uuid.UUID

type OrderPosition interface {
	ID() OrderPositionID
	OrderID() OrderID
	PositionID() PositionID
	Count() int
	Status() int
	Confirm()
	Cancel()
}

func NewOrderPosition(id OrderPositionID, orderID OrderID, positionID PositionID, count int, status int) OrderPosition {
	return &orderPosition{
		id:         id,
		orderID:    orderID,
		positionID: positionID,
		count:      count,
		status:     status,
	}
}

type orderPosition struct {
	id         OrderPositionID
	orderID    OrderID
	positionID PositionID
	count      int
	status     int
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

func (u *orderPosition) Status() int {
	return u.count
}

func (u *orderPosition) Confirm() {
	u.status = Confirmed
}

func (u *orderPosition) Cancel() {
	u.status = Cancelled
}
