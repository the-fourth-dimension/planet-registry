package repositories

type RepositoryResult[T any] struct {
	Result *T
	Error  error
}
