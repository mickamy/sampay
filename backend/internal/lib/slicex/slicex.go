package slicex

func Map[T, U any](ts []T, f func(T) U) []U {
	us := make([]U, len(ts))
	for i := range ts {
		us[i] = f(ts[i])
	}
	return us
}

func MapToValue[T any](ts []*T) []T {
	return Map(ts, func(t *T) T {
		return *t
	})
}

func MapToPointer[T any](ts []T) []*T {
	return Map(ts, func(t T) *T {
		return &t
	})
}

func FlatMap[T, U any](ts []T, f func(T) []U) []U {
	us := make([]U, 0, len(ts))
	for _, t := range ts {
		us = append(us, f(t)...)
	}
	return us
}

func GroupBy[T any, KEY comparable](items []T, getProperty func(T) KEY) map[KEY][]T {
	grouped := make(map[KEY][]T)

	for _, item := range items {
		key := getProperty(item)
		grouped[key] = append(grouped[key], item)
	}

	return grouped
}

func Find[T any](ts []T, f func(T) bool) (T, bool) {
	for _, t := range ts {
		if f(t) {
			return t, true
		}
	}
	var zero T
	return zero, false
}

func Filter[T any](ts []T, f func(T) bool) []T {
	matched := make([]T, 0, len(ts))
	for _, t := range ts {
		if f(t) {
			matched = append(matched, t)
		}
	}
	return matched
}
