package mysql

import (
	"database/sql"
	"store/pkg/common/infrastructure/mysql"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"store/pkg/user/domain"
)

func NewUserRepository(client mysql.Client) domain.UserRepository {
	return &userRepository{client: client}
}

type userRepository struct {
	client mysql.Client
}

type sqlxUser struct {
	ID        mysql.BinaryUUID `db:"id"`
	Login     string           `db:"login"`
	Firstname string           `db:"firstname"`
	Lastname  string           `db:"lastname"`
	Email     string           `db:"email"`
	Phone     string           `db:"phone"`
}

func (u *userRepository) NewID() domain.UserID {
	return domain.UserID(mysql.NewUUID())
}

func (u *userRepository) Store(user domain.User) error {
	const sqlQuery = `INSERT INTO user
	(id, login, firstname, lastname, email, phone)
	VALUES (:id, :login, :firstname, :lastname, :email, :phone)
	ON DUPLICATE KEY UPDATE firstname=VALUES(firstname), lastname=VALUES(lastname), email=VALUES(email), phone=VALUES(phone)`

	sqlUser := sqlxUser{
		ID:        mysql.BinaryUUID(user.ID()),
		Login:     user.Login(),
		Firstname: user.Firstname(),
		Lastname:  user.Lastname(),
		Email:     user.Email(),
		Phone:     user.Phone(),
	}

	_, err := u.client.NamedQuery(sqlQuery, &sqlUser)
	return errors.WithStack(err)
}

func (u *userRepository) Remove(user domain.User) error {
	const sqlQuery = `DELETE FROM user WHERE id=?`
	_, err := u.client.Exec(sqlQuery, mysql.BinaryUUID(user.ID()))

	return errors.WithStack(err)
}

func (u *userRepository) FindOne(id domain.UserID) (domain.User, error) {
	const sqlQuery = `SELECT id, login, firstname, lastname, email, phone FROM user WHERE id=?`

	var sqlUser sqlxUser

	err := u.client.Get(&sqlUser, sqlQuery, uuid.UUID(id).Bytes())
	if err == sql.ErrNoRows {
		return nil, errors.WithStack(domain.ErrUserNotFound)
	} else if err != nil {
		return nil, errors.WithStack(err)
	}

	return domain.NewUser(domain.UserID(sqlUser.ID), sqlUser.Login, sqlUser.Firstname, sqlUser.Lastname, sqlUser.Email, sqlUser.Phone), nil
}
