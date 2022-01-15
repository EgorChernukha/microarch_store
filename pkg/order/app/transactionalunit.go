package app

import (
	"store/pkg/common/app/storedevent"
)

type RepositoryProvider interface {
	UserOrderRepository() UserOrderRepository
	EventStore() storedevent.EventStore
}

type TransactionalUnit interface {
	RepositoryProvider
	Complete(err error) error
}

type TransactionalUnitFactory interface {
	NewTransactionalUnit() (TransactionalUnit, error)
}
