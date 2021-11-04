package mysql

import (
	"database/sql"
	"github.com/pkg/errors"

	"mod/pkg/store/app"
)
import _ "github.com/go-sql-driver/mysql"

func newUserQueryService(client *sql.DB) app.UserQueryService {
	return &userQueryService{client: client}
}

type userQueryService struct {
	client *sql.DB
}

func (u userQueryService) FindUser(id int) (app.UserData, error) {
	const sqlQuery = `SELECT id, login, firstname, lastname, email, phone FROM user WHERE id = ?`

	userData := app.UserData{}

	err := u.client.QueryRow(sqlQuery, id).Scan(&userData.ID, &userData.Username, &userData.Firstname, &userData.Lastname, &userData.Email, &userData.Phone)
	if err != nil {
		return app.UserData{}, errors.WithStack(err)
	}

	return userData, nil
}
