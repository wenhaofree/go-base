// 08-slices.go 展示了Go中切片的使用
package main

import "fmt"

func main() {
	// 与数组不同，切片的类型仅由它所包含的元素决定（不需要元素个数）
	// 创建一个长度为3的字符串类型切片
	s := make([]string, 3)
	fmt.Println("空切片:", s)

	// 可以像数组一样设置和获取值
	s[0] = "a"
	s[1] = "b"
	s[2] = "c"
	fmt.Println("设置:", s)
	fmt.Println("获取:", s[2])

	// 可以获取长度
	fmt.Println("长度:", len(s))

	// 除了基本操作外，切片还支持一些使其比数组更丰富的操作
	// 其中一个是append函数，它会返回一个包含了原有切片所有元素和新增元素的新切片
	s = append(s, "d")
	s = append(s, "e", "f")
	fmt.Println("添加后:", s)

	// 切片也可以被复制
	c := make([]string, len(s))
	copy(c, s)
	fmt.Println("复制后:", c)

	// 支持"切片"操作，用于获取切片的一段
	// s[low:high] 包含元素s[low]到s[high-1]
	l := s[2:5]
	fmt.Println("切片2->5:", l)

	// 省略下标代表从0开始
	l = s[:5]
	fmt.Println("切片:5:", l)

	// 省略上标代表直到末尾
	l = s[2:]
	fmt.Println("切片2::", l)

	// 声明并初始化一个切片
	t := []string{"g", "h", "i"}
	fmt.Println("声明:", t)

	// 创建二维切片
	twoD := make([][]int, 3)
	for i := 0; i < 3; i++ {
		innerLen := i + 1
		twoD[i] = make([]int, innerLen)
		for j := 0; j < innerLen; j++ {
			twoD[i][j] = i + j
		}
	}
	fmt.Println("二维切片:", twoD)
} 