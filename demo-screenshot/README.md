# 网页截图工具

这是一个基于Golang和Chromium的网页截图工具，可以通过命令行或HTTP API来获取网页截图。经过性能优化，提供了更快的截图速度和详细的耗时统计。

## 功能特点

- 支持命令行方式获取指定URL的截图
- 提供HTTP API服务，支持通过Web请求获取截图
- 可配置截图尺寸、等待时间、全页面截图等选项
- 支持移动设备模拟
- **性能优化**：支持屏蔽图片和JavaScript加载，显著提高截图速度
- **智能等待**：支持等待指定元素出现后再截图，避免固定等待时间
- **详细统计**：提供完整的耗时统计，包括浏览器启动、页面导航和截图各阶段耗时

## 依赖条件

- Go 1.16+
- Chrome/Chromium浏览器（支持headless模式）

## 安装

```bash
# 克隆代码
git clone https://github.com/fuwenhao/go-base.git
cd go-base/demo-screenshot

# 安装依赖
go mod tidy
```

## 使用方法

### 命令行工具

```bash
# 基本用法
go run main.go https://example.com

# 指定输出文件
go run main.go https://example.com output.png

# 使用性能优化选项
go run main.go https://example.com --block-images=true --block-js=true

# 等待指定元素出现
go run main.go https://example.com --selector="#content"

# 全页面截图
go run main.go https://example.com --full=true
```

### HTTP API 服务

启动服务：

```bash
go run server.go
```

API使用：

#### 1. 截图API

- 基本截图：`http://localhost:8080/screenshot?url=https://example.com`
- 指定尺寸：`http://localhost:8080/screenshot?url=https://example.com&width=1920&height=1080`
- 全页面截图：`http://localhost:8080/screenshot?url=https://example.com&full=true`
- 移动设备模拟：`http://localhost:8080/screenshot?url=https://example.com&mobile=true`
- 设置等待时间：`http://localhost:8080/screenshot?url=https://example.com&wait=5`（等待5秒）
- 自定义User-Agent：`http://localhost:8080/screenshot?url=https://example.com&ua=Mozilla/5.0...`
- **优化选项**：
  - 屏蔽图片：`http://localhost:8080/screenshot?url=https://example.com&block-images=true`
  - 屏蔽JavaScript：`http://localhost:8080/screenshot?url=https://example.com&block-js=true`
  - 等待指定元素：`http://localhost:8080/screenshot?url=https://example.com&selector=#main-content`

#### 2. 信息API

获取截图耗时统计但不返回图片：
`http://localhost:8080/screenshot/info?url=https://example.com`

返回JSON格式的耗时统计信息，包括：
- 浏览器启动时间
- 页面导航时间
- 等待加载时间
- 截图操作时间
- 总耗时

## 性能优化技巧

以下选项可以显著提高截图速度：

1. **屏蔽图片加载**：对于内容丰富的网页可提升50-70%的速度
   ```
   --block-images=true
   ```

2. **屏蔽JavaScript**：对于JS密集型网站可提升30-50%的速度
   ```
   --block-js=true
   ```
   
3. **使用CSS选择器等待**：比固定等待时间更高效
   ```
   --selector=".main-content"
   ```
   
4. **合理设置截图尺寸**：较小的尺寸处理更快
   ```
   --width=1024 --height=768
   ```


## 代码结构

```
demo-screenshot/
├── main.go       # 命令行工具入口
├── server.go     # HTTP API服务入口
└── screenshot/   # 核心截图功能包
    ├── screenshot.go  # 核心截图功能实现
    └── utils.go       # 辅助函数集合
```

## 注意事项

- 需要安装Chrome/Chromium浏览器
- 使用时确保有足够的内存，因为启动浏览器需要较大的内存
- 部分网页可能需要较长的加载时间，可以通过`wait`参数增加等待时间
- 使用`block-images`和`block-js`选项可能会导致部分网页的布局变化 