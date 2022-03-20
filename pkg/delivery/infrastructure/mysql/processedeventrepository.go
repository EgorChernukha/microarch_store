package mysql

import (
	"github.com/pkg/errors"

	"store/pkg/common/app/integrationevent"
	"store/pkg/common/infrastructure/mysql"

	"store/pkg/delivery/app"
)

func NewProcessedEventRepository(client mysql.Client) app.ProcessedEventRepository {
	return &processedEventRepository{client: client}
}

type processedEventRepository struct {
	client mysql.Client
}

func (repo *processedEventRepository) SetProcessed(eventID integrationevent.EventUID) (alreadyProcessed bool, err error) {
	const query = `INSERT IGNORE INTO processed_event (event_id) VALUES (?)`

	result, err := repo.client.Exec(query, mysql.BinaryUUID(eventID))
	if err != nil {
		return false, errors.WithStack(err)
	}

	rowsAffected, _ := result.RowsAffected()

	return rowsAffected < 1, nil
}
