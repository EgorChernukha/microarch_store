package app

import (
	stderrors "errors"

	uuid "github.com/satori/go.uuid"
)

type UserService interface {
	AddUser(login, password string) (UserID, error)
	RemoveUser(id UserID) error
	FindUserByID(id UserID) (*User, error)
	FindUserByLoginAndPassword(login, password string) (*User, error)
}

type PasswordEncoder interface {
	Encode(password string, userID UserID) string
}

func NewUserService(userRepository UserRepository, passwordEncoder PasswordEncoder) UserService {
	return &userService{userRepository: userRepository, passwordEncoder: passwordEncoder}
}

type userService struct {
	userRepository  UserRepository
	passwordEncoder PasswordEncoder
}

func (u *userService) AddUser(login, password string) (UserID, error) {
	id := UserID(uuid.NewV1())
	if err := u.validateLogin(login); err != nil {
		return id, err
	}

	if user, err := u.userRepository.FindOneByLogin(login); err != ErrUserNotFound || user != nil {
		if user != nil {
			return id, ErrUserAlreadyExists
		}
		return id, err
	}

	user := User{
		ID:       id,
		Login:    login,
		Password: u.passwordEncoder.Encode(password, id),
	}

	return id, u.userRepository.Store(&user)
}

func (u *userService) RemoveUser(id UserID) error {
	return u.userRepository.Remove(id)
}

func (u *userService) FindUserByID(id UserID) (*User, error) {
	return u.userRepository.FindOneByID(id)
}

func (u *userService) FindUserByLoginAndPassword(login, password string) (*User, error) {
	user, err := u.userRepository.FindOneByLogin(login)
	if err != nil {
		return nil, err
	}
	encodedPass := u.passwordEncoder.Encode(password, user.ID)
	if encodedPass != user.Password {
		return nil, ErrInvalidPassword
	}
	return user, nil
}
func (u *userService) validateLogin(login string) error {
	if len(login) > maxLoginLength {
		return ErrInvalidLogin
	}

	return nil
}

var ErrUserAlreadyExists = stderrors.New("user already exists")
var ErrInvalidLogin = stderrors.New("invalid login")
var ErrInvalidPassword = stderrors.New("invalid password")
