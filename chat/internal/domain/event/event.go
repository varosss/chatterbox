package event

import "time"

type Event interface {
	OccurredAt() time.Time
	Name() string
}

type BaseEvent struct {
	occurredAt time.Time
}

func NewBaseEvent() BaseEvent {
	return BaseEvent{occurredAt: time.Now()}
}

func (e BaseEvent) OccurredAt() time.Time {
	return e.occurredAt
}
