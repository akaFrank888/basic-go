package main

func Find[T any](s []T, idx int) T {
	if idx < 0 || idx >= len(s) {
		panic("下标越界")
	}
	return s[idx]
}
