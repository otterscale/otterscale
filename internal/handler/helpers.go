package handler

// deref returns the value pointed to by ptr, or def if ptr is nil.
func deref[T any](ptr *T, def T) T {
	if ptr != nil {
		return *ptr
	}
	return def
}
