package domain

import "github.com/pkg/errors"

type UserAccountService interface {
	CreateAccount(userID UserID) error
}

func NewUserAccountService(repository UserAccountRepository) UserAccountService {
	return &userAccountService{repository: repository}
}

type userAccountService struct {
	repository UserAccountRepository
}

func (u *userAccountService) CreateAccount(userID UserID) error {
	userAccount := NewUserAccount(u.repository.NewID(), userID, 0)
	err := u.repository.Store(userAccount)

	return errors.Wrap(err, "failed to create user account")
}
