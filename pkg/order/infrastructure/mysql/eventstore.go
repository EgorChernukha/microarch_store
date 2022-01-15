package mysql

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"store/pkg/common/app/integrationevent"
	"store/pkg/common/app/storedevent"
	"store/pkg/common/infrastructure/mysql"
)

func NewEventStore(client mysql.Client) storedevent.EventStore {
	return &eventStore{client: client}
}

type eventStore struct {
	client mysql.Client
}

func (store *eventStore) Add(event integrationevent.EventData) error {
	const query = `INSERT INTO stored_event (uid, type, body, confirmed) VALUES (:uid, :type, :body, :confirmed)`

	eventX := sqlxStoredEvent{
		UID:       mysql.BinaryUUID(event.UID),
		Type:      event.Type,
		Body:      event.Body,
		Confirmed: false,
	}

	_, err := store.client.NamedExec(query, &eventX)
	return errors.WithStack(err)
}

func (store *eventStore) ConfirmDelivery(id storedevent.EventID) error {
	const query = `UPDATE stored_event SET confirmed = TRUE WHERE id = ?`

	_, err := store.client.Exec(query, id)
	return errors.WithStack(err)
}

func (store *eventStore) FindByUIDs(uids []integrationevent.EventUID) ([]storedevent.Event, error) {
	const sqlQuery = `SELECT id, uid, type, body, confirmed FROM stored_event WHERE uid IN (?)`

	strUids := make([][]byte, 0, len(uids))
	for _, uid := range uids {
		strUids = append(strUids, uuid.UUID(uid).Bytes())
	}

	query, params, err := sqlx.In(sqlQuery, strUids)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	query = sqlx.Rebind(sqlx.DOLLAR, query)

	var events []*sqlxStoredEvent
	err = store.client.Select(&events, query, params...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res := make([]storedevent.Event, 0, len(events))
	for _, event := range events {
		res = append(res, sqlxStoredEventToEvent(event))
	}
	return res, nil
}

func (store *eventStore) FindAllUnconfirmedBefore(time time.Time) ([]storedevent.Event, error) {
	const sqlQuery = `SELECT id, uid, type, body, confirmed FROM stored_event WHERE confirmed = FALSE AND created_at < $1`

	var events []*sqlxStoredEvent
	err := store.client.Select(&events, sqlQuery, time)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res := make([]storedevent.Event, 0, len(events))
	for _, event := range events {
		res = append(res, sqlxStoredEventToEvent(event))
	}
	return res, nil
}

func sqlxStoredEventToEvent(event *sqlxStoredEvent) storedevent.Event {
	return storedevent.Event{
		EventData: integrationevent.EventData{
			UID:  integrationevent.EventUID(event.UID),
			Type: event.Type,
			Body: event.Body,
		},
		ID:        storedevent.EventID(event.ID),
		Confirmed: event.Confirmed,
	}
}

type sqlxStoredEvent struct {
	ID        uint64           `db:"id"`
	UID       mysql.BinaryUUID `db:"Uid"`
	Type      string           `db:"type"`
	Body      string           `db:"body"`
	Confirmed bool             `db:"confirmed"`
}
