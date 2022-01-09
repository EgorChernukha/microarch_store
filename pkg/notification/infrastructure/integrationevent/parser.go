package integrationevent

import (
	"encoding/json"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"store/pkg/notification/app/integrationevent"
)

type eventBody struct {
	Type    string           `json:"type"`
	Payload *json.RawMessage `json:"payload"`
}

type userOrderPaymentSucceeded struct {
	UserID  string `json:"user_id"`
	OrderID string `json:"order_id"`
	Email   string `json:"email"`
}

type userOrderPaymentFailed struct {
	UserID  string `json:"user_id"`
	OrderID string `json:"order_id"`
	Email   string `json:"email"`
}

const (
	TypeUserOrderPaymentSucceeded = "user_order.payment_succeeded"
	TypeUserOrderPaymentFailed    = "user_order.payment_failed"
)

func (h *handler) parseClientMessage(msg string) (integrationevent.IntegrationEvent, error) {
	var body eventBody
	err := json.Unmarshal([]byte(msg), &body)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	switch t := body.Type; t {
	case TypeUserOrderPaymentSucceeded:
		return parseUserOrderPaymentSucceeded(body)
	case TypeUserOrderPaymentFailed:
		return parseUserOrderPaymentFailed(body)
	default:
		return nil, nil
	}
}

func parseUserOrderPaymentSucceeded(body eventBody) (integrationevent.IntegrationEvent, error) {
	if body.Payload == nil {
		return nil, errors.New("failed to parse user registered message: no payload")
	}
	var msg userOrderPaymentSucceeded
	err := json.Unmarshal(*body.Payload, &msg)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	userID, err := uuid.FromString(msg.UserID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	orderID, err := uuid.FromString(msg.OrderID)
	return integrationevent.NewUserOrderPaymentSucceeded(userID, orderID, msg.Email), err
}

func parseUserOrderPaymentFailed(body eventBody) (integrationevent.IntegrationEvent, error) {
	if body.Payload == nil {
		return nil, errors.New("failed to parse user registered message: no payload")
	}
	var msg userOrderPaymentFailed
	err := json.Unmarshal(*body.Payload, &msg)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	userID, err := uuid.FromString(msg.UserID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	orderID, err := uuid.FromString(msg.OrderID)
	return integrationevent.NewUserOrderPaymentFailed(userID, orderID, msg.Email), err
}
