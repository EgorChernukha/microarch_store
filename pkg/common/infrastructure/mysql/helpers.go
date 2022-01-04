package mysql

import (
	"database/sql/driver"
	uuid "github.com/satori/go.uuid"
)

type binaryUUID uuid.UUID

func newUUID() binaryUUID {
	return binaryUUID(uuid.NewV1())
}

func (uid binaryUUID) Value() (driver.Value, error) {
	return uuid.UUID(uid).Bytes(), nil
}

func (uid *binaryUUID) Scan(src interface{}) error {
	var result uuid.UUID
	err := result.Scan(src)
	*uid = binaryUUID(result)
	return err
}

type nullBinaryUUID uuid.NullUUID

func (uid nullBinaryUUID) Value() (driver.Value, error) {
	if !uid.Valid {
		return nil, nil
	}
	// Delegate to UUID Value function
	return uid.UUID.Bytes(), nil
}

func (uid *nullBinaryUUID) Scan(src interface{}) error {
	var result uuid.NullUUID
	err := result.Scan(src)
	*uid = nullBinaryUUID(result)
	return err
}
