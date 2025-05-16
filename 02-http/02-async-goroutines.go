// 02-async-goroutines.go 展示了Go中使用Goroutines实现异步并行HTTP请求
package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// 定义一个结构体来存储请求结果
type RequestResult struct {
	URL      string
	Response string
	Error    error
	Duration time.Duration
}

// 发起HTTP GET请求并返回结果
func fetchURL(url string) RequestResult {
	startTime := time.Now()
	
	// 创建一个请求结果对象
	result := RequestResult{
		URL: url,
	}
	
	// 发起HTTP请求
	resp, err := http.Get(url)
	if err != nil {
		result.Error = err
		result.Duration = time.Since(startTime)
		return result
	}
	defer resp.Body.Close()
	
	// 读取响应体
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1024))
	if err != nil {
		result.Error = err
		result.Duration = time.Since(startTime)
		return result
	}
	
	result.Response = string(body)
	result.Duration = time.Since(startTime)
	return result
}

func main() {
	fmt.Println("开始异步HTTP请求示例 (使用Goroutines)")
	
	// 要请求的URL列表
	urls := []string{
		"https://httpbin.org/get",
		"https://httpbin.org/delay/1",
		"https://httpbin.org/delay/2",
		"https://httpbin.org/headers",
	}
	
	fmt.Printf("将并行请求 %d 个URL\n\n", len(urls))
	
	// 记录开始时间
	startTime := time.Now()
	
	// 创建一个等待组，用于等待所有goroutine完成
	var wg sync.WaitGroup
	
	// 创建一个通道，用于接收结果
	resultChan := make(chan RequestResult, len(urls))
	
	// 为每个URL启动一个goroutine
	for _, url := range urls {
		wg.Add(1) // 增加等待计数
		
		// 启动goroutine发起请求
		go func(url string) {
			defer wg.Done() // 完成时减少等待计数
			
			// 发起请求并将结果发送到通道
			result := fetchURL(url)
			resultChan <- result
			
		}(url) // 将URL作为参数传递给匿名函数
	}
	
	// 启动另一个goroutine来关闭结果通道
	go func() {
		// 等待所有请求goroutine完成
		wg.Wait()
		// 关闭结果通道
		close(resultChan)
	}()
	
	// 从通道接收结果并打印
	for result := range resultChan {
		if result.Error != nil {
			fmt.Printf("请求 %s 失败: %v (耗时: %v)\n", 
				result.URL, result.Error, result.Duration)
		} else {
			fmt.Printf("请求 %s 成功 (耗时: %v)\n", 
				result.URL, result.Duration)
			fmt.Printf("响应体摘要: %s...\n\n", 
				result.Response[:min(50, len(result.Response))])
		}
	}
	
	// 计算总耗时
	totalDuration := time.Since(startTime)
	fmt.Printf("所有请求完成，总耗时: %v\n", totalDuration)
}

// min函数返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
} 