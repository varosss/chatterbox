package valueobject

import "github.com/google/uuid"

type UserID string

func NewUserID() UserID {
	return UserID(uuid.NewString())
}

func ParseUserID(id string) (UserID, error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return "", err
	}

	return UserID(parsedUUID.String()), nil
}

func (id UserID) String() string {
	return string(id)
}
