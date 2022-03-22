package app

type RepositoryProvider interface {
	ProcessedEventRepository() ProcessedEventRepository
	OrderDeliveryRepository() OrderDeliveryRepository
}

type TransactionalUnit interface {
	RepositoryProvider
	Complete(err error) error
}

type TransactionalUnitFactory interface {
	NewTransactionalUnit() (TransactionalUnit, error)
}
