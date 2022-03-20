package mysql

import (
	"github.com/pkg/errors"

	commonmysql "store/pkg/common/infrastructure/mysql"

	"store/pkg/stock/app"
)

func NewTransactionalUnitFactory(client commonmysql.TransactionalClient) app.TransactionalUnitFactory {
	return &transactionalUnitFactory{client: client}
}

type transactionalUnitFactory struct {
	client commonmysql.TransactionalClient
}

func (d *transactionalUnitFactory) NewTransactionalUnit() (app.TransactionalUnit, error) {
	transaction, err := d.client.BeginTransaction()
	if err != nil {
		return nil, err
	}
	return &transactionalUnit{transaction: transaction}, nil
}

type transactionalUnit struct {
	transaction commonmysql.Transaction
}

func (t *transactionalUnit) ProcessedEventRepository() app.ProcessedEventRepository {
	return NewProcessedEventRepository(t.transaction)
}

func (t *transactionalUnit) Complete(err error) error {
	if err != nil {
		rollbackErr := t.transaction.Rollback()
		if rollbackErr != nil {
			return errors.Wrap(err, rollbackErr.Error())
		}
		return err
	}

	return errors.WithStack(t.transaction.Commit())
}
