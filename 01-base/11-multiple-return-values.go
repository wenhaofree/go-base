// 11-multiple-return-values.go 展示了Go中多返回值函数的使用
package main

import "fmt"

// (int, int)在这个函数签名中表示这个函数返回2个int
func vals() (int, int) {
	return 3, 7
}

// 这是一个返回两个字符串的函数
func getNames() (string, string) {
	return "张三", "李四"
}

// 你可以使用命名返回值来声明返回值变量的名字
// 这样在函数中可以像使用普通变量一样使用它们
// 它们会被视作已经被声明的变量
func split(sum int) (x, y int) {
	x = sum * 4 / 9
	y = sum - x
	// 没有参数的return语句会返回命名返回值的当前值
	// 这就是所谓的"裸"返回
	return
}

func main() {
	// 获取函数返回的多个值
	a, b := vals()
	fmt.Println("a =", a)
	fmt.Println("b =", b)

	// 如果你只想要返回值的一部分，使用空白标识符_
	firstName, _ := getNames()
	fmt.Println("名字 =", firstName)

	// 使用命名返回值的函数
	x, y := split(17)
	fmt.Println("x =", x)
	fmt.Println("y =", y)
} 