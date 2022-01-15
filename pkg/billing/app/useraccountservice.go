package app

import (
	uuid "github.com/satori/go.uuid"

	"store/pkg/billing/domain"
)

type UserAccountService interface {
	CreateAccount(userID uuid.UUID) error
}

func NewUserAccountService(domainService domain.UserAccountService) UserAccountService {
	return &userAccountService{domainService: domainService}
}

type userAccountService struct {
	domainService domain.UserAccountService
}

func (u *userAccountService) CreateAccount(userID uuid.UUID) error {
	return u.domainService.CreateAccount(domain.UserID(userID))
}
