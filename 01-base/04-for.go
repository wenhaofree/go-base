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
} 