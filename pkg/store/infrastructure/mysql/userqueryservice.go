package mysql

import (
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"store/pkg/store/app"
)
import _ "github.com/go-sql-driver/mysql"

func NewUserQueryService(client Client) app.UserQueryService {
	return &userQueryService{client: client}
}

type userQueryService struct {
	client Client
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
