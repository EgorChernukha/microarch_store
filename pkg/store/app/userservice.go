package app

import "mod/pkg/store/domain"

type UserService interface {
	AddUser(username, firstname, lastname, email, phone string) (int, error)
	UpdateUser(id int, firstname, lastname, email, phone string) error
	RemoveUser(id int) error
}

func NewUserService(domainService domain.UserService) UserService {
	return &userService{domainService: domainService}
}

type userService struct {
	domainService domain.UserService
}

func (u *userService) AddUser(username, firstname, lastname, email, phone string) (int, error) {
	userID, err := u.domainService.AddUser(username, firstname, lastname, email, phone)

	return int(userID), err
}

func (u *userService) UpdateUser(id int, firstname, lastname, email, phone string) error {
	return u.domainService.UpdateUser(domain.UserID(id), firstname, lastname, email, phone)
}

func (u *userService) RemoveUser(id int) error {
	return u.domainService.RemoveUser(domain.UserID(id))
}
