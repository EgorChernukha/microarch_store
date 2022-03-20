package mysql

import (
	"database/sql"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"store/pkg/common/infrastructure/mysql"
	"store/pkg/stock/app"
)

func NewPositionRepository(client mysql.Client) app.PositionRepository {
	return &positionRepository{client: client}
}

type sqlxPosition struct {
	ID    mysql.BinaryUUID `db:"id"`
	Total int              `db:"total"`
}

type positionRepository struct {
	client mysql.Client
}

func (p *positionRepository) NewID() app.PositionID {
	return app.PositionID(mysql.NewUUID())
}

func (p *positionRepository) Store(position app.Position) error {
	const sqlQuery = `INSERT INTO position
	(id, total, updated_at)
	VALUES (:id, :total, NOW())
	ON DUPLICATE KEY UPDATE total=VALUES(total)`

	sqlPosition := sqlxPosition{
		ID:    mysql.BinaryUUID(position.ID()),
		Total: position.Total(),
	}

	_, err := p.client.NamedQuery(sqlQuery, &sqlPosition)
	return errors.WithStack(err)
}

func (p *positionRepository) FindByID(id app.PositionID) (app.Position, error) {
	const sqlQuery = `SELECT id, total FROM position WHERE id=?`

	var sqlPosition sqlxPosition

	err := p.client.Get(&sqlPosition, sqlQuery, uuid.UUID(id).Bytes())
	if err == sql.ErrNoRows {
		return nil, errors.WithStack(app.ErrPositionNotFound)
	} else if err != nil {
		return nil, errors.WithStack(err)
	}

	return app.NewPosition(
		app.PositionID(sqlPosition.ID),
		sqlPosition.Total,
	), nil
}
