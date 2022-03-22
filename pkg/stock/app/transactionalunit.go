package app

type RepositoryProvider interface {
	ProcessedEventRepository() ProcessedEventRepository
	PositionRepository() PositionRepository
	OrderPositionRepository() OrderPositionRepository
}

type TransactionalUnit interface {
	RepositoryProvider
	Complete(err error) error
}

type TransactionalUnitFactory interface {
	NewTransactionalUnit() (TransactionalUnit, error)
}
