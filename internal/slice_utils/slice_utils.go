package slice_utils

// Map applies fn to every item in values, returning a new slice with the returned values.
// If ret is provided, the mapped values are appended to it.
func Map[T, V any](values []T, fn func(T) V, ret []V) []V {
	if ret == nil {
		ret = make([]V, 0, len(values))
	}

	for _, v := range values {
		ret = append(ret, fn(v))
	}

	return ret
}
