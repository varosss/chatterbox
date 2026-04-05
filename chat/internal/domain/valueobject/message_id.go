package valueobject

import "github.com/google/uuid"

type MessageID string

func NewMessageID() MessageID {
	return MessageID(uuid.NewString())
}

func (id MessageID) String() string {
	return string(id)
}
