package stores

type MemoryStore[T comparable] struct {
	items []T
	zz    chan T
}

func NewMemoryStore[T comparable]() *MemoryStore[T] {
	m := MemoryStore[T]{
		items: []T{},
		zz:    make(chan T),
	}
	go m.start()
	return &m
}

func (m *MemoryStore[T]) start() {
	for {
		select {
		case x := <-m.zz:
			m.items = append(m.items, x)
		}
	}
}

func (m *MemoryStore[T]) Store(x T) error {
	m.zz <- x
    return nil
}

func (m *MemoryStore[T]) View() ([]T, error) {
	dst := make([]T, len(m.items))
	copy(dst, m.items)
	return dst, nil
}

func (m *MemoryStore[T]) Exist(x T) (bool, error) {
	for _, item := range m.items {
		if item == x {
			return true, nil
		}
	}
	return false, nil
}
