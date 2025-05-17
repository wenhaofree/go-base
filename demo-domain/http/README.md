# HTTP域名查询工具

这是一个使用HTTP请求实现的域名查询工具，它通过调用第三方API来快速检查多个域名的注册状态。

## 功能特点

- 使用Server-Sent Events (SSE)技术进行初步的异步实时查询。
- 对已初步确认注册的域名，通过调用独立的WHOIS API (`https://instant.who.sb/api/v1/whois`) 获取更详细、准确的WHOIS信息。
- 支持同时查询多个域名后缀
- 显示域名注册状态、注册商、精确的注册时间、到期时间及详细的域名状态
- 以表格形式清晰展示结果

## 支持的域名后缀

支持以下扩展后的主流及常用域名后缀 (已尝试按热门程度和类型排列):

**核心通用 TLDs:**
- .com, .net, .org, .info

**常用国家代码 TLDs (部分具有通用性):**
- .io, .co, .ai, .cn, .uk, .de, .jp, .au, .ca, .fr, .eu, .us, .me, .tv, .cc

**新通用 TLDs:**
- .app, .xyz, .club, .online, .store, .tech, .site, .space, .website, .agency, .dev
- .pro, .top, .run, .so, .live, .news, .global, .today

**行业特定或小众 TLDs:**
- .video, .domains, .link, .shop, .art, .blog, .design, .photography, .guru, .biz, .mobi

## 使用方法

```bash
cd demo-domain/http
go run test.go [关键词]
```

例如，查询关键词"example"在所有支持的域名后缀下的注册状态：

```bash
go run test.go example
```

## 技术实现

- 使用Go的`net/http`包发送HTTP GET请求
- 通过`bufio.Reader`解析SSE格式的响应流
- 使用goroutine实现异步处理
- 使用channel进行数据传输和错误处理
- 设置超时机制确保程序不会无限期等待

## API来源

该工具使用`https://instant.who.sb/api/v1/check`接口获取域名信息。API支持SSE实时查询，每个域名结果会单独作为一个事件返回。

## 示例输出

```
正在查询关键词 'example' 的域名信息...

正在发送SSE请求...
开始接收数据流...
---------------------------------------------------------------------------------------------------------------------------------------------
域名               状态      注册商                         注册时间                  到期时间                  详细状态
---------------------------------------------------------------------------------------------------------------------------------------------
example.com        已注册    Cloudflare, Inc.               2020-09-15T04:00:00Z      2025-09-15T04:00:00Z      clientTransferProhibited
example.org        已注册    Public Interest Registry       2019-05-08T10:00:00Z      2026-05-08T10:00:00Z      clientUpdateProhibited,clientDeleteProhibited,clientTransferProhibited
example.net        未注册    -                              -                         -                         -
example.ai         已注册    GoDaddy.com, LLC               2017-12-15T22:38:47Z      2025-08-02T22:38:47Z      clientUpdateProhibited,clientDeleteProhibited,clientTransferProhibited
example.xyz        已注册    2014-03-20T13:49:27Z      2025-03-20T13:49:27Z
---------------------------------------------------------------------------------------------------------------------------------------------
所有域名检查完成
``` 