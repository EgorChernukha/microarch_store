package integrationevent

import (
	uuid "github.com/satori/go.uuid"
)

type IntegrationEvent interface {
}

type UserOrderPaymentSucceeded struct {
	userID  uuid.UUID
	orderID uuid.UUID
	email   string
}

func NewUserOrderPaymentSucceeded(userID, orderID uuid.UUID, email string) UserOrderPaymentSucceeded {
	return UserOrderPaymentSucceeded{userID, orderID, email}
}

type UserOrderPaymentFailed struct {
	userID  uuid.UUID
	orderID uuid.UUID
	email   string
}

func NewUserOrderPaymentFailed(userID, orderID uuid.UUID, email string) UserOrderPaymentFailed {
	return UserOrderPaymentFailed{userID, orderID, email}
}
