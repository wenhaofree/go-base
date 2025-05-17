// 03-concurrent-http.go 展示了使用goroutine并发获取多个URL内容
package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

func main() {
	fmt.Println("=== Go并发HTTP请求示例 ===")

	// 准备一些要访问的URL
	urls := []string{
		"https://www.baidu.com",
		"https://www.qq.com",
		"https://www.163.com",
		"https://cn.bing.com",
		"https://www.zhihu.com",
	}

	// 1. 顺序获取URL内容
	fmt.Println("\n1. 顺序获取URL内容")
	sequentialFetch(urls)

	// 2. 使用goroutine和WaitGroup并发获取
	fmt.Println("\n2. 使用goroutine和WaitGroup并发获取")
	concurrentFetchWithWaitGroup(urls)

	// 3. 使用goroutine和channel获取结果
	fmt.Println("\n3. 使用goroutine和channel获取结果")
	concurrentFetchWithChannel(urls)

	fmt.Println("所有示例完成")
}

// sequentialFetch 按顺序获取每个URL的内容
func sequentialFetch(urls []string) {
	start := time.Now()

	for i, url := range urls {
		fmt.Printf("开始获取URL(%d): %s\n", i+1, url)
		
		// 获取URL内容
		body, size, err := fetchURL(url)
		if err != nil {
			fmt.Printf("获取URL(%d)失败: %v\n", i+1, err)
			continue
		}

		fmt.Printf("URL(%d): %s - 获取成功, 大小: %d 字节, 前100字符: %.100s\n",
			i+1, url, size, body)
	}

	elapsed := time.Since(start)
	fmt.Printf("顺序获取完成，总耗时: %v\n", elapsed)
}

// concurrentFetchWithWaitGroup 使用goroutine和WaitGroup并发获取URL内容
func concurrentFetchWithWaitGroup(urls []string) {
	start := time.Now()
	var wg sync.WaitGroup

	// 为每个URL启动一个goroutine
	for i, url := range urls {
		wg.Add(1)
		
		go func(index int, fetchURL string) {
			defer wg.Done()
			
			fmt.Printf("开始获取URL(%d): %s\n", index+1, fetchURL)
			
			// 获取URL内容
			body, size, err := fetchURL(fetchURL)
			if err != nil {
				fmt.Printf("获取URL(%d)失败: %v\n", index+1, err)
				return
			}

			fmt.Printf("URL(%d): %s - 获取成功, 大小: %d 字节, 前100字符: %.100s\n",
				index+1, fetchURL, size, body)
		}(i, url)
	}

	// 等待所有请求完成
	wg.Wait()
	
	elapsed := time.Since(start)
	fmt.Printf("并发获取(WaitGroup)完成，总耗时: %v\n", elapsed)
}

// concurrentFetchWithChannel 使用goroutine和channel并发获取URL内容
func concurrentFetchWithChannel(urls []string) {
	start := time.Now()
	
	// 创建结果通道
	type Result struct {
		URL      string
		Content  string
		Size     int
		Error    error
		Index    int
	}
	
	resultChan := make(chan Result, len(urls))
	
	// 为每个URL启动一个goroutine
	for i, url := range urls {
		go func(index int, fetchURL string) {
			fmt.Printf("开始获取URL(%d): %s\n", index+1, fetchURL)
			
			// 获取URL内容
			body, size, err := fetchURL(fetchURL)
			
			// 将结果发送到通道
			resultChan <- Result{
				URL:     fetchURL,
				Content: body,
				Size:    size,
				Error:   err,
				Index:   index,
			}
		}(i, url)
	}
	
	// 收集结果
	for i := 0; i < len(urls); i++ {
		result := <-resultChan
		if result.Error != nil {
			fmt.Printf("获取URL(%d)失败: %v\n", result.Index+1, result.Error)
			continue
		}

		fmt.Printf("URL(%d): %s - 获取成功, 大小: %d 字节, 前100字符: %.100s\n",
			result.Index+1, result.URL, result.Size, result.Content)
	}
	
	elapsed := time.Since(start)
	fmt.Printf("并发获取(Channel)完成，总耗时: %v\n", elapsed)
}

// fetchURL 获取指定URL的内容
func fetchURL(url string) (string, int, error) {
	// 发起HTTP GET请求
	resp, err := http.Get(url)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()
	
	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}
	
	return string(body), len(body), nil
} 