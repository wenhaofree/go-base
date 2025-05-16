// 02-variables.go 展示了Go中变量的声明和使用
package main

import "fmt"

func main() {
	// var 声明一个或多个变量
	var a = "initial"
	fmt.Println(a)

	// 你可以一次性声明多个变量
	var b, c int = 1, 2
	fmt.Println(b, c)

	// Go会自动推断已经初始化的变量类型
	var d = true
	fmt.Println(d)

	// 声明变量且没有给出初始值的，将会被赋予零值
	// 例如，int的零值是0
	var e int
	fmt.Println(e)

	// := 简写语法可以用于替代var声明和初始化
	f := "short"
	fmt.Println(f)
} 