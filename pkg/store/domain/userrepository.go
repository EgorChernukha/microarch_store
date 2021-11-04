package domain

import "errors"

type UserRepository interface {
	Store(user User) (UserID, error)
	Remove(user User) error
	FindOne(id UserID) (User, error)
}

var ErrUserNotFound = errors.New("user not found")
