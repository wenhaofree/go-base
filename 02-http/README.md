# Go HTTP 示例

本目录包含了Go语言中HTTP客户端和服务器的各种示例，涵盖了同步和异步请求、并发控制、超时管理等内容。

## 示例文件列表

1. `01-sync-client.go`: 同步HTTP客户端示例，演示了基本的GET、POST请求和自定义HTTP客户端配置。
2. `02-async-goroutines.go`: 使用Goroutines和通道实现异步并行HTTP请求。
3. `03-async-waitgroup.go`: 使用sync.WaitGroup实现异步并行HTTP请求，展示了另一种同步方式。
4. `04-async-rate-limit.go`: 使用信号量(semaphore)模式实现限流的并发HTTP请求。
5. `05-http-server.go`: 基本的HTTP服务器示例，包含路由、中间件、静态文件服务和JSON API。
6. `06-http-client-context.go`: 使用context包控制HTTP请求的超时、取消和并发协调。

## 如何运行示例

每个示例都可以单独运行。使用以下命令运行指定的示例：

```bash
go run 01-sync-client.go
go run 02-async-goroutines.go
# 以此类推...
```

对于HTTP服务器示例，运行后可以使用浏览器或工具如curl来测试：

```bash
# 启动服务器
go run 05-http-server.go

# 然后在另一个终端中测试API
curl http://localhost:8080/api/time
curl http://localhost:8080/api/data/greeting
curl -X POST -H "Content-Type: application/json" -d '{"value":"新问候"}' http://localhost:8080/api/data/greeting
```

## 并发模型比较

Go提供了多种方式来处理并发HTTP请求：

1. **同步顺序请求** (`01-sync-client.go`): 最简单但效率最低的方式，每个请求必须等待前一个请求完成。

2. **Goroutines + Channel** (`02-async-goroutines.go`): 最基本的Go并发模式，为每个请求启动一个goroutine，通过通道收集结果。

3. **WaitGroup** (`03-async-waitgroup.go`): 使用sync.WaitGroup来等待所有goroutine完成，适合需要等待所有请求完成的场景。

4. **限流模式** (`04-async-rate-limit.go`): 使用缓冲通道作为"令牌桶"来限制并发请求数量，避免过多的并发请求导致资源耗尽。

5. **Context控制** (`06-http-client-context.go`): 使用context包进行超时控制、请求取消和请求协调，适合更复杂的请求场景。

## 学习要点

- Go的HTTP请求默认不设置超时，生产环境中应当始终设置适当的超时时间
- 使用`defer resp.Body.Close()`确保响应体被关闭，避免资源泄漏
- 并行请求可以显著提高效率，但需要合理控制并发数量
- 使用context包可以优雅地处理请求的超时和取消
- Go的HTTP服务器使用goroutine处理每个请求，默认即为并发模式

## 参考资料

- [Go标准库文档: net/http](https://golang.org/pkg/net/http/)
- [Go By Example: HTTP客户端](https://gobyexample.com/http-clients)
- [Go By Example: HTTP服务器](https://gobyexample.com/http-servers) 