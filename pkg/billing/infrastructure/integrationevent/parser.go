package integrationevent

import (
	"encoding/json"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"store/pkg/billing/app/integrationevent"
)

type eventBody struct {
	Type    string           `json:"type"`
	Payload *json.RawMessage `json:"payload"`
}

type userOrderCreated struct {
	UserID  string  `json:"user_id"`
	OrderID string  `json:"order_id"`
	Price   float64 `json:"price"`
}

type userOrderCanceled struct {
	UserID  string  `json:"user_id"`
	OrderID string  `json:"order_id"`
	Price   float64 `json:"price"`
}

const (
	TypeUserOrderCreated  = "user_order.created"
	TypeUserOrderCanceled = "user_order.canceled"
)

func (h *handler) parseClientMessage(msg string) (integrationevent.IntegrationEvent, error) {
	var body eventBody
	err := json.Unmarshal([]byte(msg), &body)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	switch t := body.Type; t {
	case TypeUserOrderCreated:
		return parseUserOrderCreated(body)
	case TypeUserOrderCanceled:
		return parseUserOrderCanceled(body)
	default:
		return nil, nil
	}
}

func parseUserOrderCreated(body eventBody) (integrationevent.IntegrationEvent, error) {
	if body.Payload == nil {
		return nil, errors.New("failed to parse user registered message: no payload")
	}
	var msg userOrderCreated
	err := json.Unmarshal(*body.Payload, &msg)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	userID, err := uuid.FromString(msg.UserID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	orderID, err := uuid.FromString(msg.OrderID)
	return integrationevent.NewUserOrderCreated(userID, orderID, msg.Price), err
}

func parseUserOrderCanceled(body eventBody) (integrationevent.IntegrationEvent, error) {
	if body.Payload == nil {
		return nil, errors.New("failed to parse user registered message: no payload")
	}
	var msg userOrderCanceled
	err := json.Unmarshal(*body.Payload, &msg)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	userID, err := uuid.FromString(msg.UserID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	orderID, err := uuid.FromString(msg.OrderID)
	return integrationevent.NewUserOrderCanceled(userID, orderID, msg.Price), err
}
