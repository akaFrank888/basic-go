package main

func Add[T any](s []T, idx int, val T) []T {

	if idx < 0 || idx > len(s) {
		panic("下标越界")
	}

	// 先扩容
	var zero T
	s = append(s, zero)

	for i := len(s) - 1; i > idx; i-- {
		s[i] = s[i-1]
	}
	s[idx] = val

	return s

}
