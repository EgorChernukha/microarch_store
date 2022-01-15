package integrationevent

import (
	"encoding/json"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"store/pkg/billing/app"
	"store/pkg/common/app/integrationevent"
)

const typeUserRegistered = "auth.user_registered"

func NewEventParser() app.IntegrationEventParser {
	return eventParser{}
}

type eventParser struct {
}

func (e eventParser) ParseIntegrationEvent(event integrationevent.EventData) (app.UserEvent, error) {
	switch event.Type {
	case typeUserRegistered:
		return parseUserRegisteredEvent(event.Body)
	default:
		return nil, nil
	}
}

func parseUserRegisteredEvent(strBody string) (app.UserEvent, error) {
	var body userRegisteredEventBody
	err := json.Unmarshal([]byte(strBody), &body)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	userID, err := uuid.FromString(body.UserID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return app.NewUserRegisteredEvent(userID, body.Login), nil
}

type userRegisteredEventBody struct {
	UserID string `json:"user_id"`
	Login  string `json:"login"`
}
