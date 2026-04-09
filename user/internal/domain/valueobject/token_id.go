package valueobject

import "github.com/google/uuid"

type TokenID string

func NewTokenID() TokenID {
	return TokenID(uuid.NewString())
}

func ParseTokenID(id string) (TokenID, error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return "", err
	}

	return TokenID(parsedUUID.String()), nil
}

func (id TokenID) String() string {
	return string(id)
}
