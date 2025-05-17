// 04-mutex.go 展示了互斥锁在并发环境中的使用
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	fmt.Println("=== Go互斥锁(Mutex)示例 ===")

	// 1. 没有互斥锁的问题
	fmt.Println("\n1. 没有互斥锁的并发问题")
	dataRaceExample()

	// 2. 使用互斥锁保护共享资源
	fmt.Println("\n2. 使用互斥锁保护共享资源")
	mutexExample()

	// 3. 使用读写锁优化读多写少的场景
	fmt.Println("\n3. 使用读写锁(RWMutex)优化读多写少的场景")
	rwMutexExample()

	fmt.Println("所有示例完成")
}

// dataRaceExample 展示没有互斥锁的数据竞争问题
func dataRaceExample() {
	// 共享的计数器变量
	counter := 0
	
	// 创建WaitGroup来同步goroutines
	var wg sync.WaitGroup
	
	// 启动多个goroutines同时增加计数器
	numGoroutines := 1000
	wg.Add(numGoroutines)
	
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			
			// 增加计数器（非线程安全操作）
			counter++
		}()
	}
	
	// 等待所有goroutines完成
	wg.Wait()
	
	// 输出最终计数器值
	// 注意：由于数据竞争，这个值很可能小于numGoroutines
	fmt.Printf("预期计数器值: %d, 实际计数器值: %d\n", numGoroutines, counter)
}

// mutexExample 展示如何使用互斥锁保护共享数据
func mutexExample() {
	// 创建一个包含计数器和互斥锁的结构体
	type SafeCounter struct {
		count int
		mutex sync.Mutex
	}
	
	// 初始化安全计数器
	counter := SafeCounter{count: 0}
	
	// 创建WaitGroup来同步goroutines
	var wg sync.WaitGroup
	
	// 启动多个goroutines同时增加计数器
	numGoroutines := 1000
	wg.Add(numGoroutines)
	
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			
			// 锁定互斥锁，保护计数器增加操作
			counter.mutex.Lock()
			counter.count++
			counter.mutex.Unlock()
		}()
	}
	
	// 等待所有goroutines完成
	wg.Wait()
	
	// 输出最终计数器值
	fmt.Printf("使用互斥锁后 - 预期计数器值: %d, 实际计数器值: %d\n", numGoroutines, counter.count)
}

// rwMutexExample 展示读写锁的使用
func rwMutexExample() {
	// 创建一个包含数据、访问计数和读写锁的结构体
	type SafeData struct {
		data        map[string]int
		readCount   int
		writeCount  int
		rwMutex     sync.RWMutex
	}
	
	// 初始化安全数据结构
	safeData := SafeData{
		data:       make(map[string]int),
		readCount:  0,
		writeCount: 0,
	}
	
	// 添加一些初始数据
	safeData.data["初始值"] = 100
	
	// 创建WaitGroup来同步goroutines
	var wg sync.WaitGroup
	
	// 启动多个读取goroutines
	numReaders := 100
	wg.Add(numReaders)
	
	for i := 0; i < numReaders; i++ {
		go func(id int) {
			defer wg.Done()
			
			// 使用读锁(多个goroutine可以同时读取)
			safeData.rwMutex.RLock()
			_ = safeData.data["初始值"] // 读取数据
			safeData.readCount++    // 在真实场景中，这也应该使用互斥锁保护
			safeData.rwMutex.RUnlock()
			
			// 随机休眠一小段时间
			time.Sleep(time.Duration(id%10) * time.Millisecond)
		}(i)
	}
	
	// 启动少量写入goroutines
	numWriters := 10
	wg.Add(numWriters)
	
	for i := 0; i < numWriters; i++ {
		go func(id int) {
			defer wg.Done()
			
			// 使用写锁(独占访问)
			safeData.rwMutex.Lock()
			safeData.data[fmt.Sprintf("key-%d", id)] = id * 100
			safeData.writeCount++ // 在真实场景中，这也应该使用互斥锁保护
			safeData.rwMutex.Unlock()
			
			// 随机休眠一小段时间
			time.Sleep(time.Duration(id*5) * time.Millisecond)
		}(i)
	}
	
	// 等待所有goroutines完成
	wg.Wait()
	
	// 输出最终结果
	fmt.Printf("读写锁示例 - 读取操作: %d, 写入操作: %d\n", safeData.readCount, safeData.writeCount)
	fmt.Printf("数据条目数: %d\n", len(safeData.data))
	
	// 打印部分数据
	safeData.rwMutex.RLock()
	fmt.Println("部分数据内容:")
	count := 0
	for k, v := range safeData.data {
		if count < 5 {
			fmt.Printf("  %s: %d\n", k, v)
			count++
		} else {
			break
		}
	}
	safeData.rwMutex.RUnlock()
} 