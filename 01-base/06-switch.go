// 06-switch.go 展示了Go中switch语句的使用
package main

import (
	"fmt"
	"time"
)

func main() {
	// 基本的switch语句
	i := 2
	fmt.Print("写 ", i, " 为：")
	switch i {
	case 1:
		fmt.Println("一")
	case 2:
		fmt.Println("二")
	case 3:
		fmt.Println("三")
	}

	// 在同一个case语句中，你可以使用逗号分隔的列表
	switch time.Now().Weekday() {
	case time.Saturday, time.Sunday:
		fmt.Println("今天是周末")
	default:
		fmt.Println("今天是工作日")
	}

	// 不带表达式的switch等同于switch true
	// 这种形式可以用来编写if-else逻辑的另一种方式
	t := time.Now()
	switch {
	case t.Hour() < 12:
		fmt.Println("现在是上午")
	default:
		fmt.Println("现在是下午")
	}

	// 类型开关(type switch)用于比较类型而非值
	whatType := func(i interface{}) {
		switch t := i.(type) {
		case bool:
			fmt.Println("我是布尔类型")
		case int:
			fmt.Println("我是整数类型")
		default:
			fmt.Printf("我是未知类型 %T\n", t)
		}
	}
	whatType(true)
	whatType(1)
	whatType("hey")
} 