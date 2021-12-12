package app

import (
	uuid "github.com/satori/go.uuid"

	"store/pkg/store/domain"
)

type UserService interface {
	UpdateUser(id uuid.UUID, username, firstname, lastname, email, phone string) error
	RemoveUser(id uuid.UUID) error
}

func NewUserService(domainService domain.UserService) UserService {
	return &userService{domainService: domainService}
}

type userService struct {
	domainService domain.UserService
}

func (u *userService) UpdateUser(id uuid.UUID, username, firstname, lastname, email, phone string) error {
	return u.domainService.UpdateUser(domain.UserID(id), username, firstname, lastname, email, phone)
}

func (u *userService) RemoveUser(id uuid.UUID) error {
	return u.domainService.RemoveUser(domain.UserID(id))
}
