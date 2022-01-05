package mysql

import (
	"database/sql"
	"time"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"store/pkg/common/infrastructure/mysql"

	"store/pkg/billing/app"
)

func NewUserAccountQueryService(client mysql.Client) app.UserAccountQueryService {
	return &userAccountQueryService{client: client}
}

type userAccountQueryService struct {
	client mysql.Client
}

type sqlxUserAccountData struct {
	ID        mysql.BinaryUUID `db:"id"`
	UserID    mysql.BinaryUUID `db:"user_id"`
	Balance   float64          `db:"balance"`
	UpdatedAt time.Time        `db:"updated_at"`
}

func (u *userAccountQueryService) FindUserAccountByUserID(userID uuid.UUID) (app.UserAccountData, error) {
	const sqlQuery = `SELECT id, user_id, balance, updated_at FROM user_account WHERE user_id = ?`
	var sqlUserAccountData sqlxUserAccountData

	err := u.client.Get(&sqlUserAccountData, sqlQuery, userID.Bytes())
	if err == sql.ErrNoRows {
		return app.UserAccountData{}, app.ErrUserAccountNotExists
	} else if err != nil {
		return app.UserAccountData{}, errors.WithStack(err)
	}

	return sqlxUserAccountDataToUserAccountData(&sqlUserAccountData), nil
}

func (u *userAccountQueryService) ListUserAccounts() ([]app.UserAccountData, error) {
	const sqlQuery = `SELECT id, user_id, balance, updated_at FROM user_account`

	var userAccounts []*sqlxUserAccountData
	err := u.client.Select(&userAccounts, sqlQuery)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res := make([]app.UserAccountData, 0, len(userAccounts))
	for _, userAccount := range userAccounts {
		res = append(res, sqlxUserAccountDataToUserAccountData(userAccount))
	}
	return res, nil
}

func sqlxUserAccountDataToUserAccountData(userAccount *sqlxUserAccountData) app.UserAccountData {
	return app.UserAccountData{
		ID:        uuid.UUID(userAccount.ID),
		UserID:    uuid.UUID(userAccount.UserID),
		Balance:   userAccount.Balance,
		UpdatedAt: userAccount.UpdatedAt,
	}
}
