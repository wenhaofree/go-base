// 01-values.go 展示了Go的基本值类型
package main

import "fmt"

func main() {
	// 字符串可以通过 + 连接
	fmt.Println("go" + "lang")

	// 整数和浮点数
	fmt.Println("1+1 =", 1+1)
	fmt.Println("7.0/3.0 =", 7.0/3.0)

	// 布尔值以及布尔运算
	fmt.Println(true && false) // 逻辑与
	fmt.Println(true || false) // 逻辑或
	fmt.Println(!true)         // 逻辑非
} 