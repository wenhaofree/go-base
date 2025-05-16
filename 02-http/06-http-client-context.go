// 06-http-client-context.go 展示了Go中使用context控制HTTP请求
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

// 定义API响应结构
type APIResponse struct {
	Args     map[string]string `json:"args"`
	Headers  map[string]string `json:"headers"`
	Origin   string            `json:"origin"`
	URL      string            `json:"url"`
	Data     string            `json:"data,omitempty"`
	Json     interface{}       `json:"json,omitempty"`
	Response string
	Err      error
}

func main() {
	fmt.Println("开始HTTP客户端Context控制示例")
	
	// 创建一个可取消的根context
	rootCtx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc() // 确保所有context最终被取消
	
	// 设置信号处理，以便在接收到中断信号时取消所有请求
	setupSignalHandler(cancelFunc)
	
	// 示例1: 带超时的HTTP请求
	fmt.Println("\n===== 示例1: 带超时的HTTP请求 =====")
	timeoutExample(rootCtx)
	
	// 示例2: 可取消的HTTP请求
	fmt.Println("\n===== 示例2: 可取消的HTTP请求 =====")
	cancellationExample(rootCtx)
	
	// 示例3: 并发请求并等待第一个成功结果
	fmt.Println("\n===== 示例3: 并发请求首个成功结果 =====")
	firstSuccessExample(rootCtx)
	
	fmt.Println("\n所有示例已完成")
}

// 设置信号处理，以便在程序接收到中断信号时取消所有请求
func setupSignalHandler(cancelFunc context.CancelFunc) {
	signalChan := make(chan os.Signal, 1)
	
	// 注册要监听的信号
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	
	// 启动goroutine来处理信号
	go func() {
		<-signalChan
		fmt.Println("\n收到中断信号，取消所有请求...")
		cancelFunc()
		fmt.Println("已取消所有请求，正在退出...")
		os.Exit(0)
	}()
}

// 示例1: 带超时的HTTP请求
func timeoutExample(parentCtx context.Context) {
	// 创建一个3秒超时的context
	ctx, cancel := context.WithTimeout(parentCtx, 3*time.Second)
	defer cancel()
	
	fmt.Println("发起一个带3秒超时的请求到 httpbin.org/delay/5 (将会超时)")
	
	// 创建请求
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://httpbin.org/delay/5", nil)
	if err != nil {
		fmt.Printf("创建请求失败: %v\n", err)
		return
	}
	
	// 发送请求并计时
	startTime := time.Now()
	resp, err := http.DefaultClient.Do(req)
	duration := time.Since(startTime)
	
	// 检查结果
	if err != nil {
		fmt.Printf("请求失败 (耗时 %.2f 秒): %v\n", duration.Seconds(), err)
		// 检查是否是因为context超时导致的
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Println("确认失败原因: context超时")
		}
	} else {
		defer resp.Body.Close()
		fmt.Printf("请求成功 (耗时 %.2f 秒): 状态码 %d\n", duration.Seconds(), resp.StatusCode)
	}
}

// 示例2: 可取消的HTTP请求
func cancellationExample(parentCtx context.Context) {
	// 创建一个可取消的context
	ctx, cancel := context.WithCancel(parentCtx)
	
	fmt.Println("发起请求到 httpbin.org/delay/3，并在1秒后手动取消")
	
	// 创建请求
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://httpbin.org/delay/3", nil)
	if err != nil {
		fmt.Printf("创建请求失败: %v\n", err)
		return
	}
	
	// 在1秒后自动取消请求
	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("手动取消请求...")
		cancel()
	}()
	
	// 发送请求并计时
	startTime := time.Now()
	resp, err := http.DefaultClient.Do(req)
	duration := time.Since(startTime)
	
	// 检查结果
	if err != nil {
		fmt.Printf("请求失败 (耗时 %.2f 秒): %v\n", duration.Seconds(), err)
		// 检查是否是因为context取消导致的
		if ctx.Err() == context.Canceled {
			fmt.Println("确认失败原因: context被取消")
		}
	} else {
		defer resp.Body.Close()
		fmt.Printf("请求成功 (耗时 %.2f 秒): 状态码 %d\n", duration.Seconds(), resp.StatusCode)
	}
}

// 示例3: 并发请求并等待第一个成功结果
func firstSuccessExample(parentCtx context.Context) {
	// 创建一个带超时的context，防止所有请求都失败时无限等待
	ctx, cancel := context.WithTimeout(parentCtx, 10*time.Second)
	defer cancel()
	
	// 准备多个服务器URL
	urls := []string{
		"https://httpbin.org/delay/2",
		"https://httpbin.org/delay/3", 
		"https://httpbin.org/delay/1", // 这个应该是最快的
		"https://httpbin.org/status/500", // 这个会返回错误
	}
	
	fmt.Printf("并发请求 %d 个URL，返回第一个成功的结果\n", len(urls))
	
	// 创建一个通道接收结果
	resultChan := make(chan APIResponse, len(urls))
	
	// 启动所有请求
	var wg sync.WaitGroup
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			
			// 为每个请求创建一个可取消的context
			reqCtx, reqCancel := context.WithCancel(ctx)
			defer reqCancel()
			
			startTime := time.Now()
			
			// 创建请求
			req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, url, nil)
			if err != nil {
				resultChan <- APIResponse{URL: url, Err: err}
				return
			}
			
			// 发送请求
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				resultChan <- APIResponse{URL: url, Err: err}
				return
			}
			defer resp.Body.Close()
			
			// 检查响应状态码
			if resp.StatusCode >= 400 {
				resultChan <- APIResponse{
					URL: url, 
					Err: fmt.Errorf("HTTP错误: %d", resp.StatusCode),
				}
				return
			}
			
			// 读取响应体
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				resultChan <- APIResponse{URL: url, Err: err}
				return
			}
			
			// 解析JSON响应
			var apiResp APIResponse
			err = json.Unmarshal(body, &apiResp)
			if err != nil {
				apiResp = APIResponse{
					Response: string(body),
					Err:      nil,
				}
			}
			
			apiResp.URL = url
			fmt.Printf("完成请求: %s (耗时: %.2f 秒)\n", url, time.Since(startTime).Seconds())
			
			// 发送结果
			resultChan <- apiResp
			
		}(url)
	}
	
	// 启动一个goroutine来关闭结果通道
	go func() {
		wg.Wait()
		close(resultChan)
	}()
	
	// 等待第一个成功的结果
	startTime := time.Now()
	var firstSuccessResult *APIResponse
	
	for result := range resultChan {
		if result.Err == nil {
			firstSuccessResult = &result
			// 取消所有其他请求
			cancel()
			break
		} else {
			fmt.Printf("请求 %s 失败: %v\n", result.URL, result.Err)
		}
	}
	
	// 打印结果
	if firstSuccessResult != nil {
		fmt.Printf("\n获得首个成功的结果 (总耗时: %.2f 秒):\n", time.Since(startTime).Seconds())
		fmt.Printf("URL: %s\n", firstSuccessResult.URL)
		
		// 显示部分结果
		responseExcerpt := firstSuccessResult.Response
		if responseExcerpt == "" {
			responseBytes, _ := json.MarshalIndent(firstSuccessResult, "", "  ")
			responseExcerpt = string(responseBytes)
		}
		
		if len(responseExcerpt) > 500 {
			responseExcerpt = responseExcerpt[:500] + "...[截断]"
		}
		fmt.Printf("响应: %s\n", responseExcerpt)
	} else {
		fmt.Printf("\n没有获得成功的结果 (总耗时: %.2f 秒)\n", time.Since(startTime).Seconds())
	}
}

// 发送带有自定义请求体的POST请求
func sendPostRequest(ctx context.Context, url string, data map[string]interface{}) (*APIResponse, error) {
	// 将数据转换为JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("JSON编码失败: %w", err)
	}
	
	// 创建请求
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		strings.NewReader(string(jsonData)),
	)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	
	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Go-HTTP-Client/Context-Example")
	
	// 发送请求
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()
	
	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}
	
	// 解析响应
	var result APIResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return &APIResponse{
			Response: string(body),
		}, nil
	}
	
	return &result, nil
} 