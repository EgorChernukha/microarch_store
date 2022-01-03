package mysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"store/pkg/user/app"
)

func NewUserQueryService(client Client) app.UserQueryService {
	return &userQueryService{client: client}
}

type userQueryService struct {
	client Client
}

func (u userQueryService) FindUser(id uuid.UUID) (app.UserData, error) {
	const sqlQuery = `SELECT id, login, firstname, lastname, email, phone FROM user WHERE id = ?`
	var sqlUser sqlxUser

	err := u.client.Get(&sqlUser, sqlQuery, id.Bytes())
	if err == sql.ErrNoRows {
		return app.UserData{}, app.ErrUserNotExists
	} else if err != nil {
		return app.UserData{}, errors.WithStack(err)
	}

	return app.UserData{
		ID:        uuid.UUID(sqlUser.ID),
		Username:  sqlUser.Login,
		Firstname: sqlUser.Firstname,
		Lastname:  sqlUser.Lastname,
		Email:     sqlUser.Email,
		Phone:     sqlUser.Phone,
	}, nil
}
