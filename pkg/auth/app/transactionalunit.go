package app

import "store/pkg/common/app/storedevent"

type RepositoryProvider interface {
	UserRepository() UserRepository
	EventStore() storedevent.EventStore
}

type TransactionalUnit interface {
	RepositoryProvider
	Complete(err error) error
}

type TransactionalUnitFactory interface {
	NewTransactionalUnit() (TransactionalUnit, error)
}
