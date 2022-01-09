package integrationevent

import (
	uuid "github.com/satori/go.uuid"
)

type IntegrationEvent interface {
}

type UserOrderPaymentSucceeded struct {
	userID  uuid.UUID
	orderID uuid.UUID
}

func NewUserOrderPaymentSucceeded(userID, orderID uuid.UUID) UserOrderPaymentSucceeded {
	return UserOrderPaymentSucceeded{userID, orderID}
}

type UserOrderPaymentFailed struct {
	userID  uuid.UUID
	orderID uuid.UUID
}

func NewUserOrderPaymentFailed(userID, orderID uuid.UUID) UserOrderPaymentFailed {
	return UserOrderPaymentFailed{userID, orderID}
}
