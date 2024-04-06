package main

import "fmt"

func main() {

	s := []string{"111", "222", "333"}
	res, sResult := Delete[string](s, 1)
	fmt.Printf("删除元素%v后，切片为：%v", res, sResult)

}

func Slice() {

	// 1. 两种初始化
	s1 := []int{1, 2, 3}
	// 下面是初始化的推荐写法
	s2 := make([]int, 3, 4) // len = 3 , cap = 4 ; s2 := make([]int, 3)  ==>  len = 3

	// 2. []访问元素
	println(s1[0], s2[0])

	// 3. append追加元素
	s1 = append(s1, 4)

	// 4. 子切片
	s2_child := s2[1:]

	// 5. for-range遍历
	for idx, val := range s2_child {
		println(idx, val)
	}
}
