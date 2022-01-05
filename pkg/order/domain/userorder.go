package domain

import uuid "github.com/satori/go.uuid"

type ID uuid.UUID
type UserID uuid.UUID
type OrderID uuid.UUID

type UserOrder interface {
	ID() ID
	UserID() UserID
	OrderID() OrderID
	Price() float64
	Status() Status
}

func NewUserOrder(id ID, userID UserID, orderID OrderID, price float64, status Status) UserOrder {
	return &userOrder{
		id:      id,
		userID:  userID,
		orderID: orderID,
		price:   price,
		status:  status,
	}
}

type userOrder struct {
	id      ID
	userID  UserID
	orderID OrderID
	price   float64
	status  Status
}

func (u *userOrder) ID() ID {
	return u.id
}

func (u *userOrder) UserID() UserID {
	return u.userID
}

func (u *userOrder) OrderID() OrderID {
	return u.orderID
}

func (u *userOrder) Price() float64 {
	return u.price
}

func (u *userOrder) Status() Status {
	return u.status
}
