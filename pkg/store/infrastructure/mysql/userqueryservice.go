package mysql

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"mod/pkg/store/app"
)
import _ "github.com/go-sql-driver/mysql"

func NewUserQueryService(client *sqlx.DB) app.UserQueryService {
	return &userQueryService{client: client}
}

type userQueryService struct {
	client *sqlx.DB
}

func (u userQueryService) FindUser(id uuid.UUID) (app.UserData, error) {
	const sqlQuery = `SELECT id, login, firstname, lastname, email, phone FROM user WHERE id = ?`
	var sqlUser sqlxUser

	err := u.client.Get(&sqlUser, sqlQuery, id)
	if err != nil {
		return app.UserData{}, errors.WithStack(err)
	}

	return app.UserData{
		ID:        uuid.UUID(sqlUser.id),
		Username:  sqlUser.Login,
		Firstname: sqlUser.Firstname,
		Lastname:  sqlUser.Lastname,
		Email:     sqlUser.Email,
		Phone:     sqlUser.Phone,
	}, nil
}
