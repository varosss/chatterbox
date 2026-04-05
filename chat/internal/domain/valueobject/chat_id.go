package valueobject

import "github.com/google/uuid"

type ChatID string

func NewChatID() ChatID {
	return ChatID(uuid.NewString())
}

func ParseChatID(id string) (ChatID, error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return "", err
	}

	return ChatID(parsedUUID.String()), nil
}

func (id ChatID) String() string {
	return string(id)
}
