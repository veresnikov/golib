package event

type Handler interface {
	Handle(event Event) error
}

type Dispatcher interface {
	Dispatch(event Event) error
	Subscribe(handler Handler)
}

func NewDispatcher() Dispatcher {
	return &dispatcher{}
}

type dispatcher struct {
	handlers []Handler
}

func (d *dispatcher) Dispatch(event Event) error {
	for _, handler := range d.handlers {
		err := handler.Handle(event)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *dispatcher) Subscribe(handler Handler) {
	d.handlers = append(d.handlers, handler)
}
