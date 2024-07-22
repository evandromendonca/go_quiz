package handlers

type RepositoryInterface[T any] interface {
	Create(T) (T, error)
}
