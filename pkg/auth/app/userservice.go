package app

import (
	stderrors "errors"
	"store/pkg/common/app/storedevent"

	"github.com/pkg/errors"
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

func NewUserService(trUnitFactory TransactionalUnitFactory, userReadRepository UserRepository, eventSender storedevent.Sender, passwordEncoder PasswordEncoder) UserService {
	return &userService{trUnitFactory: trUnitFactory, userReadRepository: userReadRepository, eventSender: eventSender, passwordEncoder: passwordEncoder}
}

type userService struct {
	trUnitFactory      TransactionalUnitFactory
	userReadRepository UserRepository
	eventSender        storedevent.Sender
	passwordEncoder    PasswordEncoder
}

func (u *userService) AddUser(login, password string) (UserID, error) {
	id := UserID(uuid.NewV1())
	if err := u.validateLogin(login); err != nil {
		return id, err
	}

	if user, err := u.userReadRepository.FindOneByLogin(login); errors.Cause(err) != ErrUserNotFound || user != nil {
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

	err := u.executeInTransaction(func(provider RepositoryProvider) error {
		err2 := provider.UserRepository().Store(&user)
		if err2 != nil {
			return err2
		}

		event := NewUserRegisteredEvent(id, login)
		err2 = provider.EventStore().Add(event)
		if err2 != nil {
			return err2
		}
		u.eventSender.EventStored(event.UID)

		return nil
	})
	if err != nil {
		return id, err
	}

	u.eventSender.SendStoredEvents()
	return id, err
}

func (u *userService) RemoveUser(id UserID) error {
	return u.userReadRepository.Remove(id)
}

func (u *userService) FindUserByID(id UserID) (*User, error) {
	return u.userReadRepository.FindOneByID(id)
}

func (u *userService) FindUserByLoginAndPassword(login, password string) (*User, error) {
	user, err := u.userReadRepository.FindOneByLogin(login)
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

func (u *userService) executeInTransaction(f func(RepositoryProvider) error) (err error) {
	var trUnit TransactionalUnit
	trUnit, err = u.trUnitFactory.NewTransactionalUnit()
	if err != nil {
		return err
	}
	defer func() {
		err = trUnit.Complete(err)
	}()
	err = f(trUnit)
	return err
}

var ErrUserAlreadyExists = stderrors.New("user already exists")
var ErrInvalidLogin = stderrors.New("invalid login")
var ErrInvalidPassword = stderrors.New("invalid password")
