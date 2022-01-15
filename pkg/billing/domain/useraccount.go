package domain

import (
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

type ID uuid.UUID
type UserID uuid.UUID

type UserAccount interface {
	ID() ID
	UserID() UserID
	Balance() float64
	Withdraw(amount float64) error
	TopUp(amount float64) error
}

var ErrNotEnoughBalance = errors.New("not enough balance for payment")
var ErrInvalidAmount = errors.New("invalid amount")

func NewUserAccount(id ID, userID UserID, balance float64) UserAccount {
	return &userAccount{id: id, userID: userID, balance: balance}
}

type userAccount struct {
	id      ID
	userID  UserID
	balance float64
}

func (u *userAccount) Withdraw(amount float64) error {
	if amount < 0 {
		return ErrInvalidAmount
	}

	if u.balance < amount {
		return ErrNotEnoughBalance
	}

	u.balance = u.balance - amount
	return nil
}

func (u *userAccount) TopUp(amount float64) error {
	if amount < 0 {
		return ErrInvalidAmount
	}

	u.balance = u.balance + amount
	return nil
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
