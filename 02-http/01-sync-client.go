// 01-sync-client.go 展示了Go中同步HTTP客户端请求的基本用法
package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func main() {
	fmt.Println("开始同步HTTP请求示例")
	
	// 记录开始时间
	startTime := time.Now()

	// 简单的GET请求
	// http.Get是一个便捷的方法，内部使用默认的http.Client
	fmt.Println("\n1. 基本GET请求：")
	resp, err := http.Get("https://httpbin.org/get")
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return
	}
	// 务必关闭响应体，避免资源泄漏
	defer resp.Body.Close()
	
	// 读取并打印响应状态
	fmt.Printf("状态码: %d\n", resp.StatusCode)
	fmt.Printf("状态: %s\n", resp.Status)
	
	// 读取并打印响应体（最多1024字节）
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1024))
	if err != nil {
		fmt.Printf("读取响应体失败: %v\n", err)
		return
	}
	fmt.Printf("响应体: %s\n", body)

	// 使用自定义客户端配置进行GET请求
	fmt.Println("\n2. 使用自定义HTTP客户端：")
	// 创建一个自定义的HTTP客户端，设置超时
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	// 使用自定义客户端发起请求
	resp2, err := client.Get("https://httpbin.org/delay/1")
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return
	}
	defer resp2.Body.Close()
	
	fmt.Printf("状态: %s\n", resp2.Status)
	body2, _ := io.ReadAll(io.LimitReader(resp2.Body, 1024))
	fmt.Printf("响应体: %s\n", body2)

	// 发起POST请求示例
	fmt.Println("\n3. POST请求示例：")
	resp3, err := http.Post(
		"https://httpbin.org/post",
		"application/json",
		// 使用字符串作为请求体
		// 也可以使用bytes.NewReader或strings.NewReader传递请求体
		// bytes.NewBuffer也是常见选择
		// 大型请求可以使用io.Pipe创建一个reader/writer对
		strings.NewReader(`{"name":"张三", "age":30}`),
	)
	if err != nil {
		fmt.Printf("POST请求失败: %v\n", err)
		return
	}
	defer resp3.Body.Close()
	
	fmt.Printf("状态: %s\n", resp3.Status)
	body3, _ := io.ReadAll(io.LimitReader(resp3.Body, 1024))
	fmt.Printf("响应体: %s\n", body3)

	// 使用http.NewRequest创建自定义请求
	fmt.Println("\n4. 自定义请求头：")
	// 创建一个新的请求
	req, err := http.NewRequest("GET", "https://httpbin.org/headers", nil)
	if err != nil {
		fmt.Printf("创建请求失败: %v\n", err)
		return
	}
	
	// 添加自定义请求头
	req.Header.Add("User-Agent", "GoExample-Client/1.0")
	req.Header.Add("X-Custom-Header", "自定义值")
	
	// 发送请求
	resp4, err := client.Do(req)
	if err != nil {
		fmt.Printf("发送请求失败: %v\n", err)
		return
	}
	defer resp4.Body.Close()
	
	fmt.Printf("状态: %s\n", resp4.Status)
	body4, _ := io.ReadAll(io.LimitReader(resp4.Body, 1024))
	fmt.Printf("响应体: %s\n", body4)

	// 计算总耗时
	elapsedTime := time.Since(startTime)
	fmt.Printf("\n所有请求总耗时: %v\n", elapsedTime)
} 