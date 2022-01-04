package mysql

import (
	"database/sql/driver"
	uuid "github.com/satori/go.uuid"
)

type BinaryUUID uuid.UUID

func NewUUID() BinaryUUID {
	return BinaryUUID(uuid.NewV1())
}

func (uid BinaryUUID) Value() (driver.Value, error) {
	return uuid.UUID(uid).Bytes(), nil
}

func (uid *BinaryUUID) Scan(src interface{}) error {
	var result uuid.UUID
	err := result.Scan(src)
	*uid = BinaryUUID(result)
	return err
}

type NullBinaryUUID uuid.NullUUID

func (uid NullBinaryUUID) Value() (driver.Value, error) {
	if !uid.Valid {
		return nil, nil
	}
	// Delegate to UUID Value function
	return uid.UUID.Bytes(), nil
}

func (uid *NullBinaryUUID) Scan(src interface{}) error {
	var result uuid.NullUUID
	err := result.Scan(src)
	*uid = NullBinaryUUID(result)
	return err
}
