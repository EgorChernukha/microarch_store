package app

import (
	"errors"
	"time"

	uuid "github.com/satori/go.uuid"
)

type PositionData struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Total     int       `json:"total"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PositionQueryService interface {
	FindPositionByID(positionID uuid.UUID) (PositionData, error)
	ListPositions() ([]PositionData, error)
}

var ErrPositionNotExists = errors.New("position does not exist")
