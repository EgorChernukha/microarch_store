package mysql

import (
	"database/sql"
	"time"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"store/pkg/common/infrastructure/mysql"
	"store/pkg/delivery/app"
)

func NewOrderDeliveryQueryService(client mysql.Client) app.OrderDeliveryQueryService {
	return &orderDeliveryQueryService{client: client}
}

type sqlxOrderDeliveryData struct {
	ID        mysql.BinaryUUID `db:"id"`
	OrderID   mysql.BinaryUUID `db:"order_id"`
	UserID    mysql.BinaryUUID `db:"user_id"`
	Status    int              `db:"status"`
	UpdatedAt time.Time        `db:"updated_at"`
}

type orderDeliveryQueryService struct {
	client mysql.Client
}

func (d *orderDeliveryQueryService) FindByOrderID(orderID uuid.UUID) (app.OrderDeliveryData, error) {
	const sqlQuery = `SELECT id, order_id, user_id, status, updated_at FROM order_delivery WHERE order_id=?`

	var sqlOrderDeliveryData sqlxOrderDeliveryData

	err := d.client.Get(&sqlOrderDeliveryData, sqlQuery, orderID.Bytes())
	if err == sql.ErrNoRows {
		return app.OrderDeliveryData{}, errors.WithStack(app.ErrOrderDeliveryNotExists)
	} else if err != nil {
		return app.OrderDeliveryData{}, errors.WithStack(err)
	}

	return sqlxOrderDeliveryDataToOrderDeliveryData(&sqlOrderDeliveryData), nil
}

func (d *orderDeliveryQueryService) FindByUserID(userID uuid.UUID) ([]app.OrderDeliveryData, error) {
	const sqlQuery = `SELECT id, order_id, user_id, status, updated_at FROM order_delivery WHERE user_id=?`

	var sqlOrderDeliveryDataList []*sqlxOrderDeliveryData
	err := d.client.Select(&sqlOrderDeliveryDataList, sqlQuery, userID.Bytes())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res := make([]app.OrderDeliveryData, 0, len(sqlOrderDeliveryDataList))
	for _, orderDeliveryData := range sqlOrderDeliveryDataList {
		res = append(res, sqlxOrderDeliveryDataToOrderDeliveryData(orderDeliveryData))
	}
	return res, nil
}

func (d *orderDeliveryQueryService) ListOrderDelivery() ([]app.OrderDeliveryData, error) {
	const sqlQuery = `SELECT id, order_id, user_id, status, updated_at FROM order_delivery`

	var sqlOrderDeliveryDataList []*sqlxOrderDeliveryData
	err := d.client.Select(&sqlOrderDeliveryDataList, sqlQuery)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res := make([]app.OrderDeliveryData, 0, len(sqlOrderDeliveryDataList))
	for _, orderDeliveryData := range sqlOrderDeliveryDataList {
		res = append(res, sqlxOrderDeliveryDataToOrderDeliveryData(orderDeliveryData))
	}
	return res, nil
}

func sqlxOrderDeliveryDataToOrderDeliveryData(orderDelivery *sqlxOrderDeliveryData) app.OrderDeliveryData {
	return app.OrderDeliveryData{
		ID:        uuid.UUID(orderDelivery.ID),
		OrderID:   uuid.UUID(orderDelivery.OrderID),
		UserID:    uuid.UUID(orderDelivery.UserID),
		Status:    orderDelivery.Status,
		UpdatedAt: orderDelivery.UpdatedAt,
	}
}
