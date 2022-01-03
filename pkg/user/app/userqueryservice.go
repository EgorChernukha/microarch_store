package app

import (
	"errors"

	uuid "github.com/satori/go.uuid"
)

type UserData struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Firstname string    `json:"firstname"`
	Lastname  string    `json:"lastname"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
}

type UserQueryService interface {
	FindUser(id uuid.UUID) (UserData, error)
}

var ErrUserNotExists = errors.New("user does not exist")
