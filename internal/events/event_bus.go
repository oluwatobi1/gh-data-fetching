package events

import (
	"sync"
)

type EventHandler func(Event)

type EventBus struct {
	handler        map[string][]EventHandler
	lock           sync.RWMutex
	workerPoolSize int
}

func NewEventBus(workerPoolSize int) *EventBus {
	return &EventBus{
		handler:        make(map[string][]EventHandler),
		workerPoolSize: workerPoolSize,
	}
}

func (bus *EventBus) Emit(event Event) {
	bus.lock.RLock()

	handlers, ok := bus.handler[event.EventType()]
	bus.lock.RUnlock()
	if !ok {
		return
	}

	// Channel to control/ limit the number of active workers
	workerChan := make(chan struct{}, bus.workerPoolSize)
	eventChan := make(chan func(Event), len(handlers))

	for i := 0; i < bus.workerPoolSize; i++ {
		go func() {
			for handler := range eventChan {
				workerChan <- struct{}{}
				handler(event)
				<-workerChan
			}
		}()
	}

	for _, handler := range handlers {
		eventChan <- handler
	}

	close(eventChan)
}
func (bus *EventBus) Register(eventType string, handler EventHandler) {
	bus.lock.Lock()
	defer bus.lock.Unlock()

	if _, ok := bus.handler[eventType]; !ok {
		bus.handler[eventType] = []EventHandler{}
	}
	bus.handler[eventType] = append(bus.handler[eventType], handler)
}
