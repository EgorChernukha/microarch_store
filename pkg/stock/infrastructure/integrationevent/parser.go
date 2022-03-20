package integrationevent

import (
	"encoding/json"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"store/pkg/common/app/integrationevent"

	"store/pkg/stock/app"
)

const typeOrderConfirmed = "order.order_confirmed"

func NewEventParser() app.IntegrationEventParser {
	return eventParser{}
}

type eventParser struct {
}

func (e eventParser) ParseIntegrationEvent(event integrationevent.EventData) (app.UserEvent, error) {
	switch event.Type {
	case typeOrderConfirmed:
		return parseOrderConfirmedEvent(event.Body)
	default:
		return nil, nil
	}
}

func parseOrderConfirmedEvent(strBody string) (app.UserEvent, error) {
	body, err := parseOrderEvent(strBody)
	if err != nil {
		return nil, err
	}

	userID, err := uuid.FromString(body.UserID)
	if err != nil {
		return nil, err
	}

	orderID, err := uuid.FromString(body.OrderID)
	if err != nil {
		return nil, err
	}

	return app.NewOrderConfirmedEvent(userID, orderID), nil
}

func parseOrderEvent(strBody string) (orderEventBody, error) {
	var body orderEventBody
	err := json.Unmarshal([]byte(strBody), &body)
	if err != nil {
		return body, errors.WithStack(err)
	}
	_, err = uuid.FromString(body.UserID)
	if err != nil {
		return body, errors.WithStack(err)
	}
	_, err = uuid.FromString(body.OrderID)
	return body, errors.WithStack(err)
}

type orderEventBody struct {
	OrderID string `json:"order_id"`
	UserID  string `json:"user_id"`
}
