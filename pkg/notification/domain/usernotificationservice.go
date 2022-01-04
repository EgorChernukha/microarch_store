package domain

import "github.com/pkg/errors"

type UserNotificationService interface {
	CreateUserNotification(userID UserID, orderID OrderID, email, message string) error
}

func NewUserNotificationService(repository UserNotificationRepository) UserNotificationService {
	return &userNotificationService{repository: repository}
}

type userNotificationService struct {
	repository UserNotificationRepository
}

func (s *userNotificationService) CreateUserNotification(userID UserID, orderID OrderID, email, message string) error {
	userNotification := NewUserNotification(userID, orderID, email, message)
	err := s.repository.Store(userNotification)

	return errors.Wrap(err, "failed to create user notification")
}
