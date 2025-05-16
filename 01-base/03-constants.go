// 03-constants.go 展示了Go中常量的声明和使用
package main

import (
	"fmt"
	"math"
)

// 常量可以在任何var变量可以使用的地方使用
const s string = "constant"

func main() {
	fmt.Println(s)

	// 常量表达式可以执行任意精度的运算
	const n = 500000000

	// 常数表达式可以是类型不确定的，直到被赋予一个类型
	// 例如，根据上下文确定类型
	const d = 3e20 / n
	fmt.Println(d)

	// 数值型常量没有具体类型，除非指定一个类型
	// 或者通过上下文来确定，比如变量赋值或函数调用
	fmt.Println(int64(d))

	// 数学运算中也可以使用常量
	fmt.Println(math.Sin(n))
} 