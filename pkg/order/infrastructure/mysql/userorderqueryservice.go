package mysql

import (
	"database/sql"
	"time"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"store/pkg/common/infrastructure/mysql"

	"store/pkg/order/app"
)

func NewUserOrderQueryService(client mysql.Client) app.UserOrderQueryService {
	return &userOrderQueryService{client: client}
}

type userOrderQueryService struct {
	client mysql.Client
}

type sqlxUserOrderData struct {
	ID        mysql.BinaryUUID `db:"id"`
	UserID    mysql.BinaryUUID `db:"user_id"`
	OrderID   mysql.BinaryUUID `db:"order_id"`
	Price     float64          `db:"price"`
	Status    int              `db:"status"`
	CreatedAt time.Time        `db:"created_at"`
	UpdatedAt time.Time        `db:"updated_at"`
}

func (u *userOrderQueryService) FindUserOrderByOrderID(orderID uuid.UUID) (app.UserOrderData, error) {
	const sqlQuery = `SELECT id, user_id, order_id, price, status, created_at, updated_at FROM user_order where order_id = ?`

	var userOrder sqlxUserOrderData

	err := u.client.Get(&userOrder, sqlQuery, orderID.Bytes())
	if err == sql.ErrNoRows {
		return app.UserOrderData{}, app.ErrUserOrderNotExists
	} else if err != nil {
		return app.UserOrderData{}, errors.WithStack(err)
	}

	return sqlxUserOrderDataToUserOrderData(&userOrder), nil
}

func (u *userOrderQueryService) ListUserOrdersByUserIDs(userID uuid.UUID) ([]app.UserOrderData, error) {
	const sqlQuery = `SELECT id, user_id, order_id, price, status, created_at, updated_at FROM user_order where user_id = ?`

	var userOrders []*sqlxUserOrderData
	err := u.client.Select(&userOrders, sqlQuery, userID.Bytes())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res := make([]app.UserOrderData, 0, len(userOrders))
	for _, userOrder := range userOrders {
		res = append(res, sqlxUserOrderDataToUserOrderData(userOrder))
	}
	return res, nil
}

func (u *userOrderQueryService) ListUserOrders() ([]app.UserOrderData, error) {
	const sqlQuery = `SELECT id, user_id, order_id, price, status, created_at, updated_at FROM user_order`

	var userOrders []*sqlxUserOrderData
	err := u.client.Select(&userOrders, sqlQuery)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res := make([]app.UserOrderData, 0, len(userOrders))
	for _, userOrder := range userOrders {
		res = append(res, sqlxUserOrderDataToUserOrderData(userOrder))
	}
	return res, nil
}

func sqlxUserOrderDataToUserOrderData(userOrder *sqlxUserOrderData) app.UserOrderData {
	return app.UserOrderData{
		ID:        uuid.UUID(userOrder.ID),
		UserID:    uuid.UUID(userOrder.UserID),
		OrderID:   uuid.UUID(userOrder.OrderID),
		Price:     userOrder.Price,
		Status:    userOrder.Status,
		CreatedAt: userOrder.CreatedAt,
		UpdatedAt: userOrder.UpdatedAt,
	}
}
