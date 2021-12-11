package app

import (
	"errors"
	"time"

	uuid "github.com/satori/go.uuid"
)

type SessionID uuid.UUID

type Session struct {
	ID        SessionID
	UserID    UserID
	ValidTill time.Time
}

type SessionRepository interface {
	Store(session *Session) error
	Remove(id SessionID) error
	FindOneByID(id SessionID) (*Session, error)
}

var ErrSessionNotFound = errors.New("session not found")
