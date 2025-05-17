// 01-basic-goroutine.go 展示了Go中goroutine的基本使用
package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("=== Go基本goroutine示例 ===")
	
	// 1. 启动一个简单的goroutine
	fmt.Println("\n1. 启动一个简单的goroutine")
	go sayHello("世界")
	
	// 主程序继续执行
	fmt.Println("主goroutine继续执行")
	
	// 2. 启动多个goroutine
	fmt.Println("\n2. 启动多个goroutine")
	for i := 0; i < 5; i++ {
		// 注意：在循环中使用变量时需要传递副本
		go func(n int) {
			fmt.Printf("goroutine %d 正在执行\n", n)
		}(i)
	}
	
	// 3. 使用匿名goroutine
	fmt.Println("\n3. 使用匿名goroutine")
	go func() {
		fmt.Println("这是一个匿名goroutine")
	}()
	
	// 等待一段时间让goroutines完成
	// 注意：这只是为了演示，实际应用中应该使用WaitGroup或其他同步方法
	time.Sleep(1 * time.Second)
	fmt.Println("主程序结束")
}

// sayHello 是一个将在goroutine中执行的函数
func sayHello(name string) {
	for i := 0; i < 3; i++ {
		fmt.Printf("Hello, %s - 消息 %d\n", name, i)
		// 模拟一些工作
		time.Sleep(100 * time.Millisecond)
	}
} 