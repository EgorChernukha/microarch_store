package mysql

import (
	"database/sql"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"store/pkg/common/infrastructure/mysql"

	"store/pkg/auth/app"
)

func NewUserRepository(client mysql.Client) app.UserRepository {
	return &userRepository{client: client}
}

type userRepository struct {
	client mysql.Client
}

type sqlxUser struct {
	ID       mysql.BinaryUUID `db:"id"`
	Login    string           `db:"login"`
	Password string           `db:"password"`
}

func (u *userRepository) Store(user *app.User) error {
	const sqlQuery = `INSERT INTO user_auth
	(id, login, password)
	VALUES (:id, :login, :password)
	ON DUPLICATE KEY UPDATE login=VALUES(login), password=VALUES(password)`

	sqlUser := sqlxUser{
		ID:       mysql.BinaryUUID(user.ID),
		Login:    user.Login,
		Password: user.Password,
	}

	_, err := u.client.NamedQuery(sqlQuery, &sqlUser)
	return errors.WithStack(err)
}

func (u *userRepository) Remove(id app.UserID) error {
	const sqlQuery = `DELETE FROM user_auth WHERE id=?`
	_, err := u.client.Exec(sqlQuery, mysql.BinaryUUID(id))

	return errors.WithStack(err)
}

func (u *userRepository) FindOneByID(id app.UserID) (*app.User, error) {
	const sqlQuery = `SELECT id, login, password FROM user_auth WHERE id=?`

	var sqlUser sqlxUser

	err := u.client.Get(&sqlUser, sqlQuery, uuid.UUID(id).Bytes())
	if err == sql.ErrNoRows {
		return nil, errors.WithStack(app.ErrUserNotFound)
	} else if err != nil {
		return nil, errors.WithStack(err)
	}

	return &app.User{
		ID:       app.UserID(sqlUser.ID),
		Login:    sqlUser.Login,
		Password: sqlUser.Password,
	}, nil
}

func (u *userRepository) FindOneByLogin(login string) (*app.User, error) {
	const sqlQuery = `SELECT id, login, password FROM user_auth WHERE login=?`

	var sqlUser sqlxUser

	err := u.client.Get(&sqlUser, sqlQuery, login)
	if err == sql.ErrNoRows {
		return nil, errors.WithStack(app.ErrUserNotFound)
	} else if err != nil {
		return nil, errors.WithStack(err)
	}

	return &app.User{
		ID:       app.UserID(sqlUser.ID),
		Login:    sqlUser.Login,
		Password: sqlUser.Password,
	}, nil
}
