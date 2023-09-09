package eventbus

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/nil-nil/ticket/internal/domain"
)

var (
	ErrNotAStruct = errors.New("subject is not a struct")
	ErrNoIDField  = errors.New("subject has no field tagged as eventbus:\"id\"")
)

func NewEventBus[T any](driver EventBusDriver[T]) *EventBus[T] {
	return &EventBus[T]{
		driver: driver,
	}
}

type EventBus[T any] struct {
	driver EventBusDriver[T]
}

func (e *EventBus[T]) Publish(data T, eventType domain.EventType) error {
	subjectType := reflect.TypeOf(data)
	if subjectType.Kind() != reflect.Struct {
		return ErrNotAStruct
	}

	var idField reflect.StructField
	for _, f := range reflect.VisibleFields(subjectType) {
		if f.Tag.Get("eventbus") == "id" {
			idField = f
			break
		}
	}
	if idField.Index == nil {
		return ErrNoIDField
	}

	idValue := reflect.ValueOf(data).FieldByIndex(idField.Index).Interface()
	eventKey := fmt.Sprintf("%s:%v:%s", subjectType, idValue, eventType)

	err := e.driver.Publish(eventKey, data)
	if err != nil {
		return err
	}

	return nil
}

func (e *EventBus[T]) Subscribe(subject T, wildcardID bool, eventTypes []domain.EventType, callback func(subject T, eventType domain.EventType)) error {
	subjectType := reflect.TypeOf(subject)
	if subjectType.Kind() != reflect.Struct {
		return ErrNotAStruct
	}

	var idField reflect.StructField
	for _, f := range reflect.VisibleFields(subjectType) {
		if f.Tag.Get("eventbus") == "id" {
			idField = f
			break
		}
	}
	if idField.Index == nil {
		return ErrNoIDField
	}

	var idString string
	if wildcardID {
		idString = "*"
	} else {
		idValue := reflect.ValueOf(subject).FieldByIndex(idField.Index).Interface()
		idString = fmt.Sprintf("%v", idValue)
	}

	for _, eventType := range eventTypes {
		err := e.driver.Subscribe(fmt.Sprintf("%s:%s:%s", subjectType, idString, eventType), func(eventKey string, data T) {
			keyParts := strings.Split(eventKey, ":")
			if len(keyParts) != 3 {
				return
			}
			callback(data, domain.ParseEventString(keyParts[2]))
		})
		if err != nil {
			return err
		}
	}

	return nil
}

type EventBusDriver[T any] interface {
	Publish(subject string, data T) error
	Subscribe(subject string, callback func(eventKey string, data T)) error
}
