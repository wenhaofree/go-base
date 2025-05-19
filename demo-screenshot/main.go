package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/fuwenhao/go-base/demo-screenshot/screenshot"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("使用方法: go run main.go <URL> [输出文件名.png] [选项]")
		fmt.Println("选项:")
		fmt.Println("  --width=数值     : 设置截图宽度")
		fmt.Println("  --height=数值    : 设置截图高度")
		fmt.Println("  --full=true/false: 是否全页面截图")
		fmt.Println("  --mobile=true/false: 是否使用移动设备模式")
		fmt.Println("  --wait=数值      : 等待时间(秒)")
		fmt.Println("  --block-images=true/false: 是否屏蔽图片加载")
		fmt.Println("  --block-js=true/false: 是否屏蔽JavaScript")
		fmt.Println("  --selector=CSS选择器: 等待指定元素出现后截图")
		os.Exit(1)
	}

	url := os.Args[1]
	outputFile := "screenshot.png"
	if len(os.Args) >= 3 && !strings.HasPrefix(os.Args[2], "--") {
		outputFile = os.Args[2]
	}

	// 使用默认选项
	options := screenshot.DefaultOptions()
	
	// 解析命令行参数
	for i := 2; i < len(os.Args); i++ {
		arg := os.Args[i]
		if !strings.HasPrefix(arg, "--") {
			continue
		}
		
		parts := strings.SplitN(strings.TrimPrefix(arg, "--"), "=", 2)
		if len(parts) != 2 {
			continue
		}
		
		key, value := parts[0], parts[1]
		
		switch key {
		case "width":
			if w, err := strconv.Atoi(value); err == nil && w > 0 {
				options.Width = w
			}
		case "height":
			if h, err := strconv.Atoi(value); err == nil && h > 0 {
				options.Height = h
			}
		case "full":
			options.FullPage = (value == "true" || value == "1")
		case "mobile":
			options.MobileMode = (value == "true" || value == "1")
		case "wait":
			if w, err := strconv.Atoi(value); err == nil && w > 0 {
				options.WaitTime = screenshot.ParseDuration(w, "s")
			}
		case "block-images":
			options.BlockImages = (value == "true" || value == "1")
		case "block-js":
			options.BlockJS = (value == "true" || value == "1")
		case "selector":
			options.Selector = value
		}
	}

	// 获取网页截图
	fmt.Printf("开始截图: %s\n", url)
	
	buf, timing, err := screenshot.CaptureScreenshot(url, options)
	if err != nil {
		log.Fatalf("截图失败: %v", err)
	}

	// 保存截图到文件
	if err := os.WriteFile(outputFile, buf, 0644); err != nil {
		log.Fatalf("无法保存截图: %v", err)
	}

	// 显示耗时统计
	fmt.Printf("成功截图网页 %s 并保存到 %s\n", url, outputFile)
	fmt.Printf("\n=== 耗时统计 ===\n")
	fmt.Printf("启动浏览器: %.2f 秒\n", timing.BrowserStart.Seconds())
	fmt.Printf("页面导航: %.2f 秒\n", timing.Navigation.Seconds())
	fmt.Printf("页面等待: %.2f 秒\n", timing.WaitComplete.Seconds())
	fmt.Printf("截图操作: %.2f 秒\n", timing.ScreenshotTime.Seconds())
	fmt.Printf("总耗时: %.2f 秒\n", timing.TotalTime.Seconds())
	
	// 显示优化统计
	if options.BlockImages || options.BlockJS {
		fmt.Printf("\n=== 优化策略 ===\n")
		if options.BlockImages {
			fmt.Println("- 已屏蔽图片加载")
		}
		if options.BlockJS {
			fmt.Println("- 已屏蔽JavaScript")
		}
	}
} 