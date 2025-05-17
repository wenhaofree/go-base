package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"go-base/demo-domain/whois"
)

// RunCmd 运行命令行模式的查询
func RunCmd() bool {
	var domain string
	var showFull bool

	// 解析命令行参数
	flag.StringVar(&domain, "domain", "", "要查询的域名")
	flag.BoolVar(&showFull, "full", false, "是否显示完整WHOIS信息")
	flag.Parse()

	// 如果没有提供domain参数，返回false表示不是命令行模式
	if domain == "" {
		return false
	}

	// 如果没有后缀，添加.com作为默认后缀
	if !strings.Contains(domain, ".") {
		domain = domain + ".com"
	}

	// 执行查询
	fmt.Printf("正在查询域名: %s\n", domain)
	result, err := whois.Query(domain)
	if err != nil {
		fmt.Printf("查询失败: %s\n", err)
		os.Exit(1)
	}

	// 显示结果
	if result.IsRegistered {
		fmt.Printf("状态: 已注册\n")
		fmt.Printf("注册时间: %s\n", result.CreationDate)
		fmt.Printf("到期时间: %s\n", result.ExpirationDate)
		fmt.Printf("注册人: %s\n", result.Registrant)
		fmt.Printf("注册商: %s\n", result.Registrar)
		
		if showFull {
			fmt.Println("\n完整WHOIS信息:")
			fmt.Println("----------------------------------------")
			fmt.Println(result.RawText)
			fmt.Println("----------------------------------------")
		}
	} else {
		fmt.Printf("状态: 未注册 (可注册)\n")
	}

	return true
} 