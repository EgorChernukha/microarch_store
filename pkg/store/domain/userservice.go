package domain

import "github.com/pkg/errors"

type UserService interface {
	UpdateUser(id UserID, username, firstname, lastname, email, phone string) error
	RemoveUser(id UserID) error
}

func NewUserService(repository UserRepository) UserService {
	return &userService{repository: repository}
}

type userService struct {
	repository UserRepository
}

func (u *userService) UpdateUser(id UserID, username, firstname, lastname, email, phone string) error {
	user, err := u.repository.FindOne(id)
	if err != nil && errors.Cause(err) != ErrUserNotFound {
		return errors.Wrap(err, "failed to fetch user")
	}

	if user == nil {
		user := NewUser(id, username, firstname, lastname, email, phone)
		err = u.repository.Store(user)
	} else {
		user.Update(firstname, lastname, email, phone)
		err = u.repository.Store(user)
	}

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
