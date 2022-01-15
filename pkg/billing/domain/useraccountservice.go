package domain

import "github.com/pkg/errors"

type UserAccountService interface {
	CreateAccount(userID UserID) error
	TopUpAccount(userID UserID, amount float64) error
	ProcessPayment(userID UserID, amount float64) error
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

func (u *userAccountService) TopUpAccount(userID UserID, amount float64) error {
	account, err := u.repository.FindOneByUserID(userID)
	if err != nil {
		return err
	}

	account.TopUp(amount)
	return u.repository.Store(account)
}

func (u *userAccountService) ProcessPayment(userID UserID, amount float64) error {
	account, err := u.repository.FindOneByUserID(userID)
	if err != nil {
		return err
	}

	account.Withdraw(amount)
	return u.repository.Store(account)
}
