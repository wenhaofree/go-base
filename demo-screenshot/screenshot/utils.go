package screenshot

import (
	"time"
)

// ParseDuration 将数值和单位转换为时间间隔
func ParseDuration(value int, unit string) time.Duration {
	switch unit {
	case "ms":
		return time.Duration(value) * time.Millisecond
	case "s":
		return time.Duration(value) * time.Second
	case "m":
		return time.Duration(value) * time.Minute
	default:
		return time.Duration(value) * time.Second
	}
}

// GetCurrentTime 获取当前时间
func GetCurrentTime() time.Time {
	return time.Now()
}

// FormatDuration 将时间间隔格式化为字符串
func FormatDuration(d time.Duration) string {
	if d < time.Millisecond {
		return "< 1ms"
	} else if d < time.Second {
		return d.Round(time.Millisecond).String()
	} else {
		return d.Round(10 * time.Millisecond).String()
	}
}

// TimingToMap 将耗时信息转换为字典
func TimingToMap(timing TimingInfo) map[string]interface{} {
	return map[string]interface{}{
		"browser_start": timing.BrowserStart.Seconds(),
		"navigation":    timing.Navigation.Seconds(),
		"wait_complete": timing.WaitComplete.Seconds(),
		"screenshot":    timing.ScreenshotTime.Seconds(),
		"total":         timing.TotalTime.Seconds(),
	}
} 