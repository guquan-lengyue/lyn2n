package status

var WindowsHideStatus = GlobalStatus[bool]{}

type GlobalStatus[T any] struct {
	v       T
	handles map[string]func(T)
}

func (s *GlobalStatus[T]) Set(v T) {
	s.v = v
	for _, f := range s.handles {
		f(s.v)
	}
}
func (s *GlobalStatus[T]) Get() T {
	return s.v
}

func (s *GlobalStatus[T]) Listen(id string, handle func(T)) bool {
	if s.handles == nil {
		s.handles = make(map[string]func(T))
	}
	if _, ok := s.handles[id]; !ok {
		s.handles[id] = handle
		return true
	} else {
		return false
	}
}

func (s *GlobalStatus[T]) UnListen(id string) {
	delete(s.handles, id)
}
