package domain

import (
	uuid "github.com/satori/go.uuid"
)

type UserID uuid.UUID
type OrderID uuid.UUID

type UserNotification interface {
	UserID() UserID
	OrderID() OrderID
	Message() string
}

func NewUserNotification(userID UserID, orderID OrderID, message string) UserNotification {
	return &userNotification{
		userID:  userID,
		orderID: orderID,
		message: message,
	}
}

type userNotification struct {
	userID  UserID
	orderID OrderID
	message string
}

func (n *userNotification) UserID() UserID {
	return n.userID
}

func (n *userNotification) OrderID() OrderID {
	return n.orderID
}

func (n *userNotification) Message() string {
	return n.message
}
