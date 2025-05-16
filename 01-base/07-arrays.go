// 07-arrays.go 展示了Go中数组的使用
package main

import "fmt"

func main() {
	// 在Go中，数组是一个具有明确长度的元素序列
	// 默认情况下，数组的零值是零值元素的重复
	// 对于整数，零值是0
	var a [5]int
	fmt.Println("初始化:", a)

	// 可以使用数组[index] = value 语法来设置数组指定索引上的值
	// 或者用 array[index] 得到值
	a[4] = 100
	fmt.Println("设置:", a)
	fmt.Println("获取:", a[4])

	// 内置函数len返回数组的长度
	fmt.Println("长度:", len(a))

	// 声明并初始化数组
	b := [5]int{1, 2, 3, 4, 5}
	fmt.Println("声明:", b)

	// 声明多维数组
	var twoD [2][3]int
	for i := 0; i < 2; i++ {
		for j := 0; j < 3; j++ {
			twoD[i][j] = i + j
		}
	}
	fmt.Println("二维数组:", twoD)
} 