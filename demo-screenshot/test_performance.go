package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fuwenhao/go-base/demo-screenshot/screenshot"
)

// 测试配置
type TestCase struct {
	Name        string
	URL         string
	BlockImages bool
	BlockJS     bool
	UseSelector bool
	Selector    string
}

func main() {
	// 检查参数
	if len(os.Args) < 2 {
		fmt.Println("使用方法: go run test_performance.go <URL>")
		os.Exit(1)
	}

	url := os.Args[1]
	fmt.Printf("测试网站: %s\n\n", url)

	// 定义测试用例
	testCases := []TestCase{
		{
			Name:        "标准模式",
			URL:         url,
			BlockImages: false,
			BlockJS:     false,
			UseSelector: false,
		},
		{
			Name:        "屏蔽图片",
			URL:         url,
			BlockImages: true,
			BlockJS:     false,
			UseSelector: false,
		},
		{
			Name:        "屏蔽JS",
			URL:         url,
			BlockImages: false,
			BlockJS:     true,
			UseSelector: false,
		},
		{
			Name:        "屏蔽图片和JS",
			URL:         url,
			BlockImages: true,
			BlockJS:     true,
			UseSelector: false,
		},
	}

	// 存储结果
	var results []struct {
		Name   string
		Timing screenshot.TimingInfo
	}

	// 运行测试
	for _, tc := range testCases {
		fmt.Printf("测试案例: %s\n", tc.Name)
		
		options := screenshot.DefaultOptions()
		options.BlockImages = tc.BlockImages
		options.BlockJS = tc.BlockJS
		
		if tc.UseSelector {
			options.Selector = tc.Selector
		}
		
		testStart := time.Now()
		
		_, timing, err := screenshot.CaptureScreenshot(tc.URL, options)
		if err != nil {
			log.Printf("测试失败 [%s]: %v\n", tc.Name, err)
			continue
		}
		
		results = append(results, struct {
			Name   string
			Timing screenshot.TimingInfo
		}{
			Name:   tc.Name,
			Timing: timing,
		})
		
		fmt.Printf("- 总耗时: %.2f 秒\n", timing.TotalTime.Seconds())
		fmt.Printf("- 浏览器启动: %.2f 秒\n", timing.BrowserStart.Seconds())
		fmt.Printf("- 页面导航: %.2f 秒\n", timing.Navigation.Seconds())
		fmt.Printf("- 页面等待: %.2f 秒\n", timing.WaitComplete.Seconds())
		fmt.Printf("- 截图操作: %.2f 秒\n", timing.ScreenshotTime.Seconds())
		fmt.Printf("- 实际执行时间: %.2f 秒\n", time.Since(testStart).Seconds())
		fmt.Println()
	}

	// 输出对比结果
	if len(results) > 1 {
		fmt.Println("\n=== 性能优化对比 ===")
		baseline := results[0].Timing.TotalTime.Seconds()
		
		for i, r := range results {
			if i == 0 {
				fmt.Printf("%s: %.2f 秒 (基准)\n", r.Name, r.Timing.TotalTime.Seconds())
			} else {
				improved := baseline - r.Timing.TotalTime.Seconds()
				percent := (improved / baseline) * 100
				fmt.Printf("%s: %.2f 秒 (提升 %.1f%%, 节省 %.2f 秒)\n", 
					r.Name, 
					r.Timing.TotalTime.Seconds(),
					percent,
					improved,
				)
			}
		}
	}
} 