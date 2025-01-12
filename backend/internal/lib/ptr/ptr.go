package ptr

func Of[T any](val T) *T {
	return &val
}

func Unwrap[T any](ptr *T) T {
	if ptr == nil {
		return *new(T)
	}
	return *ptr
}

func Map[T any, U any](val *T, mapper func(*T) U) U {
	if val == nil {
		return *new(U)
	}
	return mapper(val)
}
