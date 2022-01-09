package integrationevent

import (
	uuid "github.com/satori/go.uuid"

	"store/pkg/notification/app"
)

const (
	paymentSucceededMessage = "payment succeeded"
	paymentFailedMessage    = "payment failed"
)

type Handler interface {
	Handle(event IntegrationEvent) error
}

func NewHandler(userNotificationService app.UserNotificationService) Handler {
	return handler{userNotificationService: userNotificationService}
}

type handler struct {
	userNotificationService app.UserNotificationService
}

func (h handler) Handle(event IntegrationEvent) error {
	switch concreteEvent := event.(type) {
	case UserOrderPaymentSucceeded:
		return h.handleUserOrderPaymentSucceeded(concreteEvent.userID, concreteEvent.orderID, concreteEvent.email)
	case UserOrderPaymentFailed:
		return h.handleUserOrderPaymentFailed(concreteEvent.userID, concreteEvent.orderID, concreteEvent.email)
	default:
		return nil
	}
}

func (h handler) handleUserOrderPaymentSucceeded(userID, orderID uuid.UUID, email string) error {
	return h.userNotificationService.CreateUserNotification(userID, orderID, email, paymentSucceededMessage)
}

func (h handler) handleUserOrderPaymentFailed(userID, orderID uuid.UUID, email string) error {
	return h.userNotificationService.CreateUserNotification(userID, orderID, email, paymentFailedMessage)
}
