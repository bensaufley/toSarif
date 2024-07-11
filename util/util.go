package util

func Keys[K string, V any](m map[K]V) []K {
	vals := make([]K, len(m))
	i := 0
	for k := range m {
		vals[i] = k
		i++
	}
	return vals
}

func AddUnique[T comparable](s []T, x T) ([]T, int) {
	for i, y := range s {
		if x == y {
			return s, i
		}
	}
	return append(s, x), len(s)
}
