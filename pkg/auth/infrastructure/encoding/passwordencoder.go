package encoding

import (
	"crypto/sha256"
	"fmt"
	uuid "github.com/satori/go.uuid"

	"store/pkg/auth/app"
)

func NewPasswordEncoder() app.PasswordEncoder {
	return sha256PasswordEncoder{}
}

type sha256PasswordEncoder struct {
}

func (m sha256PasswordEncoder) Encode(password string, userID app.UserID) string {
	data := []byte(uuid.UUID(userID).String() + password)
	return fmt.Sprintf("%x", sha256.Sum256(data))
}
