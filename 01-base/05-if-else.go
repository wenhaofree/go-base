// 05-if-else.go 展示了Go中条件语句的使用
package main

import "fmt"

func main() {
	// 最基本的例子
	if 7%2 == 0 {
		fmt.Println("7是偶数")
	} else {
		fmt.Println("7是奇数")
	}

	// 你可以只用if语句，不需要else
	if 8%4 == 0 {
		fmt.Println("8能被4整除")
	}

	// 条件前可以有一个声明语句；在这个声明语句定义的变量
	// 可以在所有的条件分支中使用
	if num := 9; num < 0 {
		fmt.Println(num, "是负数")
	} else if num < 10 {
		fmt.Println(num, "是个位数")
	} else {
		fmt.Println(num, "有多位")
	}

	// 注意：Go中没有"三元if"语句，所以即使是基本条件判断也需要用完整的if语句
} 