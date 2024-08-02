package events

import "reflect"

type EventBus struct {
	handler map[reflect.Type][]interface{}
}

func NewEventBus() *EventBus {
	return &EventBus{
		handler: make(map[reflect.Type][]interface{}),
	}
}

func (bus *EventBus) Emit(event interface{}) {
	eventType := reflect.TypeOf(event)
	if handlers, found := bus.handler[eventType]; found {
		for _, handler := range handlers {
			go func(handler interface{}) {
				handlerValue := reflect.ValueOf(handler)
				handlerValue.Call([]reflect.Value{reflect.ValueOf(event)})
			}(handler)
		}
	}
}

func (bus *EventBus) Register(eventType reflect.Type, handler interface{}) {
	bus.handler[eventType] = append(bus.handler[eventType], handler)
}
