package integrationevent

import uuid "github.com/satori/go.uuid"

type EventID uuid.UUID

type EventData struct {
	UID  EventID
	Type string
	Body string
}
