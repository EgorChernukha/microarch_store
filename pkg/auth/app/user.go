package app

import (
	"errors"

	uuid "github.com/satori/go.uuid"
)

type UserID uuid.UUID

const maxLoginLength = 255

type User struct {
	ID       UserID
	Login    string
	Password string
}

type UserRepository interface {
	Store(user *User) error
	Remove(id UserID) error
	FindOneByID(id UserID) (*User, error)
	FindOneByLogin(login string) (*User, error)
}

var ErrUserNotFound = errors.New("user not found")
