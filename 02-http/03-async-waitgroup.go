// 03-async-waitgroup.go 展示了Go中使用sync.WaitGroup实现异步并行HTTP请求
package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// 定义一个结构体来存储请求结果
type Result struct {
	URL      string
	Response string
	Error    error
	Duration time.Duration
}

func main() {
	fmt.Println("开始异步HTTP请求示例 (使用WaitGroup)")
	
	// 要请求的URL列表
	urls := []string{
		"https://httpbin.org/get",
		"https://httpbin.org/delay/1",
		"https://httpbin.org/delay/2",
		"https://httpbin.org/ip",
	}
	
	fmt.Printf("将并行请求 %d 个URL\n\n", len(urls))
	
	// 记录开始时间
	startTime := time.Now()
	
	// 创建一个等待组来同步goroutines
	var wg sync.WaitGroup
	
	// 使用互斥锁保护结果切片
	var mu sync.Mutex
	results := make([]Result, 0, len(urls))
	
	// 为每个URL创建一个goroutine进行请求
	for _, url := range urls {
		// 增加等待计数
		wg.Add(1)
		
		// 启动goroutine
		go func(url string) {
			// 确保在此goroutine完成时调用Done
			defer wg.Done()
			
			// 记录此请求的开始时间
			requestStart := time.Now()
			
			// 创建HTTP客户端并设置超时
			client := &http.Client{
				Timeout: 30 * time.Second,
			}
			
			// 发起GET请求
			resp, err := client.Get(url)
			
			// 初始化结果
			result := Result{
				URL: url,
			}
			
			// 处理错误情况
			if err != nil {
				result.Error = err
				result.Duration = time.Since(requestStart)
				
				// 使用互斥锁安全地添加结果
				mu.Lock()
				results = append(results, result)
				mu.Unlock()
				
				return
			}
			
			// 确保响应体被关闭
			defer resp.Body.Close()
			
			// 读取响应体
			body, err := io.ReadAll(io.LimitReader(resp.Body, 1024))
			if err != nil {
				result.Error = err
				result.Duration = time.Since(requestStart)
				
				mu.Lock()
				results = append(results, result)
				mu.Unlock()
				
				return
			}
			
			// 设置成功结果
			result.Response = string(body)
			result.Duration = time.Since(requestStart)
			
			// 安全地添加结果
			mu.Lock()
			results = append(results, result)
			mu.Unlock()
			
		}(url) // 将当前URL传递给goroutine
	}
	
	// 等待所有goroutine完成
	wg.Wait()
	
	// 计算总耗时
	totalDuration := time.Since(startTime)
	
	// 打印结果
	fmt.Println("所有请求已完成:")
	for _, result := range results {
		if result.Error != nil {
			fmt.Printf("× 请求 %s 失败: %v (耗时: %v)\n", 
				result.URL, result.Error, result.Duration)
		} else {
			fmt.Printf("✓ 请求 %s 成功 (耗时: %v)\n", 
				result.URL, result.Duration)
			fmt.Printf("  响应体摘要: %s...\n\n", 
				truncateString(result.Response, 50))
		}
	}
	
	fmt.Printf("\n总耗时: %v (如果是串行请求会更长)\n", totalDuration)
	
	// 对比同步执行的耗时
	var totalSyncDuration time.Duration
	for _, result := range results {
		totalSyncDuration += result.Duration
	}
	fmt.Printf("如果串行执行，预计总耗时: %v\n", totalSyncDuration)
	fmt.Printf("并行执行节省时间: %v\n", totalSyncDuration-totalDuration)
}

// 截断字符串到指定长度
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
} 