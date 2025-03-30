package stores

type Store[T any] interface {
	Store(x T) error
	View() ([]T, error)
	Exist(x T) (bool, error)
}

