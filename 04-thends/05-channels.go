// 05-channels.go 展示了在Go中使用channel进行goroutine之间的通信
package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("=== Go通道(Channel)示例 ===")

	// 1. 基本通道操作
	fmt.Println("\n1. 基本通道操作")
	basicChannelExample()

	// 2. 单向通道
	fmt.Println("\n2. 单向通道")
	directionalChannelExample()

	// 3. 带缓冲的通道
	fmt.Println("\n3. 带缓冲的通道")
	bufferedChannelExample()

	// 4. 使用select多路复用通道
	fmt.Println("\n4. 使用select多路复用通道")
	selectExample()

	// 5. 使用通道实现工作池模式
	fmt.Println("\n5. 使用通道实现工作池模式")
	workerPoolExample()

	fmt.Println("所有示例完成")
}

// basicChannelExample 展示基本的通道操作
func basicChannelExample() {
	// 创建一个无缓冲的字符串通道
	messages := make(chan string)

	// 启动一个goroutine发送数据
	go func() {
		fmt.Println("发送者: 准备发送消息")
		// 向通道发送数据
		messages <- "你好，通道!"
		fmt.Println("发送者: 消息已发送")
	}()

	// 从通道接收数据
	fmt.Println("接收者: 准备接收消息")
	msg := <-messages
	fmt.Println("接收者: 收到消息:", msg)

	// 通道作为同步机制的例子
	done := make(chan bool)
	go func() {
		fmt.Println("Worker: 开始执行任务...")
		time.Sleep(500 * time.Millisecond)
		fmt.Println("Worker: 任务完成")
		done <- true
	}()

	// 等待工作完成
	fmt.Println("主程序: 等待工作完成...")
	<-done
	fmt.Println("主程序: 收到完成信号")
}

// directionalChannelExample 展示单向通道的使用
func directionalChannelExample() {
	// 创建一个双向通道
	channel := make(chan string)

	// 启动发送者goroutine，使用通道的只发送端
	go sender(channel)

	// 启动接收者goroutine，使用通道的只接收端
	go receiver(channel)

	// 给时间让goroutines运行
	time.Sleep(1 * time.Second)
}

// sender 函数只向通道发送数据
func sender(ch chan<- string) { // chan<- 表示只能发送
	ch <- "消息1"
	ch <- "消息2"
	ch <- "消息3"
}

// receiver 函数只从通道接收数据
func receiver(ch <-chan string) { // <-chan 表示只能接收
	for i := 0; i < 3; i++ {
		msg := <-ch
		fmt.Println("接收到:", msg)
	}
}

// bufferedChannelExample 展示带缓冲通道的使用
func bufferedChannelExample() {
	// 创建一个缓冲大小为3的通道
	bufferedCh := make(chan string, 3)

	// 向缓冲通道发送多条消息，不会阻塞直到超出缓冲区
	fmt.Println("向带缓冲的通道发送数据:")
	bufferedCh <- "数据1"
	fmt.Println("- 已发送数据1")
	bufferedCh <- "数据2"
	fmt.Println("- 已发送数据2")
	bufferedCh <- "数据3"
	fmt.Println("- 已发送数据3")

	// 现在通道已满，再发送一条会阻塞
	fmt.Println("通道缓冲已满，开始接收数据:")

	// 从通道接收数据
	fmt.Println("接收:", <-bufferedCh)
	fmt.Println("接收:", <-bufferedCh)
	fmt.Println("接收:", <-bufferedCh)

	// 检查通道容量和长度
	fmt.Printf("通道容量: %d, 当前长度: %d\n", cap(bufferedCh), len(bufferedCh))
}

// selectExample 展示使用select多路复用通道
func selectExample() {
	// 创建两个通道
	ch1 := make(chan string)
	ch2 := make(chan string)

	// 在一个goroutine中向通道1发送数据
	go func() {
		time.Sleep(100 * time.Millisecond)
		ch1 <- "来自通道1的消息"
	}()

	// 在另一个goroutine中向通道2发送数据
	go func() {
		time.Sleep(200 * time.Millisecond)
		ch2 <- "来自通道2的消息"
	}()

	// 使用select监听两个通道
	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-ch1:
			fmt.Println("接收到:", msg1)
		case msg2 := <-ch2:
			fmt.Println("接收到:", msg2)
		}
	}

	// 使用select的超时例子
	fmt.Println("\nselect超时示例:")
	ch := make(chan string)
	
	go func() {
		time.Sleep(500 * time.Millisecond)
		ch <- "迟到的消息"
	}()

	select {
	case msg := <-ch:
		fmt.Println("接收到消息:", msg)
	case <-time.After(200 * time.Millisecond):
		fmt.Println("接收超时!")
	}
}

// workerPoolExample 展示使用通道实现工作池模式
func workerPoolExample() {
	const numJobs = 10
	const numWorkers = 3

	// 创建工作和结果通道
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	// 启动工作协程池
	fmt.Printf("启动 %d 个工作协程...\n", numWorkers)
	for w := 1; w <= numWorkers; w++ {
		go worker(w, jobs, results)
	}

	// 发送工作到jobs通道
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs) // 关闭通道表示没有更多的工作
	fmt.Println("所有工作已发送")

	// 收集所有结果
	for a := 1; a <= numJobs; a++ {
		result := <-results
		fmt.Printf("收到结果: %d\n", result)
	}
}

// worker 工作协程，从jobs接收工作，处理后将结果发送到results
func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Printf("工作协程 %d 开始处理工作 %d\n", id, j)
		// 模拟工作处理时间
		time.Sleep(time.Duration(100) * time.Millisecond)
		// 计算并发送结果
		results <- j * 2
		fmt.Printf("工作协程 %d 完成工作 %d\n", id, j)
	}
} 