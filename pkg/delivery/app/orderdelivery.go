package app

import (
	uuid "github.com/satori/go.uuid"
)

type ID uuid.UUID
type OrderID uuid.UUID
type UserID uuid.UUID

const (
	Created = iota
	Sent
	Received
	Rejected
)

type OrderDelivery interface {
	ID() ID
	OrderID() OrderID
	UserID() UserID
	Status() int
	MarkAsSent()
	MarkAsReceived()
	MarkAsRejected()
}

func NewOrderDelivery(id ID, orderID OrderID, userID UserID, status int) OrderDelivery {
	return &orderDelivery{
		id:      id,
		orderID: orderID,
		userID:  userID,
		status:  status,
	}
}

type orderDelivery struct {
	id      ID
	orderID OrderID
	userID  UserID
	status  int
}

func (o *orderDelivery) ID() ID {
	return o.id
}

func (o *orderDelivery) OrderID() OrderID {
	return o.orderID
}

func (o *orderDelivery) UserID() UserID {
	return o.userID
}

func (o *orderDelivery) Status() int {
	return o.status
}

func (o *orderDelivery) MarkAsSent() {
	o.status = Sent
}

func (o *orderDelivery) MarkAsReceived() {
	o.status = Received
}

func (o *orderDelivery) MarkAsRejected() {
	o.status = Rejected
}
