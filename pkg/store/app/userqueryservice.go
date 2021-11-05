package app

import uuid "github.com/satori/go.uuid"

type UserData struct {
	ID        uuid.UUID
	Username  string
	Firstname string
	Lastname  string
	Email     string
	Phone     string
}

type UserQueryService interface {
	FindUser(id uuid.UUID) (UserData, error)
}
