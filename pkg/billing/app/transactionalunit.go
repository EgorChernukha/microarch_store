package app

import "store/pkg/billing/domain"

type RepositoryProvider interface {
	UserAccountRepository() domain.UserAccountRepository
	ProcessedEventRepository() ProcessedEventRepository
}

type TransactionalUnit interface {
	RepositoryProvider
	Complete(err error) error
}

type TransactionalUnitFactory interface {
	NewTransactionalUnit() (TransactionalUnit, error)
}
