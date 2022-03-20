package mysql

import (
	"database/sql"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"store/pkg/common/infrastructure/mysql"
	"store/pkg/stock/app"
)

func NewOrderPositionRepository(client mysql.Client) app.OrderPositionRepository {
	return &orderPositionRepository{client: client}
}

type sqlxOrderPosition struct {
	ID         mysql.BinaryUUID `db:"id"`
	OrderID    mysql.BinaryUUID `db:"order_id"`
	PositionID mysql.BinaryUUID `db:"position_id"`
	Count      int              `db:"count"`
}

type orderPositionRepository struct {
	client mysql.Client
}

func (o *orderPositionRepository) NewID() app.OrderPositionID {
	return app.OrderPositionID(mysql.NewUUID())
}

func (o *orderPositionRepository) Store(orderPosition app.OrderPosition) error {
	const sqlQuery = `INSERT INTO order_position
	(id, order_id, position_id, count, updated_at)
	VALUES (:id, :order_id, :position_id, :count, NOW())
	ON DUPLICATE KEY UPDATE count=VALUES(count)`

	sqlOrderPosition := sqlxOrderPosition{
		ID:         mysql.BinaryUUID(orderPosition.ID()),
		OrderID:    mysql.BinaryUUID(orderPosition.OrderID()),
		PositionID: mysql.BinaryUUID(orderPosition.PositionID()),
		Count:      orderPosition.Count(),
	}

	_, err := o.client.NamedQuery(sqlQuery, &sqlOrderPosition)
	return errors.WithStack(err)
}

func (o *orderPositionRepository) FindByOrderID(orderID app.OrderID) ([]app.OrderPosition, error) {
	const sqlQuery = `SELECT id, order_id, position_id, count FROM order_position WHERE AND order_id=?`

	var orderPositions []*sqlxOrderPosition
	err := o.client.Select(&orderPositions, sqlQuery, uuid.UUID(orderID).Bytes())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res := make([]app.OrderPosition, 0, len(orderPositions))
	for _, orderPosition := range orderPositions {
		res = append(res, sqlxOrderPositionToOrderPosition(orderPosition))
	}
	return res, nil
}

func (o *orderPositionRepository) FindByPositionID(positionID app.PositionID) ([]app.OrderPosition, error) {
	const sqlQuery = `SELECT id, order_id, position_id, count FROM order_position WHERE AND position_id=?`

	var orderPositions []*sqlxOrderPosition
	err := o.client.Select(&orderPositions, sqlQuery, uuid.UUID(positionID).Bytes())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res := make([]app.OrderPosition, 0, len(orderPositions))
	for _, orderPosition := range orderPositions {
		res = append(res, sqlxOrderPositionToOrderPosition(orderPosition))
	}
	return res, nil
}

func (o *orderPositionRepository) FindByOrderIDAndPositionID(orderID app.OrderID, positionID app.PositionID) (app.OrderPosition, error) {
	const sqlQuery = `SELECT id, order_id, position_id, count FROM order_position WHERE order_id=? AND position_id=?`

	var sqlOrderPosition sqlxOrderPosition

	err := o.client.Get(&sqlOrderPosition, sqlQuery, uuid.UUID(orderID).Bytes(), uuid.UUID(positionID).Bytes())
	if err == sql.ErrNoRows {
		return nil, errors.WithStack(app.ErrOrderPositionNotFound)
	} else if err != nil {
		return nil, errors.WithStack(err)
	}

	return sqlxOrderPositionToOrderPosition(&sqlOrderPosition), nil
}

func sqlxOrderPositionToOrderPosition(orderPosition *sqlxOrderPosition) app.OrderPosition {
	return app.NewOrderPosition(
		app.OrderPositionID(orderPosition.ID),
		app.OrderID(orderPosition.OrderID),
		app.PositionID(orderPosition.PositionID),
		orderPosition.Count,
	)
}
