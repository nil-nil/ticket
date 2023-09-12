package hubeventbus

import (
	"github.com/leandro-lugaresi/hub"
)

type HubEventBusDriver struct {
	hub *hub.Hub
}

func NewHubEventBusDriver() (*HubEventBusDriver, error) {
	return &HubEventBusDriver{hub.New()}, nil
}

func (h *HubEventBusDriver) Publish(subject string, data interface{}) error {
	h.hub.Publish(hub.Message{
		Name:   subject,
		Fields: hub.Fields{"data": data},
	})
	return nil
}

func (h *HubEventBusDriver) Subscribe(subject string, callback func(eventKey string, data interface{})) error {
	go func(subject string, callback func(eventKey string, data interface{})) {
		sub := h.hub.Subscribe(10, subject)
		for msg := range sub.Receiver {
			data, ok := msg.Fields["data"]
			if !ok {
				return
			}
			callback(msg.Name, data)
		}
	}(subject, callback)

	return nil
}
