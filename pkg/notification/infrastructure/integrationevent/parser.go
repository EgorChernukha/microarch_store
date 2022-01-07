package integrationevent

import (
	"encoding/json"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"store/pkg/notification/app/integrationevent"
)

type eventBody struct {
	Type      string           `json:"type"`
	AccountID string           `json:"account_id"`
	Payload   *json.RawMessage `json:"payload"`
}

type userRegisteredMessage struct {
	UserID       string `json:"user_id"`
	DepartmentID string `json:"department_id"`
}

const (
	TypeUserRegistered = "user.user_registered"
)

func (h *handler) parseClientMessage(msg string) (integrationevent.IntegrationEvent, error) {
	var body eventBody
	err := json.Unmarshal([]byte(msg), &body)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	switch t := body.Type; t {
	case TypeUserRegistered:
		return parseUserRegisteredMessage(body)
	default:
		return nil, nil
	}
}

func parseUserRegisteredMessage(body eventBody) (integrationevent.IntegrationEvent, error) {
	if body.Payload == nil {
		return nil, errors.New("failed to parse user registered message: no payload")
	}
	var msg userRegisteredMessage
	err := json.Unmarshal(*body.Payload, &msg)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	userID, err := uuid.FromString(msg.UserID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	departmentID, err := uuid.FromString(msg.DepartmentID)
	return integrationevent.NewUserRegistered(body.AccountID, userID, departmentID), err
}
