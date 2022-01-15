package app

import uuid "github.com/satori/go.uuid"

type BillingClient interface {
	ProcessOrderPayment(userID uuid.UUID, price float64) (succeeded bool, err error)
}
