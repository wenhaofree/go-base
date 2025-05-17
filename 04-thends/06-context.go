// 06-context.go 展示了使用context包控制goroutine生命周期
package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func main() {
	fmt.Println("=== Go Context示例 ===")

	// 1. 基本context和取消操作
	fmt.Println("\n1. 基本context和取消操作")
	basicContextExample()

	// 2. 带超时的context
	fmt.Println("\n2. 带超时的context")
	timeoutContextExample()

	// 3. 带截止时间的context
	fmt.Println("\n3. 带截止时间的context")
	deadlineContextExample()

	// 4. 带值的context
	fmt.Println("\n4. 带值的context")
	valueContextExample()

	// 5. context与HTTP请求
	fmt.Println("\n5. context与HTTP请求")
	httpContextExample()

	fmt.Println("所有示例完成")
}

// basicContextExample 展示基本的context使用和取消操作
func basicContextExample() {
	// 创建带取消功能的context
	ctx, cancel := context.WithCancel(context.Background())

	// 在一个goroutine中运行任务，并监听context取消信号
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("任务收到取消信号，退出...")
				return
			default:
				fmt.Println("任务正在执行...")
				time.Sleep(500 * time.Millisecond)
			}
		}
	}(ctx)

	// 让任务执行几次
	time.Sleep(1500 * time.Millisecond)

	// 取消context
	fmt.Println("正在取消任务...")
	cancel()

	// 给时间让goroutine响应取消
	time.Sleep(500 * time.Millisecond)
}

// timeoutContextExample 展示带超时的context
func timeoutContextExample() {
	// 创建一个2秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel() // 记得调用cancel以释放资源

	// 在goroutine中执行长时间任务
	go func(ctx context.Context) {
		// 模拟一个3秒的任务
		taskDuration := 3 * time.Second
		fmt.Printf("开始执行一个需要%v的任务...\n", taskDuration)

		select {
		case <-time.After(taskDuration):
			fmt.Println("任务完成!")
		case <-ctx.Done():
			fmt.Printf("任务被中断: %v\n", ctx.Err())
		}
	}(ctx)

	// 等待足够的时间让goroutine执行和超时
	time.Sleep(3 * time.Second)
}

// deadlineContextExample 展示带截止时间的context
func deadlineContextExample() {
	// 设置2秒后的截止时间
	deadline := time.Now().Add(2 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel() // 记得调用cancel以释放资源

	fmt.Printf("设置截止时间: %v\n", deadline.Format("15:04:05.000"))

	// 创建一个显示倒计时的goroutine
	go func(ctx context.Context) {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				timeLeft := time.Until(deadline)
				if timeLeft > 0 {
					fmt.Printf("距离截止还有: %.1f秒\n", timeLeft.Seconds())
				}
			case <-ctx.Done():
				fmt.Printf("截止时间到达: %v\n", ctx.Err())
				return
			}
		}
	}(ctx)

	// 等待足够的时间让goroutine执行和超时
	time.Sleep(3 * time.Second)
}

// valueContextExample 展示带值的context
func valueContextExample() {
	// 创建基础context
	rootCtx := context.Background()

	// 创建第一级带值的context
	userCtx := context.WithValue(rootCtx, "user_id", "12345")

	// 创建第二级带值的context
	sessionCtx := context.WithValue(userCtx, "session_id", "abc-xyz-789")

	// 模拟把context传递给一个处理函数
	processRequest(sessionCtx)
}

// processRequest 处理带有user_id和session_id的请求
func processRequest(ctx context.Context) {
	// 从context中获取值
	userID, ok := ctx.Value("user_id").(string)
	if !ok {
		fmt.Println("用户ID不存在或类型错误")
		return
	}

	sessionID, ok := ctx.Value("session_id").(string)
	if !ok {
		fmt.Println("会话ID不存在或类型错误")
		return
	}

	fmt.Printf("处理用户ID为'%s'，会话ID为'%s'的请求\n", userID, sessionID)

	// 模拟进一步处理，传递context给子函数
	validateSession(ctx)
}

// validateSession 验证会话有效性
func validateSession(ctx context.Context) {
	// 从context获取会话ID
	sessionID, ok := ctx.Value("session_id").(string)
	if !ok {
		fmt.Println("子函数: 会话ID不存在或类型错误")
		return
	}

	fmt.Printf("子函数: 验证会话ID '%s' 有效性\n", sessionID)
	// 实际应用中这里可能会查询数据库等
}

// httpContextExample 展示context与HTTP请求的结合使用
func httpContextExample() {
	// 创建一个简单的HTTP客户端
	client := &http.Client{}

	// 创建一个带有1秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// 创建HTTP请求
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/users/github", nil)
	if err != nil {
		fmt.Printf("创建请求失败: %v\n", err)
		return
	}

	fmt.Println("发送HTTP请求到GitHub API (带1秒超时)...")
	
	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// 输出响应状态
	fmt.Printf("收到响应: %s\n", resp.Status)
	
	// 在实际应用中，这里会处理响应内容
	fmt.Println("成功完成HTTP请求")
} 