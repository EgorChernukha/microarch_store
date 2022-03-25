package app

import (
	uuid "github.com/satori/go.uuid"

	"store/pkg/common/app/integrationevent"
)

type ProcessedEventRepository interface {
	SetProcessed(eventID integrationevent.EventUID) (alreadyProcessed bool, err error)
}

type UserEvent interface {
	UserID() uuid.UUID
}

func NewOrderConfirmedEvent(userID uuid.UUID, orderID uuid.UUID) UserEvent {
	return orderConfirmedEvent{userID: userID, orderID: orderID}
}

func NewOrderRejectedEvent(userID uuid.UUID, orderID uuid.UUID) UserEvent {
	return orderRejectedEvent{userID: userID, orderID: orderID}
}

type orderConfirmedEvent struct {
	userID  uuid.UUID
	orderID uuid.UUID
}

func (e orderConfirmedEvent) UserID() uuid.UUID {
	return e.userID
}

type orderRejectedEvent struct {
	userID  uuid.UUID
	orderID uuid.UUID
}

func (e orderRejectedEvent) UserID() uuid.UUID {
	return e.userID
}
