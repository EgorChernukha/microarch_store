package domain

import uuid "github.com/satori/go.uuid"

type ID uuid.UUID
type UserID uuid.UUID

type UserAccount interface {
	ID() ID
	UserID() UserID
	Balance() float64
}

func NewUserAccount(id ID, userID UserID, balance float64) UserAccount {
	return &userAccount{id: id, userID: userID, balance: balance}
}

type userAccount struct {
	id      ID
	userID  UserID
	balance float64
}

func (u *userAccount) ID() ID {
	return u.id
}

func (u *userAccount) UserID() UserID {
	return u.userID
}

func (u *userAccount) Balance() float64 {
	return u.balance
}
