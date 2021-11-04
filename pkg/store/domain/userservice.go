package domain

import "github.com/pkg/errors"

type UserService interface {
	AddUser(username, firstname, lastname, email, phone string) (UserID, error)
	UpdateUser(id UserID, firstname, lastname, email, phone string) error
	RemoveUser(id UserID) error
}

func NewUserService(repository UserRepository) UserService {
	return &userService{repository: repository}
}

type userService struct {
	repository UserRepository
}

func (u *userService) AddUser(username, firstname, lastname, email, phone string) (UserID, error) {
	user := NewUser(username, firstname, lastname, email, phone)
	id, err := u.repository.Store(user)

	return id, errors.Wrap(err, "failed to add new user")
}

func (u *userService) UpdateUser(id UserID, firstname, lastname, email, phone string) error {
	user, err := u.repository.FindOne(id)
	if err != nil {
		return errors.Wrap(err, "failed to fetch user")
	}

	user.Update(firstname, lastname, email, phone)
	_, err = u.repository.Store(user)

	return errors.Wrap(err, "failed to update user")
}

func (u *userService) RemoveUser(id UserID) error {
	user, err := u.repository.FindOne(id)
	if err != nil {
		return errors.Wrap(err, "failed to fetch user")
	}

	err = u.repository.Remove(user)

	return errors.Wrap(err, "failed to remove user")
}

func (u *userService) getUser(id UserID) (User, error) {
	return u.repository.FindOne(id)
}
