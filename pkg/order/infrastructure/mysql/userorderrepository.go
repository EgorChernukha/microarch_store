package mysql

import (
	"database/sql"
	"store/pkg/order/app"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"store/pkg/common/infrastructure/mysql"
)

func NewUserOrderRepository(client mysql.Client) app.UserOrderRepository {
	return &userOrderRepository{client: client}
}

type userOrderRepository struct {
	client mysql.Client
}

type sqlxUserOrder struct {
	ID      mysql.BinaryUUID `db:"id"`
	UserID  mysql.BinaryUUID `db:"user_id"`
	OrderID mysql.BinaryUUID `db:"order_id"`
	Price   float64          `db:"price"`
	Status  int              `db:"status"`
}

func (r *userOrderRepository) NewID() app.ID {
	return app.ID(mysql.NewUUID())
}

func (r *userOrderRepository) Store(userOrder app.UserOrder) error {
	const sqlQuery = `INSERT INTO user_order
	(id, user_id, order_id, price, status, updated_at)
	VALUES (:id, :user_id, :order_id, :price, :status, NOW())
	ON DUPLICATE KEY UPDATE price=VALUES(price), status=VALUES(status)`

	sqlUserOrder := sqlxUserOrder{
		ID:      mysql.BinaryUUID(userOrder.ID()),
		UserID:  mysql.BinaryUUID(userOrder.UserID()),
		OrderID: mysql.BinaryUUID(userOrder.OrderID()),
		Price:   userOrder.Price(),
		Status:  int(userOrder.Status()),
	}

	_, err := r.client.NamedQuery(sqlQuery, &sqlUserOrder)
	return errors.WithStack(err)
}

func (r *userOrderRepository) FindOneByOrderID(orderID app.OrderID) (app.UserOrder, error) {
	const sqlQuery = `SELECT id, user_id, order_id, price, status FROM user_order WHERE order_id=?`

	var sqlUserOrder sqlxUserOrder

	err := r.client.Get(&sqlUserOrder, sqlQuery, uuid.UUID(orderID).Bytes())
	if err == sql.ErrNoRows {
		return nil, errors.WithStack(app.ErrUserOrderNotFound)
	} else if err != nil {
		return nil, errors.WithStack(err)
	}

	return app.NewUserOrder(
		app.ID(sqlUserOrder.ID),
		app.UserID(sqlUserOrder.UserID),
		app.OrderID(sqlUserOrder.OrderID),
		sqlUserOrder.Price,
		app.Status(sqlUserOrder.Status),
	), nil
}
