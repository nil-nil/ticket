package ticketeventbus

import (
	"strings"
	"sync"

	"github.com/nil-nil/ticket/internal/domain"
)

// NewBus creates a new event bus to publish and observe events.
//
// Separator is required for publishing and subscribing, and should match the separator provided here.
// If an empty string is passed, the default separator ":" will be used
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

// sub adds a callback for a topic
//
// l1, l2, and l3 are the 3 parts of the topic. Each can be a wildcard "*".
func (t *ticketEventBus) sub(l1, l2, l3 string, f func(eventKey string, data interface{})) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if _, ok := t.subs[l1]; !ok {
		t.subs[l1] = map[string]map[string][]func(eventKey string, data interface{}){}
	}
	if _, ok := t.subs[l1][l2]; !ok {
		t.subs[l1][l2] = map[string][]func(eventKey string, data interface{}){}
	}
	t.subs[l1][l2][l3] = append(t.subs[l1][l2][l3], f)

	return nil
}

// match returns all the callbacks matching a topic
//
// l1, l2, and l3 are the 3 parts of the topic. Each can be a wildcard "*".
func (t *ticketEventBus) match(l1, l2, l3 string) (subs []func(eventKey string, data interface{})) {
	if _, ok := t.subs[l1]; ok {
		if _, ok := t.subs[l1][l2]; ok {
			if _, ok := t.subs[l1][l2][l3]; ok {
				subs = append(subs, t.subs[l1][l2][l3]...)
			}
			if _, ok := t.subs[l1][l2]["*"]; ok && l3 != "*" {
				subs = append(subs, t.subs[l1][l2]["*"]...)
			}
		}
		if _, ok := t.subs[l1]["*"]; ok && l2 != "*" {
			if _, ok := t.subs[l1]["*"][l3]; ok {
				subs = append(subs, t.subs[l1]["*"][l3]...)
			}
			if _, ok := t.subs[l1]["*"]["*"]; ok && l3 != "*" {
				subs = append(subs, t.subs[l1]["*"]["*"]...)
			}
		}
	}
	if _, ok := t.subs["*"]; ok && l1 != "*" {
		if _, ok := t.subs["*"][l2]; ok {
			if _, ok := t.subs["*"][l2][l3]; ok {
				subs = append(subs, t.subs["*"][l2][l3]...)
			}
			if _, ok := t.subs["*"][l2]["*"]; ok && l3 != "*" {
				subs = append(subs, t.subs["*"][l2]["*"]...)
			}
		}
		if _, ok := t.subs["*"]["*"]; ok && l2 != "*" {
			if _, ok := t.subs["*"]["*"][l3]; ok {
				subs = append(subs, t.subs["*"]["*"][l3]...)
			}
			if _, ok := t.subs["*"]["*"]["*"]; ok && l3 != "*" {
				subs = append(subs, t.subs["*"]["*"]["*"]...)
			}
		}
	}
	return
}

// Publish sends an event to all the callbacks registered for matching topics
//
// The subject should be a 3 parts separated by the Bus's separator. Each part can be a wildcard "*" or a value.
// e.g. if the separator is ":" the subject could be "example.*.5" or "example.test.1"
func (t *ticketEventBus) Publish(subject string, data interface{}) error {
	parts := strings.Split(subject, t.separator)
	if len(parts) != 3 {
		return domain.ErrEventKeyInvalid
	}

	funcs := t.match(parts[0], parts[1], parts[2])
	for _, f := range funcs {
		go f(subject, data)
	}

	return nil
}

// Subscribe registers a callback to receive events
//
// The subject should be a 3 parts separated by the Bus's separator. Each part can be a wildcard "*" or a value.
// e.g. if the separator is ":" the subject could be "example.*.1" or "example.test.1"
func (t *ticketEventBus) Subscribe(subject string, callback func(eventKey string, data interface{})) error {
	parts := strings.Split(subject, t.separator)
	if len(parts) != 3 {
		return domain.ErrEventKeyInvalid
	}
	return t.sub(parts[0], parts[1], parts[2], callback)
}
