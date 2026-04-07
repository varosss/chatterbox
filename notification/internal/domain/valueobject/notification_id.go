package valueobject

import "github.com/google/uuid"

type NotificationID string

func NewNotificationID() NotificationID {
	return NotificationID(uuid.NewString())
}

func ParseNotificationID(id string) (NotificationID, error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return "", err
	}

	return NotificationID(parsedUUID.String()), nil
}

func (id NotificationID) String() string {
	return string(id)
}
