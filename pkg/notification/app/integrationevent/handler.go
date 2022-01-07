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
	case UserRegistered:
		return h.handleUserRegistered(concreteEvent.accountID, concreteEvent.userID, concreteEvent.departmentID)
	default:
		return nil
	}
}

func (h handler) handleUserRegistered(accountID string, userID, departmentID uuid.UUID) error {
	return nil
}
