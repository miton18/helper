package slices

func Map[T, U any](items []T, fn func(T) U) []U {
	r := make([]U, len(items))
	for i, item := range items {
		r[i] = fn(item)
	}

	return r
}
