package app

import uuid "github.com/satori/go.uuid"

type ReserveOrderPositionInputItem struct {
	PositionID uuid.UUID
	OrderID    uuid.UUID
	Count      int
}

type ReserveOrderPositionInput struct {
	Positions []ReserveOrderPositionInputItem
}

type StockClient interface {
	ReserveOrderPositions(input ReserveOrderPositionInput) (succeeded bool, err error)
}
