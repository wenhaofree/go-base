# 域名WHOIS信息查询工具

这是一个使用Go语言实现的域名WHOIS信息查询工具。它可以查询指定关键词在主流域名后缀下的注册状态、注册时间以及详细的WHOIS信息。

## 功能特点

- 输入关键词，自动查询多个主流域名后缀（.com, .net, .org, .cn等）
- 显示域名是否已被注册
- 对于已注册的域名，显示注册时间、到期时间、注册人和注册商信息
- 提供完整的WHOIS原始信息查看选项

## 支持的域名后缀

- .com
- .net
- .org
- .cn
- .io
- .co
- .ai
- .app
- .xyz
- .run
- .me
- .pro
- .top
- .club
- .so

## 使用方法

### 交互式模式

1. 编译运行程序
```bash
cd demo-domain
go run .
```

2. 输入要查询的域名关键词（不包含后缀）

3. 程序会自动查询关键词在各个主流域名后缀下的状态

4. 对于已注册的域名，可以选择是否查看完整的WHOIS信息

### 命令行模式

直接查询指定域名：

```bash
go run . -domain example.com
```

查询并显示完整WHOIS信息：

```bash
go run . -domain example.com -full
```

也可以只提供关键词，默认会查询.com域名：

```bash
go run . -domain example
```

### 列表模式

以表格形式查询关键词在所有支持的域名后缀下的注册状态：

```bash
go run . -list example
```

或者使用已有的domain参数，同时添加showlist标志：

```bash
go run . -domain example -showlist
```

列表模式会并行查询所有域名，并以表格形式显示结果，包括域名、状态、注册时间和注册商信息。

## 示例

```
域名WHOIS信息查询工具
请输入关键词（不包含后缀）：
google

正在查询关键词 'google' 的域名信息...

检查域名: google.com
  状态: 已注册
  注册时间: 1997-09-15T04:00:00Z
  到期时间: 2028-09-14T04:00:00Z
  注册人: Google LLC
  注册商: MarkMonitor Inc.

  是否显示完整WHOIS信息? (y/n)
```

## 技术说明

该工具直接连接WHOIS服务器（端口43）进行查询，解析返回的文本信息以提取关键数据。工具使用正则表达式匹配不同WHOIS服务器返回的不同格式信息。 