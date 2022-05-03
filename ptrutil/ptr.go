package ptrutil

// Dereference dereferences a pointer, returning zero if the pointer is nil.
func Dereference[T any](t *T) T {
	if t == nil {
		var zero T
		return zero
	}

	return *t
}

// Ptr returns the pointer to the given value. Useful for when you want to get the pointer
// of a magic value: Ptr("hello").
func Ptr[T any](t T) *T {
	return &t
}
