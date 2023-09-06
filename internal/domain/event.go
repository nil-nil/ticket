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

func ParseEventString(s string) EventType {
	switch s {
	case "create":
		return CreateEvent
	case "update":
		return UpdateEvent
	case "delete":
		return DeleteEvent
	}
	return UnknownEvent
}

type EventBus interface {
	Publish(subject interface{}, eventType EventType) error
	Subscribe(subject interface{}, wildcardID bool, eventTypes []EventType, callback func(subjectType string, subjectId string, eventType EventType)) error
}
