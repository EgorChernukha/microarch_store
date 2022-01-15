package app

import (
	"encoding/json"

	uuid "github.com/satori/go.uuid"

	"store/pkg/common/app/integrationevent"
)

const typeOrderConfirmed = "order.order_confirmed"
const typeOrderRejected = "order.order_rejected"

func NewUserOrderConfirmedEvent(orderID OrderID, userID UserID) integrationevent.EventData {
	body, _ := json.Marshal(orderEventBody{
		OrderID: uuid.UUID(orderID).String(),
		UserID:  uuid.UUID(userID).String(),
	})

	return integrationevent.EventData{
		UID:  newUID(),
		Type: typeOrderConfirmed,
		Body: string(body),
	}
}

func NewUserOrderRejectedEvent(orderID OrderID, userID UserID) integrationevent.EventData {
	body, _ := json.Marshal(orderEventBody{
		OrderID: uuid.UUID(orderID).String(),
		UserID:  uuid.UUID(userID).String(),
	})

	return integrationevent.EventData{
		UID:  newUID(),
		Type: typeOrderRejected,
		Body: string(body),
	}
}

func newUID() integrationevent.EventUID {
	return integrationevent.EventUID(uuid.NewV1())
}

type orderEventBody struct {
	OrderID string `json:"order_id"`
	UserID  string `json:"user_id"`
}
