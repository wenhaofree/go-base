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

// 函数声明
func greet(name string) string {
	return fmt.Sprintf("Hello, %s!", name)
}

// 带多个返回值的函数
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("除零错误")
	}
	return a / b, nil
}

// 带命名返回值的函数
func multiply(a, b int) (result int) {
	result = a * b
	return // 裸返回
}

// 带可变参数的函数（类似剩余参数）
func sum(numbers ...int) int {
	total := 0
	for _, num := range numbers {
		total += num
	}
	return total
}

// 函数作为变量（类似函数表达式）
var add = func(a, b int) int {
	return a + b
}

func main() {
	// 使用name(args)语法调用函数
	res := plus(1, 2)
	fmt.Println("1+2 =", res)

	res = plusPlus(1, 2, 3)
	fmt.Println("1+2+3 =", res)

	fmt.Println(greet("Bob"))

	result, err := divide(10, 2)
	if err != nil {
		fmt.Println("错误:", err)
	} else {
		fmt.Println("结果:", result)
	}

	fmt.Println(multiply(4, 6))
	fmt.Println(sum(1, 2, 3, 4, 5))
	fmt.Println(add(5, 3))
}
