package mysql

import (
	"database/sql"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"store/pkg/common/infrastructure/mysql"
	"store/pkg/delivery/app"
)

func NewOrderDeliveryRepository(client mysql.Client) app.OrderDeliveryRepository {
	return &orderDeliveryRepository{client: client}
}

type sqlxOrderDelivery struct {
	ID      mysql.BinaryUUID `db:"id"`
	OrderID mysql.BinaryUUID `db:"order_id"`
	UserID  mysql.BinaryUUID `db:"user_id"`
	Status  int              `db:"status"`
}

type orderDeliveryRepository struct {
	client mysql.Client
}

func (o *orderDeliveryRepository) NewID() app.ID {
	return app.ID(mysql.NewUUID())
}

func (o *orderDeliveryRepository) Store(orderDelivery app.OrderDelivery) error {
	const sqlQuery = `INSERT INTO order_delivery
	(id, order_id, user_id, status, updated_at)
	VALUES (:id, :order_id, :user_id, :status, NOW())
	ON DUPLICATE KEY UPDATE status=VALUES(status)`

	sqlOrderDelivery := sqlxOrderDelivery{
		ID:      mysql.BinaryUUID(orderDelivery.ID()),
		OrderID: mysql.BinaryUUID(orderDelivery.OrderID()),
		UserID:  mysql.BinaryUUID(orderDelivery.UserID()),
		Status:  orderDelivery.Status(),
	}

	_, err := o.client.NamedQuery(sqlQuery, &sqlOrderDelivery)
	return errors.WithStack(err)
}

func (o *orderDeliveryRepository) FindByID(id app.ID) (app.OrderDelivery, error) {
	const sqlQuery = `SELECT id, order_id, user_id, status FROM order_delivery WHERE id=?`

	var sqlOrderDelivery sqlxOrderDelivery

	err := o.client.Get(&sqlOrderDelivery, sqlQuery, uuid.UUID(id).Bytes())
	if err == sql.ErrNoRows {
		return nil, errors.WithStack(app.ErrOrderDeliveryNotFound)
	} else if err != nil {
		return nil, errors.WithStack(err)
	}

	return sqlxOrderDeliveryToOrderDelivery(&sqlOrderDelivery), nil
}

func (o *orderDeliveryRepository) FindByOrderID(orderID app.OrderID) (app.OrderDelivery, error) {
	const sqlQuery = `SELECT id, order_id, user_id, status FROM order_delivery WHERE order_id=?`

	var sqlOrderDelivery sqlxOrderDelivery

	err := o.client.Get(&sqlOrderDelivery, sqlQuery, uuid.UUID(orderID).Bytes())
	if err == sql.ErrNoRows {
		return nil, errors.WithStack(app.ErrOrderDeliveryNotFound)
	} else if err != nil {
		return nil, errors.WithStack(err)
	}

	return sqlxOrderDeliveryToOrderDelivery(&sqlOrderDelivery), nil
}

func (o *orderDeliveryRepository) FindByUserID(userID app.UserID) ([]app.OrderDelivery, error) {
	const sqlQuery = `SELECT id, order_id, user_id, status FROM order_delivery WHERE user_id=?`

	var sqlOrderDeliveries []*sqlxOrderDelivery
	err := o.client.Select(&sqlOrderDeliveries, sqlQuery, uuid.UUID(userID).Bytes())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res := make([]app.OrderDelivery, 0, len(sqlOrderDeliveries))
	for _, sqlOrderDelivery := range sqlOrderDeliveries {
		res = append(res, sqlxOrderDeliveryToOrderDelivery(sqlOrderDelivery))
	}
	return res, nil
}

func sqlxOrderDeliveryToOrderDelivery(orderDelivery *sqlxOrderDelivery) app.OrderDelivery {
	return app.NewOrderDelivery(
		app.ID(orderDelivery.ID),
		app.OrderID(orderDelivery.OrderID),
		app.UserID(orderDelivery.UserID),
		orderDelivery.Status,
	)
}
