package app

import (
	"encoding/json"

	uuid "github.com/satori/go.uuid"

	"store/pkg/common/app/integrationevent"
)

const typeUserRegistered = "auth.user_registered"

func NewUserRegisteredEvent(userID UserID, login string) integrationevent.EventData {
	body, _ := json.Marshal(userRegisteredEventBody{
		UserID: uuid.UUID(userID).String(),
		Login:  login,
	})

	return integrationevent.EventData{
		UID:  newUID(),
		Type: typeUserRegistered,
		Body: string(body),
	}
}

func newUID() integrationevent.EventUID {
	return integrationevent.EventUID(uuid.NewV1())
}

type userRegisteredEventBody struct {
	UserID string `json:"user_id"`
	Login  string `json:"login"`
}
