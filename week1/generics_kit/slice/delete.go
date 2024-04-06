package main

func Delete[T any](s []T, idx int) (T, []T) {

	res := Find(s, idx)

	for i := idx; i < len(s)-1; i++ {
		s[i] = s[i+1]
	}

	// 去掉最后一个重复元素 or 删掉index=len-1的元素
	s = s[:len(s)-1]

	return res, s

}
