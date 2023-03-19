package utils

// Pointer is a helper to make a pointer to the given value.
func Pointer[T any](v T) *T {
	return &v
}
