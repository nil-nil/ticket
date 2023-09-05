package domain

type EventType int

const (
	UnknownEvent EventType = iota
	CreateEvent
	UpdateEvent
	DeleteEvent
)

func (e EventType) String() string {
	switch e {
	case CreateEvent:
		return "create"
	case UpdateEvent:
		return "update"
	case DeleteEvent:
		return "delete"
	}
	return "unknown"
}
