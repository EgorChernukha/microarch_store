package domain

import "errors"

type UserAccountRepository interface {
	NewID() ID
	Store(userAccount UserAccount) error
	FindOneByUserID(userID UserID) (UserAccount, error)
}

var ErrUserAccountNotFound = errors.New("user account not found")
