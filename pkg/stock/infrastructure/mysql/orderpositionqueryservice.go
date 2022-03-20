package mysql

import (
	"database/sql"
	"github.com/pkg/errors"
	"time"

	uuid "github.com/satori/go.uuid"

	"store/pkg/common/infrastructure/mysql"
	"store/pkg/stock/app"
)

func NewOrderPositionQueryService(client mysql.Client) app.OrderPositionQueryService {
	return &orderPositionQueryService{client: client}
}

type sqlxOrderPositionData struct {
	ID         mysql.BinaryUUID `db:"id"`
	OrderID    mysql.BinaryUUID `db:"order_id"`
	PositionID mysql.BinaryUUID `db:"position_id"`
	Count      int              `db:"count"`
	UpdatedAt  time.Time        `db:"updated_at"`
}

type orderPositionQueryService struct {
	client mysql.Client
}

func (o *orderPositionQueryService) FindByOrderID(orderID uuid.UUID) ([]app.OrderPositionData, error) {
	const sqlQuery = `SELECT id, order_id, position_id, count, updated_at FROM order_position WHERE order_id=?`

	var orderPositions []*sqlxOrderPositionData
	err := o.client.Select(&orderPositions, sqlQuery, orderID.Bytes())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res := make([]app.OrderPositionData, 0, len(orderPositions))
	for _, orderPosition := range orderPositions {
		res = append(res, sqlxOrderPositionDataToPositionData(orderPosition))
	}
	return res, nil
}

func (o *orderPositionQueryService) FindByPositionID(positionID uuid.UUID) ([]app.OrderPositionData, error) {
	const sqlQuery = `SELECT id, order_id, position_id, count, updated_at FROM order_position WHERE position_id=?`

	var orderPositions []*sqlxOrderPositionData
	err := o.client.Select(&orderPositions, sqlQuery, positionID.Bytes())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res := make([]app.OrderPositionData, 0, len(orderPositions))
	for _, orderPosition := range orderPositions {
		res = append(res, sqlxOrderPositionDataToPositionData(orderPosition))
	}
	return res, nil
}

func (o *orderPositionQueryService) FindByOrderIDAndPositionID(orderID uuid.UUID, positionID uuid.UUID) (app.OrderPositionData, error) {
	const sqlQuery = `SELECT id, order_id, position_id, count, updated_at FROM order_position WHERE order_id=? AND position_id=?`

	var sqlOrderPosition sqlxOrderPositionData

	err := o.client.Get(&sqlOrderPosition, sqlQuery, orderID.Bytes(), positionID.Bytes())
	if err == sql.ErrNoRows {
		return app.OrderPositionData{}, errors.WithStack(app.ErrOrderPositionNotExists)
	} else if err != nil {
		return app.OrderPositionData{}, errors.WithStack(err)
	}

	return sqlxOrderPositionDataToPositionData(&sqlOrderPosition), nil
}

func (o *orderPositionQueryService) ListOrderPositions() ([]app.OrderPositionData, error) {
	const sqlQuery = `SELECT id, order_id, position_id, count, updated_at FROM order_position`

	var orderPositions []*sqlxOrderPositionData
	err := o.client.Select(&orderPositions, sqlQuery)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res := make([]app.OrderPositionData, 0, len(orderPositions))
	for _, orderPosition := range orderPositions {
		res = append(res, sqlxOrderPositionDataToPositionData(orderPosition))
	}
	return res, nil
}

func sqlxOrderPositionDataToPositionData(orderPositionData *sqlxOrderPositionData) app.OrderPositionData {
	return app.OrderPositionData{
		ID:         uuid.UUID(orderPositionData.ID),
		OrderID:    uuid.UUID(orderPositionData.OrderID),
		PositionID: uuid.UUID(orderPositionData.PositionID),
		Count:      orderPositionData.Count,
		UpdatedAt:  orderPositionData.UpdatedAt,
	}
}
