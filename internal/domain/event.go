package domain

import (
	"errors"
	"fmt"
	"strings"
)

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

var (
	ErrEventValueInvalid  = errors.New("event data is invalid")
	ErrEventPrefixInvalid = errors.New("not a valid event prefix")
	ErrEventKeyInvalid    = errors.New("not a valid event subject key")
)

func NewEventBus[T any](prefix string, driver eventBusDriver) (*EventBus[T], error) {
	if prefix == "" {
		return nil, ErrEventPrefixInvalid
	}
	return &EventBus[T]{
		driver: driver,
		prefix: prefix,
	}, nil
}

type EventBus[T any] struct {
	driver eventBusDriver
	prefix string
}

func (e *EventBus[T]) Publish(ID string, eventType EventType, data T) error {
	if ID == "" {
		return ErrEventKeyInvalid
	}
	eventKey := fmt.Sprintf("%s:%s:%s", e.prefix, ID, eventType)

	err := e.driver.Publish(eventKey, data)
	if err != nil {
		return err
	}

	return nil
}

func (e *EventBus[T]) Subscribe(ID *string, eventTypes []EventType, callback func(eventType EventType, data T)) error {
	idString := "*"
	if ID != nil {
		idString = *ID
	}

	for _, eventType := range eventTypes {
		err := e.driver.Subscribe(fmt.Sprintf("%s:%s:%s", e.prefix, idString, eventType), func(eventKey string, data interface{}) {
			keyParts := strings.Split(eventKey, ":")
			if len(keyParts) != 3 {
				return
			}

			val, ok := data.(T)
			if !ok {
				return
			}

			callback(ParseEventString(keyParts[2]), val)
		})
		if err != nil {
			return err
		}
	}

	return nil
}

type eventBusDriver interface {
	Publish(subject string, data interface{}) error
	Subscribe(subject string, callback func(eventKey string, data interface{})) error
}
