package app

import "errors"

type PositionRepository interface {
	NewID() PositionID
	Store(position Position) error
	FindByID(id PositionID) (Position, error)
}

var ErrPositionNotFound = errors.New("position not found")
