package util

// Only takes a list of type []any, and returns a filtered list of only type T.
func Only[T any, K any](values []K) (filtered []T) {
	for i := range values {
		if t, ok := any(values[i]).(T); ok {
			filtered = append(filtered, t)
		}
	}
	return filtered
}

// Contains returns true if arr contains value.
func Contains[T comparable](arr []T, value T) bool {
	for i := range arr {
		if arr[i] == value {
			return true
		}
	}
	return false
}
