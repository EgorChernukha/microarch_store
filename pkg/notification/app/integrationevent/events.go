package integrationevent

import (
	uuid "github.com/satori/go.uuid"
)

type IntegrationEvent interface {
	AccountID() string
}

type UserRegistered struct {
	accountID    string
	userID       uuid.UUID
	departmentID uuid.UUID
}

func (event UserRegistered) AccountID() string {
	return event.accountID
}

func NewUserRegistered(accountID string, userID, departmentID uuid.UUID) UserRegistered {
	return UserRegistered{accountID, userID, departmentID}
}
