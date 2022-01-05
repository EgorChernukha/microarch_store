package mysql

import (
	"database/sql"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"store/pkg/common/infrastructure/mysql"

	"store/pkg/billing/domain"
)

func NewUserAccountRepository(client mysql.Client) domain.UserAccountRepository {
	return &userAccountRepository{client: client}
}

type userAccountRepository struct {
	client mysql.Client
}

type sqlxUserAccount struct {
	ID      mysql.BinaryUUID `db:"id"`
	UserID  mysql.BinaryUUID `db:"user_id"`
	Balance float64          `db:"balance"`
}

func (r *userAccountRepository) NewID() domain.ID {
	return domain.ID(mysql.NewUUID())
}

func (r *userAccountRepository) Store(userAccount domain.UserAccount) error {
	const sqlQuery = `INSERT INTO user_account
	(id, user_id, balance, updated_at)
	VALUES (:id, :user_id, :balance, NOW())
	ON DUPLICATE KEY UPDATE balance=VALUES(balance)`

	sqlUserAccount := sqlxUserAccount{
		ID:      mysql.BinaryUUID(userAccount.ID()),
		UserID:  mysql.BinaryUUID(userAccount.UserID()),
		Balance: userAccount.Balance(),
	}

	_, err := r.client.NamedQuery(sqlQuery, &sqlUserAccount)
	return errors.WithStack(err)
}

func (r *userAccountRepository) FindOneByUserID(userID domain.UserID) (domain.UserAccount, error) {
	const sqlQuery = `SELECT id, user_id, balance FROM user_account WHERE user_id=?`

	var sqlUserAccount sqlxUserAccount

	err := r.client.Get(&sqlUserAccount, sqlQuery, uuid.UUID(userID).Bytes())
	if err == sql.ErrNoRows {
		return nil, errors.WithStack(domain.ErrUserAccountNotFound)
	} else if err != nil {
		return nil, errors.WithStack(err)
	}

	return domain.NewUserAccount(domain.ID(sqlUserAccount.ID), domain.UserID(sqlUserAccount.UserID), sqlUserAccount.Balance), nil
}
