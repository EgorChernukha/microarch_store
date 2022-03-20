package app

import (
	"errors"
	"time"

	uuid "github.com/satori/go.uuid"
)

type OrderPositionData struct {
	ID         uuid.UUID `json:"id"`
	OrderID    uuid.UUID `json:"order_id"`
	PositionID uuid.UUID `json:"position_id"`
	Count      int       `json:"count"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type OrderPositionQueryService interface {
	FindByOrderID(orderID uuid.UUID) ([]PositionData, error)
	FindByPositionID(positionID uuid.UUID) ([]PositionData, error)
	FindByOrderIDAndPositionID(orderID uuid.UUID, positionID uuid.UUID) (PositionData, error)
	ListOrderPositions() ([]PositionData, error)
}

var ErrOrderPositionNotExists = errors.New("order position does not exist")
