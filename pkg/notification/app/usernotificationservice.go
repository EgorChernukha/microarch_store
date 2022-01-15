package app

import (
	"fmt"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"store/pkg/notification/domain"
)

type UserNotificationService interface {
	CreateUserNotification(notificationType NotificationType, userID uuid.UUID, orderID uuid.UUID) error
}

func NewUserNotificationService(domainService domain.UserNotificationService) UserNotificationService {
	return &userNotificationService{domainService: domainService}
}

type userNotificationService struct {
	domainService domain.UserNotificationService
}

func (s *userNotificationService) CreateUserNotification(notificationType NotificationType, userID uuid.UUID, orderID uuid.UUID) error {
	message, err := messageFromNotificationType(notificationType)
	if err != nil {
		return errors.WithStack(err)
	}

	return s.domainService.CreateUserNotification(domain.UserID(userID), domain.OrderID(orderID), message)
}

func messageFromNotificationType(notificationType NotificationType) (string, error) {
	switch notificationType {
	case TypeOrderConfirmed:
		return fmt.Sprintf("Order confirmed"), nil
	case TypeOrderRejected:
		return fmt.Sprintf("Order rejected"), nil
	default:
		return "", errors.New("unknown notification type")
	}
}
