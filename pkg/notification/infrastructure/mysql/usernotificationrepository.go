package mysql

import (
	"github.com/pkg/errors"

	"store/pkg/common/infrastructure/mysql"

	"store/pkg/notification/domain"
)

func NewUserNotificationRepository(client mysql.Client) domain.UserNotificationRepository {
	return &userNotificationRepository{client: client}
}

type userNotificationRepository struct {
	client mysql.Client
}

type sqlxUserNotification struct {
	UserID  mysql.BinaryUUID `db:"user_id"`
	OrderID mysql.BinaryUUID `db:"order_id"`
	Message string           `db:"message"`
}

func (r *userNotificationRepository) Store(userNotification domain.UserNotification) error {
	const sqlQuery = `INSERT INTO user_notification
	(user_id, order_id, message)
	VALUES (:user_id, :order_id, :message)`

	sqlUserNotification := sqlxUserNotification{
		UserID:  mysql.BinaryUUID(userNotification.UserID()),
		OrderID: mysql.BinaryUUID(userNotification.OrderID()),
		Message: userNotification.Message(),
	}

	_, err := r.client.NamedQuery(sqlQuery, &sqlUserNotification)
	return errors.WithStack(err)
}
