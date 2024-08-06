package events

import (
	"sync"
)

type EventHandler func(Event)

type EventBus struct {
	handler map[string][]EventHandler
	lock    sync.RWMutex
}

func NewEventBus() *EventBus {
	return &EventBus{
		handler: make(map[string][]EventHandler),
	}
}

func (bus *EventBus) Emit(event Event) {
	bus.lock.RLock()
	defer bus.lock.RUnlock()

	eventType := getType(event)
	if handlers, ok := bus.handler[eventType]; ok {
		for _, handler := range handlers {
			handler(event)
		}
	}
}

func (bus *EventBus) Register(eventType string, handler EventHandler) {
	bus.lock.Lock()
	defer bus.lock.Unlock()

	if _, ok := bus.handler[eventType]; !ok {
		bus.handler[eventType] = []EventHandler{}
	}
	bus.handler[eventType] = append(bus.handler[eventType], handler)
}

func getType(event Event) string {
	switch event.(type) {
	case AddCommitEvent:
		return "AddCommitEvent"
	case StartMonitorEvent:
		return "StartMonitorEvent"
	default:
		return "Unknown"
	}
}
