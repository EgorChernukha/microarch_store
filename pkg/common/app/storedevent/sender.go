package storedevent

import "store/pkg/common/app/integrationevent"

type Sender interface {
	EventStored(uid integrationevent.EventUID)
	SendStoredEvents()
}
