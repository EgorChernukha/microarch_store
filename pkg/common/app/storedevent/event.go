package storedevent

import (
	"time"

	"store/pkg/common/app/integrationevent"
)

type EventID uint64

type Event struct {
	integrationevent.EventData
	ID        EventID
	Confirmed bool
}

type EventStore interface {
	Add(event integrationevent.EventData) error
	ConfirmDelivery(id EventID) error
	FindByUIDs(uids []integrationevent.EventUID) ([]Event, error)
	FindAllUnconfirmedBefore(time time.Time) ([]Event, error)
}
