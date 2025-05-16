// 04-async-rate-limit.go 展示了Go中使用限流控制并发HTTP请求
package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// 定义结果结构体
type HTTPResult struct {
	URL         string
	StatusCode  int
	Body        string
	ElapsedTime time.Duration
	Error       error
}

func main() {
	fmt.Println("开始限流控制的并发HTTP请求示例")
	
	// 要访问的URL列表
	urls := []string{
		"https://httpbin.org/get",
		"https://httpbin.org/delay/1",
		"https://httpbin.org/status/200",
		"https://httpbin.org/delay/2",
		"https://httpbin.org/ip",
		"https://httpbin.org/headers",
		"https://httpbin.org/user-agent",
		"https://httpbin.org/cookies",
	}
	
	// 最大并发数 - 控制同时进行的HTTP请求数量
	maxConcurrent := 3
	fmt.Printf("URL总数: %d, 最大并发数: %d\n\n", len(urls), maxConcurrent)
	
	// 开始计时
	startTime := time.Now()
	
	// 创建结果通道，用于收集请求结果
	resultsChan := make(chan HTTPResult)
	
	// 创建等待组，用于等待所有goroutine完成
	var wg sync.WaitGroup
	
	// 创建一个令牌通道用于限流
	// 缓冲区大小决定了最大并发数
	semaphore := make(chan struct{}, maxConcurrent)
	
	// 启动一个goroutine来处理所有URL请求
	go func() {
		for _, url := range urls {
			// 增加等待组计数
			wg.Add(1)
			
			// 获取令牌（阻塞直到有可用的并发槽）
			semaphore <- struct{}{}
			
			// 为每个URL启动一个goroutine
			go func(url string) {
				// 完成时释放令牌和减少等待计数
				defer func() {
					<-semaphore // 释放令牌
					wg.Done()   // 标记此goroutine已完成
				}()
				
				// 执行HTTP请求并测量时间
				requestStart := time.Now()
				
				// 创建HTTP客户端
				client := &http.Client{
					Timeout: 10 * time.Second,
				}
				
				// 发起请求
				resp, err := client.Get(url)
				
				// 创建结果对象
				result := HTTPResult{
					URL: url,
				}
				
				// 检查错误
				if err != nil {
					result.Error = err
					result.ElapsedTime = time.Since(requestStart)
					resultsChan <- result
					return
				}
				
				// 确保响应体被关闭
				defer resp.Body.Close()
				
				// 记录状态码
				result.StatusCode = resp.StatusCode
				
				// 读取响应体（限制大小）
				body, err := io.ReadAll(io.LimitReader(resp.Body, 1024))
				if err != nil {
					result.Error = err
					result.ElapsedTime = time.Since(requestStart)
					resultsChan <- result
					return
				}
				
				// 设置成功的结果
				result.Body = string(body)
				result.ElapsedTime = time.Since(requestStart)
				
				// 将结果发送到通道
				resultsChan <- result
				
			}(url)
		}
		
		// 等待所有请求完成
		wg.Wait()
		
		// 关闭结果通道，表示没有更多结果
		close(resultsChan)
	}()
	
	// 收集结果
	var results []HTTPResult
	for result := range resultsChan {
		results = append(results, result)
	}
	
	// 计算总耗时
	totalTime := time.Since(startTime)
	
	// 打印结果
	fmt.Println("请求结果:")
	for _, r := range results {
		if r.Error != nil {
			fmt.Printf("× %s - 错误: %v (耗时: %v)\n", 
				r.URL, r.Error, r.ElapsedTime)
		} else {
			fmt.Printf("✓ %s - 状态码: %d (耗时: %v)\n", 
				r.URL, r.StatusCode, r.ElapsedTime)
		}
	}
	
	fmt.Printf("\n总耗时: %v\n", totalTime)
	fmt.Printf("平均每个请求耗时: %v\n", totalTime/time.Duration(len(urls)))
	
	// 计算理论上的串行执行时间
	var serialTime time.Duration
	for _, r := range results {
		serialTime += r.ElapsedTime
	}
	fmt.Printf("理论串行执行总耗时: %v\n", serialTime)
	fmt.Printf("并发+限流节省的时间: %v (%.1f%%)\n", 
		serialTime-totalTime, float64(serialTime-totalTime)/float64(serialTime)*100)
} 