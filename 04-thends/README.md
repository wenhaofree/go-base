# Go并发编程示例

本目录包含了Go语言中并发编程的各种示例，展示了goroutines、channels以及其他并发控制机制的使用方法。

## 文件列表

1. **01-basic-goroutine.go**：展示goroutine的基本用法
   - 简单goroutine的创建和执行
   - 匿名goroutine
   - 在循环中启动多个goroutine

2. **02-waitgroup.go**：演示使用WaitGroup同步多个goroutine
   - 基本WaitGroup使用
   - 处理多个任务并等待完成
   - 在函数中使用WaitGroup

3. **03-concurrent-http.go**：展示并发获取多个URL内容
   - 顺序获取URL内容（无并发）
   - 使用goroutine和WaitGroup并发获取
   - 使用goroutine和channel收集结果

4. **04-mutex.go**：演示互斥锁在并发环境中的使用
   - 数据竞争问题展示
   - 使用互斥锁保护共享资源
   - 使用读写锁优化读多写少的场景

5. **05-channels.go**：展示通道在goroutine间通信的用法
   - 基本通道操作
   - 单向通道的使用
   - 带缓冲的通道
   - 使用select多路复用通道
   - 工作池模式实现

6. **06-context.go**：展示使用context包控制goroutine生命周期
   - 基本context和取消操作
   - 带超时的context
   - 带截止时间的context
   - 带值的context
   - context与HTTP请求

## 如何运行

可以使用以下命令运行每个示例：

```bash
go run 01-basic-goroutine.go
go run 02-waitgroup.go
go run 03-concurrent-http.go
go run 04-mutex.go
go run 05-channels.go
go run 06-context.go
```

## 学习要点

1. **goroutine**：轻量级线程，是Go并发的基础
2. **WaitGroup**：用于等待一组goroutine完成
3. **Mutex**：互斥锁，保护共享资源
4. **Channel**：goroutine之间通信的主要机制
5. **Context**：用于跨API边界和goroutine传递截止时间、取消信号和请求相关值

## 并发编程最佳实践

1. **正确同步**：使用适当的同步机制（如WaitGroup、Mutex、Channel）来协调goroutine
2. **避免数据竞争**：共享数据访问必须同步，使用`go run -race`检测竞争条件
3. **资源管理**：确保在goroutine完成时释放所有资源
4. **错误处理**：在goroutine中妥善处理错误，并通过channel或其他方式传递给主goroutine
5. **超时控制**：使用context或select与timer结合，确保操作不会无限阻塞
6. **优雅退出**：实现干净的关闭机制，确保所有goroutine都能正确结束

这些示例旨在帮助理解Go的并发模型和工具，掌握后可以构建高效且可靠的并发程序。 