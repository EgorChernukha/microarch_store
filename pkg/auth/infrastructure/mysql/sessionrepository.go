package mysql

import (
	"database/sql"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"store/pkg/auth/app"
	"time"
)

func NewSessionRepository(client Client) app.SessionRepository {
	return &sessionRepository{client: client}
}

type sessionRepository struct {
	client Client
}

type sqlxSession struct {
	ID        binaryUUID `db:"id"`
	UserID    binaryUUID `db:"user_id"`
	ValidTill time.Time  `db:"valid_till"`
}

func (s *sessionRepository) Store(session *app.Session) error {
	const sqlQuery = `INSERT INTO user_session
	(id, user_id, valid_till)
	VALUES (:id, :user_id, :valid_till)
	ON DUPLICATE KEY UPDATE user_id=VALUES(user_id), valid_till=VALUES(valid_till)`

	sqlSession := sqlxSession{
		ID:        binaryUUID(session.ID),
		UserID:    binaryUUID(session.UserID),
		ValidTill: session.ValidTill,
	}

	_, err := s.client.NamedQuery(sqlQuery, &sqlSession)
	return errors.WithStack(err)
}

func (s *sessionRepository) Remove(id app.SessionID) error {
	const sqlQuery = `DELETE FROM user_session WHERE id=?`
	_, err := s.client.Exec(sqlQuery, binaryUUID(id))

	return errors.WithStack(err)
}

func (s *sessionRepository) FindOneByID(id app.SessionID) (*app.Session, error) {
	const sqlQuery = `SELECT id, user_id, valid_till FROM user_session WHERE id=?`

	var sqlSession sqlxSession

	err := s.client.Get(&sqlSession, sqlQuery, uuid.UUID(id).Bytes())
	if err == sql.ErrNoRows {
		return nil, errors.WithStack(app.ErrUserNotFound)
	} else if err != nil {
		return nil, errors.WithStack(err)
	}

	return &app.Session{
		ID:        app.SessionID(sqlSession.ID),
		UserID:    app.UserID(sqlSession.UserID),
		ValidTill: sqlSession.ValidTill,
	}, nil
}
