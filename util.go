package dsa

// Default return v or default if v has zero value
func Default[T comparable](v T, def T) T {
	var zero T
	if zero == v {
		return def
	}
	return v
}
