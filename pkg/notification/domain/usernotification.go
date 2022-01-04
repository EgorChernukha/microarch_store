package domain

import (
	uuid "github.com/satori/go.uuid"
)

type UserID uuid.UUID
type OrderID uuid.UUID

type UserNotification interface {
	UserID() UserID
	OrderID() OrderID
	Email() string
	Message() string
}

func NewUserNotification(userID UserID, orderID OrderID, email, message string) UserNotification {
	return &userNotification{
		userID:  userID,
		orderID: orderID,
		email:   email,
		message: message,
	}
}

type userNotification struct {
	userID  UserID
	orderID OrderID
	email   string
	message string
}

func (n *userNotification) UserID() UserID {
	return n.userID
}

func (n *userNotification) OrderID() OrderID {
	return n.orderID
}

func (n *userNotification) Email() string {
	return n.email
}

func (n *userNotification) Message() string {
	return n.message
}
