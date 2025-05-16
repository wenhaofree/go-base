// 10-functions.go 展示了Go中函数的使用
package main

import "fmt"

// 这是一个函数，接收两个int类型参数，返回它们的和
// 参数类型在变量名之后
func plus(a int, b int) int {
	// Go需要显式的return语句，它不会自动返回最后一个表达式的值
	return a + b
}

// 当多个连续的参数具有相同类型时，可以只声明最后一个参数的类型
// 前面的同类型参数可以省略类型声明
func plusPlus(a, b, c int) int {
	return a + b + c
}

func main() {
	// 使用name(args)语法调用函数
	res := plus(1, 2)
	fmt.Println("1+2 =", res)

	res = plusPlus(1, 2, 3)
	fmt.Println("1+2+3 =", res)
} 