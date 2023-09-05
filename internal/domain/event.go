package domain

type EventType int

const (
	UnknownEvent EventType = iota
	CreateEvent
	UpdateEvent
	DeleteEvent
)
