package screenshot

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

// Options 包含截图的配置选项
type Options struct {
	Width       int
	Height      int
	MobileMode  bool
	WaitTime    time.Duration
	FullPage    bool
	UserAgent   string
	Timeout     time.Duration
	BlockImages bool   // 是否屏蔽图片加载
	BlockJS     bool   // 是否屏蔽JavaScript
	Selector    string // 等待指定元素出现
}

// TimingInfo 包含截图过程的耗时信息
type TimingInfo struct {
	StartTime      time.Time
	BrowserStart   time.Duration
	Navigation     time.Duration
	WaitComplete   time.Duration
	ScreenshotTime time.Duration
	TotalTime      time.Duration
}

// ScreenshotResult 包含截图结果和耗时信息
type ScreenshotResult struct {
	Image     []byte
	Timing    TimingInfo
	Error     error
	URL       string
	Timestamp time.Time
}

// DefaultOptions 返回默认的选项设置
func DefaultOptions() Options {
	return Options{
		Width:       1280,
		Height:      800,
		MobileMode:  false,
		WaitTime:    2 * time.Second,
		FullPage:    false,
		Timeout:     30 * time.Second,
		BlockImages: false,
		BlockJS:     false,
	}
}

// CaptureScreenshot 获取指定URL的网页截图
func CaptureScreenshot(url string, options Options) ([]byte, TimingInfo, error) {
	result := ScreenshotResult{
		URL:       url,
		Timestamp: time.Now(),
		Timing: TimingInfo{
			StartTime: time.Now(),
		},
	}

	if url == "" {
		return nil, result.Timing, errors.New("URL不能为空")
	}

	// 设置超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), options.Timeout)
	defer cancel()

	// 配置浏览器选项，优化启动参数
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-web-security", true),
		chromedp.WindowSize(options.Width, options.Height),
		// 减少内存使用
		chromedp.Flag("disable-software-rasterizer", true),
		chromedp.Flag("disable-sync", true),
		chromedp.Flag("blink-settings", "imagesEnabled=true"),
		// 加速启动
		chromedp.Flag("disable-default-apps", true),
	)

	if options.UserAgent != "" {
		opts = append(opts, chromedp.UserAgent(options.UserAgent))
	}

	if options.MobileMode {
		opts = append(opts, chromedp.Flag("enable-mobile", true))
	}

	// 创建Chrome实例并记录时间
	startBrowser := time.Now()
	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	// 创建新的Chrome实例
	taskCtx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// 设置错误处理器
	chromedp.ListenTarget(taskCtx, func(ev interface{}) {
		// 处理JavaScript对话框
		if dialogEvent, ok := ev.(*page.EventJavascriptDialogOpening); ok {
			fmt.Printf("检测到对话框: %s\n", dialogEvent.Message)
			go func() {
				if err := chromedp.Run(taskCtx, page.HandleJavaScriptDialog(true)); err != nil {
					fmt.Printf("处理JavaScript对话框失败: %v\n", err)
				}
			}()
		}
	})

	result.Timing.BrowserStart = time.Since(startBrowser)

	// 执行截图
	var buf []byte
	var tasks []chromedp.Action

	// 根据选项添加网络拦截任务
	if options.BlockImages || options.BlockJS {
		tasks = append(tasks, chromedp.ActionFunc(func(ctx context.Context) error {
			// 设置网络拦截
			if err := network.Enable().Do(ctx); err != nil {
				return err
			}

			// 设置请求过滤器来阻止特定资源类型
			blockedURLs := []string{}
			
			if options.BlockImages {
				blockedURLs = append(blockedURLs, 
					"*.jpg", "*.jpeg", "*.png", "*.gif", "*.webp", "*.svg", "*.ico")
			}
			
			if options.BlockJS {
				blockedURLs = append(blockedURLs, 
					"*.js", "*.mjs", "*.jsx")
			}

			// 注册请求拦截
			if len(blockedURLs) > 0 {
				return network.SetBlockedURLs(blockedURLs).Do(ctx)
			}

			return nil
		}))
	}

	// 设置设备仿真（如果需要）
	if options.MobileMode {
		tasks = append(tasks, emulation.SetDeviceMetricsOverride(
			int64(options.Width),
			int64(options.Height),
			1.0,
			true,
		))
	}

	// 开始页面导航
	startNav := time.Now()
	tasks = append(tasks, chromedp.Navigate(url))

	// 等待页面加载
	startWait := time.Now()

	// 如果指定了选择器，等待该元素出现
	if options.Selector != "" {
		tasks = append(tasks, chromedp.WaitVisible(options.Selector))
	} else {
		// 否则使用固定等待时间
		tasks = append(tasks, chromedp.Sleep(options.WaitTime))
	}

	// 等待页面加载完成
	tasks = append(tasks, chromedp.ActionFunc(func(ctx context.Context) error {
		result.Timing.Navigation = time.Since(startNav)
		result.Timing.WaitComplete = time.Since(startWait)
		return nil
	}))

	// 截图
	startScreenshot := time.Now()
	if options.FullPage {
		tasks = append(tasks, chromedp.FullScreenshot(&buf, 100))
	} else {
		tasks = append(tasks, chromedp.CaptureScreenshot(&buf))
	}

	if err := chromedp.Run(taskCtx, tasks...); err != nil {
		result.Error = err
		return nil, result.Timing, err
	}

	// 更新耗时信息
	result.Timing.ScreenshotTime = time.Since(startScreenshot)
	result.Timing.TotalTime = time.Since(result.Timing.StartTime)
	result.Image = buf

	return buf, result.Timing, nil
} 