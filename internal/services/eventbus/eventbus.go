package eventbus

import (
	"errors"
	"fmt"
	"reflect"

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

type EventBusDriver interface {
	Publish(subject string) error
}
