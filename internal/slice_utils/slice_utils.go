package slice_utils

import "strings"

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

// ContainsPrefix reports whether v is present in s,
// or if an entry in s, when followed by sep, is a prefix for v,
// returning the matched item, if any.
func ContainsPrefix(s []string, v, sep string) (string, bool) {
	for _, want := range s {
		if strings.HasPrefix(v, want) {
			tmp := v[len(want):]
			if tmp == "" || strings.HasPrefix(tmp, sep) {
				return want, true
			}
		}
	}

	return "", false
}
