package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"go-base/demo-domain/whois"
)

func main() {
	// 先检查是否以列表模式运行
	if RunList() {
		return
	}
	
	// 再检查是否以命令行模式运行
	if RunCmd() {
		return
	}

	fmt.Println("域名WHOIS信息查询工具")
	fmt.Println("请输入关键词（不包含后缀）：")

	reader := bufio.NewReader(os.Stdin)
	keyword, _ := reader.ReadString('\n')
	keyword = strings.TrimSpace(keyword)

	if keyword == "" {
		fmt.Println("关键词不能为空")
		return
	}

	// 主流域名后缀
	tlds := []string{".com", ".net", ".org", ".cn", ".io", ".co", ".ai", ".app", 
		".xyz", ".run", ".me", ".pro", ".top", ".club", ".so"}

	fmt.Printf("正在查询关键词 '%s' 的域名信息...\n\n", keyword)

	for _, tld := range tlds {
		domain := keyword + tld
		fmt.Printf("检查域名: %s\n", domain)

		result, err := whois.Query(domain)
		if err != nil {
			fmt.Printf("  查询失败: %s\n\n", err)
			continue
		}

		if result.IsRegistered {
			fmt.Printf("  状态: 已注册\n")
			fmt.Printf("  注册时间: %s\n", result.CreationDate)
			fmt.Printf("  到期时间: %s\n", result.ExpirationDate)
			fmt.Printf("  注册人: %s\n", result.Registrant)
			fmt.Printf("  注册商: %s\n\n", result.Registrar)
			
			fmt.Println("  是否显示完整WHOIS信息? (y/n)")
			showDetails, _ := reader.ReadString('\n')
			showDetails = strings.TrimSpace(strings.ToLower(showDetails))
			
			if showDetails == "y" || showDetails == "yes" {
				fmt.Println("\n完整WHOIS信息:")
				fmt.Println("----------------------------------------")
				fmt.Println(result.RawText)
				fmt.Println("----------------------------------------\n")
			}
		} else {
			fmt.Printf("  状态: 未注册 (可注册)\n\n")
		}
	}
} 