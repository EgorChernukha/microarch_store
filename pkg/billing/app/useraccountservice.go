package app

import (
	uuid "github.com/satori/go.uuid"

	"store/pkg/billing/domain"
)

type UserAccountService interface {
	CreateAccount(userID uuid.UUID) error
	TopUpAccount(userID uuid.UUID, amount float64) error
	ProcessPayment(userID uuid.UUID, amount float64) error
}

func NewUserAccountService(domainService domain.UserAccountService) UserAccountService {
	return &userAccountService{domainService: domainService}
}

type userAccountService struct {
	domainService domain.UserAccountService
}

func (u *userAccountService) ProcessPayment(userID uuid.UUID, amount float64) error {
	return u.domainService.ProcessPayment(domain.UserID(userID), amount)
}

func (u *userAccountService) TopUpAccount(userID uuid.UUID, amount float64) error {
	return u.domainService.TopUpAccount(domain.UserID(userID), amount)
}

func (u *userAccountService) CreateAccount(userID uuid.UUID) error {
	return u.domainService.CreateAccount(domain.UserID(userID))
}
