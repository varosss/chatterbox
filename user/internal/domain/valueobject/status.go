package valueobject

type Status int

const (
	ActiveStatus  Status = 0
	BlockedStatus Status = 1
	DeletedStatus Status = 2
)

func (status Status) Int() int {
	return int(status)
}
