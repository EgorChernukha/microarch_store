package app

import (
	uuid "github.com/satori/go.uuid"

	"store/pkg/notification/domain"
)

type UserNotificationService interface {
	CreateUserNotification(userID uuid.UUID, orderID uuid.UUID, email, message string) error
}

func NewUserNotificationService(domainService domain.UserNotificationService) UserNotificationService {
	return &userNotificationService{domainService: domainService}
}

type userNotificationService struct {
	domainService domain.UserNotificationService
}

func (s *userNotificationService) CreateUserNotification(userID uuid.UUID, orderID uuid.UUID, email, message string) error {
	return s.domainService.CreateUserNotification(domain.UserID(userID), domain.OrderID(orderID), email, message)
}
