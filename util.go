package dsa

import "golang.org/x/exp/constraints"

// Default return v or default if v has zero value
func Default[T comparable](v T, def T) T {
	var zero T
	if zero == v {
		return def
	}
	return v
}

// Max return maximum values
func Max[T constraints.Ordered](v1, v2 T) T {
	if v1 > v2 {
		return v1
	}
	return v2
}

// MaxSlice return maximum values within slices
func MaxSlice[T constraints.Ordered](va ...T) T {
	var ma T

	if len(va) > 0 {
		ma = va[0]
	}
	for _, v := range va {
		if v > ma {
			ma = v
		}
	}
	return ma
}

// Min return minum between two values
func Min[T constraints.Ordered](v1, v2 T) T {
	if v1 < v2 {
		return v1
	}
	return v2
}

// MinSlice return minimum values for given slices
func MinSlice[T constraints.Ordered](va ...T) T {
	var mi T

	if len(va) > 0 {
		mi = va[0]
	}
	for _, v := range va {
		if v < mi {
			mi = v
		}
	}
	return mi
}
