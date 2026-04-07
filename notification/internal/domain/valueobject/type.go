package valueobject

type NotificationType string

const (
	NewMessageNotificationType NotificationType = "new_message"
)

func (t NotificationType) String() string {
	return string(t)
}
