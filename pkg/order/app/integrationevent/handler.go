package integrationevent

import (
	uuid "github.com/satori/go.uuid"
)

const (
	paymentSucceededMessage = "payment succeeded"
	paymentFailedMessage    = "payment failed"
)

type Handler interface {
	Handle(event IntegrationEvent) error
}

func NewHandler() Handler {
	return handler{}
}

type handler struct {
}

func (h handler) Handle(event IntegrationEvent) error {
	switch concreteEvent := event.(type) {
	case UserOrderPaymentSucceeded:
		return h.handleUserOrderPaymentSucceeded(concreteEvent.userID, concreteEvent.orderID)
	case UserOrderPaymentFailed:
		return h.handleUserOrderPaymentFailed(concreteEvent.userID, concreteEvent.orderID)
	default:
		return nil
	}
}

func (h handler) handleUserOrderPaymentSucceeded(userID, orderID uuid.UUID) error {
	return nil
}

func (h handler) handleUserOrderPaymentFailed(userID, orderID uuid.UUID) error {
	return nil
}
