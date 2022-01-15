package app

import (
	uuid "github.com/satori/go.uuid"

	"store/pkg/common/app/integrationevent"
)

type ProcessedEventRepository interface {
	SetProcessed(uid integrationevent.EventID) (alreadyProcessed bool, err error)
}

type UserEvent interface {
	UserID() uuid.UUID
}

func NewUserRegisteredEvent(userID uuid.UUID, login string) UserEvent {
	return userRegisteredEvent{userID: userID, login: login}
}

type userRegisteredEvent struct {
	userID uuid.UUID
	login  string
}

func (e userRegisteredEvent) UserID() uuid.UUID {
	return e.userID
}

func (e userRegisteredEvent) Login() string {
	return e.login
}
