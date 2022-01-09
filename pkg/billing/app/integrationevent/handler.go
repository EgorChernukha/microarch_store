package integrationevent

import uuid "github.com/satori/go.uuid"

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
	case UserOrderCreated:
		return h.handleUserOrderCreated(concreteEvent.userID, concreteEvent.orderID, concreteEvent.price)
	case UserOrderCanceled:
		return h.handleUserOrderCanceled(concreteEvent.userID, concreteEvent.orderID, concreteEvent.price)
	default:
		return nil
	}
}

func (h handler) handleUserOrderCreated(userID, orderID uuid.UUID, price float64) error {
	return nil
}

func (h handler) handleUserOrderCanceled(userID, orderID uuid.UUID, price float64) error {
	return nil
}
