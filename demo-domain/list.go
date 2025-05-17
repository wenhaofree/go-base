package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"

	"go-base/demo-domain/whois"
)

// RunList 运行列表模式，直接返回域名是否注册的列表
func RunList() bool {
	var keyword string
	var listMode bool

	// 解析命令行参数
	flag.StringVar(&keyword, "list", "", "以列表形式查询关键词在所有支持的域名后缀下的注册状态")
	flag.BoolVar(&listMode, "showlist", false, "启用列表模式")
	flag.Parse()

	// 如果没有提供list参数且未启用listMode，返回false表示不是列表模式
	if keyword == "" && !listMode {
		return false
	}

	// 如果启用了listMode但没有提供keyword，从其他参数中获取
	if keyword == "" && listMode {
		// 尝试从domain参数获取
		keyword = flag.Lookup("domain").Value.String()
		if keyword == "" {
			fmt.Println("错误: 列表模式需要提供关键词")
			os.Exit(1)
		}

		// 如果keyword包含后缀，则去掉后缀
		if strings.Contains(keyword, ".") {
			parts := strings.Split(keyword, ".")
			if len(parts) >= 2 {
				keyword = parts[0]
			}
		}
	}

	// 主流域名后缀
	tlds := []string{".com", ".net", ".org", ".cn", ".io", ".co", ".ai", ".app", 
		".xyz", ".run", ".me", ".pro", ".top", ".club", ".so"}

	fmt.Printf("关键词 '%s' 的域名注册状态列表:\n\n", keyword)
	fmt.Println("域名                 状态          注册时间                  注册商")
	fmt.Println("------------------- ------------- ------------------------- -----------------")

	// 使用WaitGroup进行并行查询
	var wg sync.WaitGroup
	// 使用互斥锁保护打印操作
	var mu sync.Mutex
	
	// 为每个TLD创建一个goroutine进行查询
	for _, tld := range tlds {
		wg.Add(1)
		go func(tld string) {
			defer wg.Done()
			domain := keyword + tld
			
			result, err := whois.Query(domain)
			
			mu.Lock()
			defer mu.Unlock()
			
			if err != nil {
				fmt.Printf("%-20s %-13s %s\n", domain, "查询失败", err.Error())
				return
			}

			if result.IsRegistered {
				registrar := result.Registrar
				if len(registrar) > 25 {
					registrar = registrar[:22] + "..."
				}
				fmt.Printf("%-20s %-13s %-25s %-20s\n", domain, "已注册", result.CreationDate, registrar)
			} else {
				fmt.Printf("%-20s %-13s %-25s %-20s\n", domain, "未注册", "-", "-")
			}
		}(tld)
	}

	// 等待所有查询完成
	wg.Wait()
	fmt.Println("\n查询完成。")
	return true
} 