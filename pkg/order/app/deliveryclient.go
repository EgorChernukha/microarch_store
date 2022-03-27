package app

import uuid "github.com/satori/go.uuid"

type DeliveryClient interface {
	ReserveDelivery(userID uuid.UUID, orderID uuid.UUID) (succeeded bool, err error)
}
