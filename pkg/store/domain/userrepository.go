package domain

import "errors"

type UserRepository interface {
	NewID() UserID
	Store(user User) error
	Remove(user User) error
	FindOne(id UserID) (User, error)
}

var ErrUserNotFound = errors.New("user not found")
