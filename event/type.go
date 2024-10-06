package event

type Event[T any] struct {
	handles map[string]func(T)
}

func (e *Event[T]) Triger(v T) {
	for _, f := range e.handles {
		f(v)
	}
}

func (e *Event[T]) Listen(id string, handle func(T)) {
	if e.handles == nil {
		e.handles = make(map[string]func(T))
	}
	e.handles[id] = handle
}

func (e *Event[T]) UnListen(id string) {
	delete(e.handles, id)
}
