package ticketeventbus

import (
	"strings"
	"sync"

	"github.com/nil-nil/ticket/internal/domain"
)

func NewBus(separator string) (*ticketEventBus, error) {
	if separator == "" {
		separator = ":"
	}
	return &ticketEventBus{separator: separator, subs: map[string]map[string]map[string][]func(eventKey string, data interface{}){}}, nil
}

type ticketEventBus struct {
	mu        sync.Mutex
	subs      map[string]map[string]map[string][]func(eventKey string, data interface{})
	separator string
}

func (t *ticketEventBus) sub(topic string, f func(eventKey string, data interface{})) error {
	parts := strings.Split(topic, t.separator)
	if len(parts) != 3 {
		return domain.ErrEventKeyInvalid
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	if _, ok := t.subs[parts[0]]; !ok {
		t.subs[parts[0]] = map[string]map[string][]func(eventKey string, data interface{}){}
	}
	if _, ok := t.subs[parts[0]][parts[1]]; !ok {
		t.subs[parts[0]][parts[1]] = map[string][]func(eventKey string, data interface{}){}
	}
	t.subs[parts[0]][parts[1]][parts[2]] = append(t.subs[parts[0]][parts[1]][parts[2]], f)

	return nil
}

func (t *ticketEventBus) match(topic string) []func(eventKey string, data interface{}) {
	return nil
}

func (t *ticketEventBus) Publish(subject string, data interface{}) error {
	return nil
}

func (t *ticketEventBus) Subscribe(subject string, callback func(eventKey string, data interface{})) error {
	return t.sub(subject, callback)
}
