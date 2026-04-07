package valueobject

import "github.com/google/uuid"

type MessageID string

func NewMessageID() MessageID {
	return MessageID(uuid.NewString())
}

func (id MessageID) String() string {
	return string(id)
}

func ParseMessageID(id string) (MessageID, error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return "", err
	}

	return MessageID(parsedUUID.String()), nil
}
