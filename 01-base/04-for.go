// 04-for.go 展示了Go中for循环的使用方法
package main

import "fmt"

func main() {
	// 最基本的for循环形式，单一条件
	i := 1
	for i <= 3 {
		fmt.Println(i)
		i = i + 1
	}

	// 经典的初始化/条件/后续形式的for循环
	for j := 7; j <= 9; j++ {
		fmt.Println(j)
	}

	// 不带条件的for循环将一直重复执行，直到使用break退出循环
	// 或者在函数中使用return返回
	for {
		fmt.Println("loop")
		break
	}

	// 也可以使用continue来跳转到下一个循环迭代
	for n := 0; n <= 5; n++ {
		if n%2 == 0 {
			continue
		}
		fmt.Println(n)
	}

	// 传统 for 循环（Go 的唯一循环结构）
	for i := 0; i < 5; i++ {
		fmt.Println(i)
	}

	// 基于 range 的 for 循环（类似 for...of）
	arr := []int{1, 2, 3, 4, 5}
	for i, val := range arr {
		fmt.Printf("索引 %d: %d\n", i, val)
	}

	// 遍历 map（类似 for...in）
	obj := map[string]int{"a": 1, "b": 2, "c": 3}
	for key, val := range obj {
		fmt.Printf("%s: %d\n", key, val)
	}

	// 遍历切片（只要索引）
	for i := range arr {
		fmt.Printf("索引: %d\n", i)
	}

	// 遍历切片（只要值）
	for _, val := range arr {
		fmt.Printf("值: %d\n", val)
	}
}
