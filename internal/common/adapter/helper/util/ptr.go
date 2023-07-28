package util

func P[T any](t T) *T {
	return &t
}

func UnP[T any](t *T) T {
	if t == nil {
		var a T
		return a
	}

	return *t
}
