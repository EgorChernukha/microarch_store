package integrationevent

import uuid "github.com/satori/go.uuid"

type EventUID uuid.UUID

type EventData struct {
	UID  EventUID
	Type string
	Body string
}
