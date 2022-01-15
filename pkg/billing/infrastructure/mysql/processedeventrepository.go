package mysql

import (
	"database/sql"

	"github.com/pkg/errors"

	"store/pkg/common/app/integrationevent"
	"store/pkg/common/infrastructure/mysql"
	"store/pkg/notification/app"
)

func NewProcessedEventRepository(client mysql.Client) app.ProcessedEventRepository {
	return &processedEventRepository{client: client}
}

type processedEventRepository struct {
	client mysql.Client
}

func (repo *processedEventRepository) SetProcessed(eventID integrationevent.EventID) (alreadyProcessed bool, err error) {
	const query = `INSERT IGNORE INTO processed_event (event_id) VALUES (:event_id)`

	var resUID string
	err = repo.client.Get(&resUID, query, mysql.BinaryUUID(eventID))
	if err != nil {
		if err == sql.ErrNoRows {
			return true, nil
		}
		return false, errors.WithStack(err)
	}
	return false, nil
}
