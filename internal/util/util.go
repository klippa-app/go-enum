package util

func Only[T any, K any](values []K) (filtered []T) {
	for i := range values {
		switch t := any(values[i]).(type) {
		case T:
			filtered = append(filtered, t)
		}
	}
	return filtered
}

func DereferenceOrNew[T any](v *T) T {
	if v == nil {
		v = new(T)
	}

	return *v
}

func Contains[T comparable](arr []T, s T) bool {
	for i := range arr {
		if arr[i] == s {
			return true
		}
	}
	return false
}
