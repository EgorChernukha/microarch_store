package app

import (
	"errors"
	"time"

	uuid "github.com/satori/go.uuid"
)

type UserAccountData struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Balance   float64   `json:"balance"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserAccountQueryService interface {
	FindUserAccountByUserID(userID uuid.UUID) (UserAccountData, error)
	ListUserAccounts() ([]UserAccountData, error)
}

var ErrUserAccountNotExists = errors.New("user account does not exist")
