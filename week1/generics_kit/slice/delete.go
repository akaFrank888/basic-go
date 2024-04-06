package main

import "fmt"

func Delete[T any](s []T, idx int) (T, []T, error) {

	if idx < 0 || idx >= len(s) {
		// 定义一个T的零值
		var zero T
		return zero, s, newErrIndexOutOfRange(idx, len(s))
	}

	deleteValue := s[idx]

	for i := idx; i < len(s)-1; i++ {
		s[i] = s[i+1]
	}

	// 去掉最后一个重复元素 or 删掉index=len-1的元素
	s = s[:len(s)-1]

	// TODO 未实现扩容

	return deleteValue, s, nil

}

func newErrIndexOutOfRange(index int, length int) error {
	return fmt.Errorf("下标超出范围，长度%d，下标%d", length, index)
}
