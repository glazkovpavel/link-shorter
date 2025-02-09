package event

const (
	EventLinkVisited = "event_link_visited"
)

type Event struct {
	Type string
	Data interface{}
}

type EventBus struct {
	bus chan Event
}

func NewEventBus() *EventBus {
	return &EventBus{
		bus: make(chan Event),
	}
}

func (eb *EventBus) Publish(event Event) {
	eb.bus <- event
}

func (eb *EventBus) Subscribe() <-chan Event {
	return eb.bus
}
