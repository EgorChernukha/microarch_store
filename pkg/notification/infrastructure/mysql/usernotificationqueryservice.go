package mysql

import (
	"time"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"store/pkg/common/infrastructure/mysql"
	"store/pkg/notification/app"
)

func NewUserNotificationQueryService(client mysql.Client) app.UserNotificationQueryService {
	return &userNotificationQueryService{client: client}
}

type userNotificationQueryService struct {
	client mysql.Client
}

type sqlxUserNotificationData struct {
	UserID    mysql.BinaryUUID `db:"user_id"`
	OrderID   mysql.BinaryUUID `db:"order_id"`
	Email     string           `db:"email"`
	Message   string           `db:"message"`
	CreatedAt time.Time        `db:"created_at"`
}

func (s *userNotificationQueryService) ListUserNotifications(userID uuid.UUID) ([]app.UserNotificationData, error) {
	const sqlQuery = `SELECT user_id, order_id, email, message, created_at FROM user_notification WHERE user_id = ?`

	var userNotifications []*sqlxUserNotificationData
	err := s.client.Select(&userNotifications, sqlQuery, userID.Bytes())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res := make([]app.UserNotificationData, 0, len(userNotifications))
	for _, userNotification := range userNotifications {
		res = append(res, sqlxUserNotificationDataToUserNotificationData(userNotification))
	}
	return res, nil
}

func sqlxUserNotificationDataToUserNotificationData(userNotification *sqlxUserNotificationData) app.UserNotificationData {
	return app.UserNotificationData{
		UserID:    uuid.UUID(userNotification.UserID),
		OrderID:   uuid.UUID(userNotification.OrderID),
		Email:     userNotification.Email,
		Message:   userNotification.Message,
		CreatedAt: userNotification.CreatedAt,
	}
}
