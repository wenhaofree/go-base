// 02-waitgroup.go 展示了如何使用WaitGroup同步多个goroutine
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	fmt.Println("=== Go WaitGroup示例 ===")
	
	// 1. 基本WaitGroup使用
	fmt.Println("\n1. 基本WaitGroup使用")
	basicWaitGroupExample()
	
	// 2. 使用WaitGroup处理多个任务
	fmt.Println("\n2. 使用WaitGroup处理多个任务")
	multipleTasksExample()
	
	// 3. 在函数中使用WaitGroup
	fmt.Println("\n3. 在函数中使用WaitGroup")
	executeTasks()
	
	fmt.Println("所有示例完成")
}

// basicWaitGroupExample 演示WaitGroup的基本用法
func basicWaitGroupExample() {
	// 创建一个WaitGroup
	var wg sync.WaitGroup
	
	// 添加要等待的goroutine数量
	wg.Add(3)
	
	// 启动三个goroutine
	for i := 0; i < 3; i++ {
		go func(id int) {
			// 确保在函数结束时调用Done
			defer wg.Done()
			
			// 模拟工作
			fmt.Printf("Goroutine %d 开始工作\n", id)
			time.Sleep(time.Duration(id*200) * time.Millisecond)
			fmt.Printf("Goroutine %d 完成工作\n", id)
		}(i)
	}
	
	// 等待所有goroutine完成
	fmt.Println("等待所有goroutine完成...")
	wg.Wait()
	fmt.Println("所有goroutine已完成!")
}

// multipleTasksExample 展示如何使用WaitGroup处理一组任务
func multipleTasksExample() {
	// 创建任务列表
	tasks := []string{"任务1", "任务2", "任务3", "任务4", "任务5"}
	
	// 创建WaitGroup
	var wg sync.WaitGroup
	
	// 设置要等待的goroutine数量
	wg.Add(len(tasks))
	
	// 为每个任务启动一个goroutine
	for i, task := range tasks {
		// 启动goroutine处理任务
		go processTask(task, i, &wg)
	}
	
	// 等待所有任务完成
	fmt.Printf("等待 %d 个任务完成...\n", len(tasks))
	wg.Wait()
	fmt.Println("所有任务已完成!")
}

// processTask 处理单个任务并在完成时通知WaitGroup
func processTask(taskName string, id int, wg *sync.WaitGroup) {
	// 确保在函数结束时通知WaitGroup
	defer wg.Done()
	
	// 模拟任务处理
	fmt.Printf("开始处理: %s\n", taskName)
	
	// 模拟处理时间(随任务ID变化)
	processingTime := time.Duration(100*(id+1)) * time.Millisecond
	time.Sleep(processingTime)
	
	fmt.Printf("完成处理: %s (耗时: %v)\n", taskName, processingTime)
}

// executeTasks 展示如何在函数中使用WaitGroup
func executeTasks() {
	var wg sync.WaitGroup
	
	// 定义一个执行任务并返回结果的函数
	executeTask := func(id int) {
		// 确保任务完成时通知WaitGroup
		defer wg.Done()
		
		// 模拟工作
		fmt.Printf("执行任务 %d 开始\n", id)
		time.Sleep(time.Duration(100*id) * time.Millisecond)
		fmt.Printf("执行任务 %d 完成\n", id)
	}
	
	// 启动多个任务
	taskCount := 4
	fmt.Printf("启动 %d 个任务...\n", taskCount)
	
	// 添加任务计数
	wg.Add(taskCount)
	
	// 启动任务
	for i := 0; i < taskCount; i++ {
		go executeTask(i)
	}
	
	// 等待所有任务完成
	fmt.Println("等待所有任务执行完毕...")
	wg.Wait()
	fmt.Println("所有任务执行完毕!")
} 