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

func (t *ticketEventBus) match(topic string) (subs []func(eventKey string, data interface{})) {
	parts := strings.Split(topic, t.separator)
	if len(parts) != 3 {
		return nil
	}

	if _, ok := t.subs[parts[0]]; ok {
		l1 := parts[0]
		if _, ok := t.subs[l1][parts[1]]; ok {
			l2 := parts[1]
			if _, ok := t.subs[l1][l2][parts[2]]; ok {
				subs = append(subs, t.subs[l1][l2][parts[2]]...)
			}
			if _, ok := t.subs[l1][l2]["*"]; ok && parts[2] != "*" {
				subs = append(subs, t.subs[l1][l2]["*"]...)
			}
		}
		if _, ok := t.subs[l1]["*"]; ok && parts[1] != "*" {
			l2 := "*"
			if _, ok := t.subs[l1][l2][parts[2]]; ok {
				subs = append(subs, t.subs[l1][l2][parts[2]]...)
			}
			if _, ok := t.subs[l1][l2]["*"]; ok && parts[2] != "*" {
				subs = append(subs, t.subs[l1][l2]["*"]...)
			}
		}
	}
	if _, ok := t.subs["*"]; ok && parts[0] != "*" {
		l1 := "*"
		if _, ok := t.subs[l1][parts[1]]; ok {
			l2 := parts[1]
			if _, ok := t.subs[l1][l2][parts[2]]; ok {
				subs = append(subs, t.subs[l1][l2][parts[2]]...)
			}
			if _, ok := t.subs[l1][l2]["*"]; ok && parts[2] != "*" {
				subs = append(subs, t.subs[l1][l2]["*"]...)
			}
		}
		if _, ok := t.subs[l1]["*"]; ok && parts[1] != "*" {
			l2 := "*"
			if _, ok := t.subs[l1][l2][parts[2]]; ok {
				subs = append(subs, t.subs[l1][l2][parts[2]]...)
			}
			if _, ok := t.subs[l1][l2]["*"]; ok && parts[2] != "*" {
				subs = append(subs, t.subs[l1][l2]["*"]...)
			}
		}
	}
	return
}

func (t *ticketEventBus) Publish(subject string, data interface{}) error {
	parts := strings.Split(subject, t.separator)
	if len(parts) != 3 {
		return domain.ErrEventKeyInvalid
	}

	funcs := t.match(subject)
	for _, f := range funcs {
		go f(subject, data)
	}

	return nil
}

func (t *ticketEventBus) Subscribe(subject string, callback func(eventKey string, data interface{})) error {
	return t.sub(subject, callback)
}
