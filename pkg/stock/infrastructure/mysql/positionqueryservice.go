package mysql

import (
	"database/sql"
	"time"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"store/pkg/common/infrastructure/mysql"
	"store/pkg/stock/app"
)

func NewPositionQueryService(client mysql.Client) app.PositionQueryService {
	return &positionQueryService{client: client}
}

type sqlxPositionData struct {
	ID        mysql.BinaryUUID `db:"id"`
	Total     int              `db:"total"`
	UpdatedAt time.Time        `db:"updated_at"`
}

type positionQueryService struct {
	client mysql.Client
}

func (p *positionQueryService) FindPositionByID(positionID uuid.UUID) (app.PositionData, error) {
	const sqlQuery = `SELECT id, total, updated_at FROM position WHERE id=?`

	var sqlPosition sqlxPositionData

	err := p.client.Get(&sqlPosition, sqlQuery, positionID.Bytes())
	if err == sql.ErrNoRows {
		return app.PositionData{}, errors.WithStack(app.ErrPositionNotExists)
	} else if err != nil {
		return app.PositionData{}, errors.WithStack(err)
	}

	return sqlxPositionDataToPositionData(&sqlPosition), nil
}

func (p *positionQueryService) ListPositions() ([]app.PositionData, error) {
	const sqlQuery = `SELECT id, total, updated_at FROM position`

	var positionDataList []*sqlxPositionData
	err := p.client.Select(&positionDataList, sqlQuery)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res := make([]app.PositionData, 0, len(positionDataList))
	for _, positionData := range positionDataList {
		res = append(res, sqlxPositionDataToPositionData(positionData))
	}
	return res, nil
}

func sqlxPositionDataToPositionData(positionData *sqlxPositionData) app.PositionData {
	return app.PositionData{
		ID:        uuid.UUID(positionData.ID),
		Total:     positionData.Total,
		UpdatedAt: positionData.UpdatedAt,
	}
}
