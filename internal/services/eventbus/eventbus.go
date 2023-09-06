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

func NewEventBus(driver EventBusDriver) *EventBus {
	return &EventBus{
		driver: driver,
	}
}

type EventBus struct {
	driver EventBusDriver
}

func (e *EventBus) Publish(subject interface{}, eventType domain.EventType) error {
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

	idValue := reflect.ValueOf(subject).FieldByIndex(idField.Index).Interface()
	eventKey := fmt.Sprintf("%s:%v:%s", subjectType, idValue, eventType)

	err := e.driver.Publish(eventKey)
	if err != nil {
		return err
	}

	return nil
}

func (e *EventBus) Subscribe(subject interface{}, wildcardID bool, eventTypes []domain.EventType, callback func(subjectType string, subjectId string, eventType domain.EventType)) error {
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
		eventKey := fmt.Sprintf("%s:%s:%s", subjectType, idString, eventType)
		err := e.driver.Subscribe(eventKey, func(subject string) {
			parts := strings.Split(subject, ":")
			if len(parts) != 3 {
				return
			}
			eventType := domain.ParseEventString(parts[2])
			callback(parts[0], parts[1], eventType)
		})
		if err != nil {
			return err
		}
	}

	return nil
}

type EventBusDriver interface {
	Publish(subject string) error
	Subscribe(subject string, callback func(subject string)) error
}
