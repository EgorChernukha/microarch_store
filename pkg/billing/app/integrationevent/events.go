package integrationevent

import uuid "github.com/satori/go.uuid"

type IntegrationEvent interface {
}

type UserOrderCreated struct {
	userID  uuid.UUID
	orderID uuid.UUID
	price   float64
}

func NewUserOrderCreated(userID, orderID uuid.UUID, price float64) UserOrderCreated {
	return UserOrderCreated{userID, orderID, price}
}

type UserOrderCanceled struct {
	userID  uuid.UUID
	orderID uuid.UUID
	price   float64
}

func NewUserOrderCanceled(userID, orderID uuid.UUID, price float64) UserOrderCanceled {
	return UserOrderCanceled{userID, orderID, price}
}
