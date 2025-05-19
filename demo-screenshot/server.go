package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/fuwenhao/go-base/demo-screenshot/screenshot"
)

// 响应结构
type Response struct {
	Success bool                   `json:"success"`
	Error   string                 `json:"error,omitempty"`
	Timing  map[string]interface{} `json:"timing,omitempty"`
}

func main() {
	// 设置HTTP路由
	http.HandleFunc("/screenshot", handleScreenshot)
	http.HandleFunc("/screenshot/info", handleScreenshotInfo)
	
	// 启动HTTP服务器
	port := 8080
	fmt.Printf("截图服务启动于 http://localhost:%d\n", port)
	fmt.Printf("- 截图API: http://localhost:%d/screenshot?url=网址\n", port)
	fmt.Printf("- 信息API: http://localhost:%d/screenshot/info?url=网址\n", port)
	
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

func handleScreenshot(w http.ResponseWriter, r *http.Request) {
	// 获取URL参数
	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "请提供有效的URL参数", http.StatusBadRequest)
		return
	}

	// 设置截图选项
	options := parseOptions(r)
	
	// 记录开始时间
	startTime := time.Now()
	
	// 捕获截图
	buf, timing, err := screenshot.CaptureScreenshot(url, options)
	if err != nil {
		log.Printf("无法获取截图: %v", err)
		http.Error(w, fmt.Sprintf("截图失败: %v", err), http.StatusInternalServerError)
		return
	}

	// 添加耗时统计到响应头
	w.Header().Set("X-Timing-Total", fmt.Sprintf("%.2fs", timing.TotalTime.Seconds()))
	w.Header().Set("X-Timing-Browser", fmt.Sprintf("%.2fs", timing.BrowserStart.Seconds()))
	w.Header().Set("X-Timing-Navigation", fmt.Sprintf("%.2fs", timing.Navigation.Seconds()))
	w.Header().Set("X-Timing-Screenshot", fmt.Sprintf("%.2fs", timing.ScreenshotTime.Seconds()))
	
	// 设置响应头并返回图片
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=screenshot-%d.png", time.Now().Unix()))
	w.Write(buf)
	
	// 记录请求完成信息
	log.Printf("截图完成: %s (耗时: %.2fs, 大小: %d KB)", 
		url, 
		time.Since(startTime).Seconds(),
		len(buf)/1024,
	)
}

func handleScreenshotInfo(w http.ResponseWriter, r *http.Request) {
	// 获取URL参数
	url := r.URL.Query().Get("url")
	if url == "" {
		sendJSONError(w, "请提供有效的URL参数", http.StatusBadRequest)
		return
	}

	// 设置截图选项
	options := parseOptions(r)
	
	// 捕获截图
	_, timing, err := screenshot.CaptureScreenshot(url, options)
	if err != nil {
		log.Printf("无法获取截图信息: %v", err)
		sendJSONError(w, fmt.Sprintf("获取信息失败: %v", err), http.StatusInternalServerError)
		return
	}

	// 返回耗时统计信息
	w.Header().Set("Content-Type", "application/json")
	response := Response{
		Success: true,
		Timing:  screenshot.TimingToMap(timing),
	}
	
	json.NewEncoder(w).Encode(response)
}

func parseOptions(r *http.Request) screenshot.Options {
	options := screenshot.DefaultOptions()
	
	// 处理附加参数
	if width := r.URL.Query().Get("width"); width != "" {
		if w, err := strconv.Atoi(width); err == nil && w > 0 {
			options.Width = w
		}
	}
	
	if height := r.URL.Query().Get("height"); height != "" {
		if h, err := strconv.Atoi(height); err == nil && h > 0 {
			options.Height = h
		}
	}
	
	if fullPage := r.URL.Query().Get("full"); fullPage == "1" || fullPage == "true" {
		options.FullPage = true
	}
	
	if mobile := r.URL.Query().Get("mobile"); mobile == "1" || mobile == "true" {
		options.MobileMode = true
	}
	
	if userAgent := r.URL.Query().Get("ua"); userAgent != "" {
		options.UserAgent = userAgent
	}
	
	if wait := r.URL.Query().Get("wait"); wait != "" {
		if waitTime, err := strconv.Atoi(wait); err == nil && waitTime > 0 {
			options.WaitTime = time.Duration(waitTime) * time.Second
		}
	}
	
	if blockImages := r.URL.Query().Get("block-images"); blockImages == "1" || blockImages == "true" {
		options.BlockImages = true
	}
	
	if blockJS := r.URL.Query().Get("block-js"); blockJS == "1" || blockJS == "true" {
		options.BlockJS = true
	}
	
	if selector := r.URL.Query().Get("selector"); selector != "" {
		options.Selector = selector
	}
	
	return options
}

func sendJSONError(w http.ResponseWriter, errorMsg string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	response := Response{
		Success: false,
		Error:   errorMsg,
	}
	
	json.NewEncoder(w).Encode(response)
} 