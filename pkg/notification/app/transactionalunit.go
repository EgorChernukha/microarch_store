package app

import "store/pkg/notification/domain"

type RepositoryProvider interface {
	UserNotificationRepository() domain.UserNotificationRepository
	ProcessedEventRepository() ProcessedEventRepository
}

type TransactionalUnit interface {
	RepositoryProvider
	Complete(err error) error
}

type TransactionalUnitFactory interface {
	NewTransactionalUnit() (TransactionalUnit, error)
}
