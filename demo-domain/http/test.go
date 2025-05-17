package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
	"io"
)

// 定义结构体来存储解析后的SSE事件
type SSEEvent struct {
	Event string
	Data  map[string]interface{}
	ID    string
}

// 定义 /api/v1/whois 接口的响应结构
type WhoisAPIResponse struct {
	Type   int                    `json:"type"`
	Prices []interface{}          `json:"prices"`
	Parsed ParsedWhoisData        `json:"parsed"`
	Raw    string                 `json:"raw"`
}

type ParsedWhoisData struct {
	ID          string `json:"id"`
	Registrar   string `json:"registrar"`
	Registered  string `json:"registered"`
	Expires     string `json:"expires"`
	Status      string `json:"status"`
	Nameservers string `json:"nameservers"`
	Name        string `json:"name"`
	Suffix      string `json:"suffix"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("请提供要查询的域名关键词")
		fmt.Println("用法: go run test.go [关键词]")
		os.Exit(1)
	}

	keyword := os.Args[1]
	fmt.Printf("正在查询关键词 '%s' 的域名信息...\n\n", keyword)

	// 构建域名列表，尝试按热门程度和类型排列
	tlds := []string{
		// 核心通用 TLDs
		".com", ".net", ".org", ".info",
		// 常用国家代码 TLDs (部分具有通用性)
		".io", ".co", ".ai", ".cn", ".uk", ".de", ".jp", ".au", ".ca", ".fr", ".eu", ".us", ".me", ".tv", ".cc", 
		// 新通用 TLDs
		".app", ".xyz", 
		// ".club", 
		// ".online", 
		// ".tech", ".site", ".space", ".website", ".dev",
		// ".pro", ".top", ".run", ".so", ".live", ".news", ".global", ".today",
		// 行业特定或小众
		// ".video", ".domains", ".link", ".shop", ".art", ".blog", ".design", ".photography", ".guru", ".biz", ".mobi",
		// ".store",
	}

	var domains []string
	for _, tld := range tlds {
		domains = append(domains, keyword+tld)
	}

	domainsParam := strings.Join(domains, "%2C")
	urlSSE := "https://instant.who.sb/api/v1/check?domain=" + domainsParam + "&sse=true&return_dates=true&return-prices=true"

	reqSSE, err := http.NewRequest("GET", urlSSE, nil)
	if err != nil {
		fmt.Printf("创建SSE请求失败: %v\n", err)
		os.Exit(1)
	}
	reqSSE.Header.Set("Accept", "text/event-stream")

	clientSSE := &http.Client{Timeout: 0}

	fmt.Println("正在发送SSE请求...")
	respSSE, err := clientSSE.Do(reqSSE)
	if err != nil {
		fmt.Printf("SSE请求失败: %v\n", err)
		os.Exit(1)
	}
	defer respSSE.Body.Close()

	if respSSE.StatusCode != http.StatusOK {
		fmt.Printf("SSE服务器返回错误状态码: %d\n", respSSE.StatusCode)
		os.Exit(1)
	}

	fmt.Println("开始接收数据流...")
	fmt.Println("---------------------------------------------------------------------------------------------------------------------------------------------")
	fmt.Println("域名               状态      注册商                         注册时间                  到期时间                  详细状态")
	fmt.Println("---------------------------------------------------------------------------------------------------------------------------------------------")

	readerSSE := bufio.NewReader(respSSE.Body)
	var currentEvent SSEEvent
	resultCount := 0
	timeout := time.After(60 * time.Second) // 增加超时时间以容纳后续的WHOIS查询
	dataChan := make(chan SSEEvent)
	errChan := make(chan error)

	go func() {
		for {
			line, err := readerSSE.ReadString('\n')
			if err != nil {
				// 如果是EOF错误，表示流正常结束，关闭数据通道并退出goroutine
				if err == io.EOF {
					// 可以在这里选择关闭dataChan，或者让主循环通过超时或resultCount来结束
					// fmt.Println("SSE流结束 (EOF)") // 可选的调试信息
					return 
				}
				// 对于其他读取错误，发送到错误通道
				errChan <- fmt.Errorf("读取SSE数据失败: %v", err)
				return
			}
			line = strings.TrimSpace(line)
			if line == "" {
				if currentEvent.Event != "" && currentEvent.Data != nil {
					dataChan <- currentEvent
					currentEvent = SSEEvent{}
				}
				continue
			}
			if strings.HasPrefix(line, "event:") {
				currentEvent.Event = strings.TrimSpace(strings.TrimPrefix(line, "event:"))
			} else if strings.HasPrefix(line, "data:") {
				jsonData := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
				var data map[string]interface{}
				if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
					errChan <- fmt.Errorf("解析SSE JSON失败: %v", err)
					return
				}
				currentEvent.Data = data
			} else if strings.HasPrefix(line, "id:") {
				currentEvent.ID = strings.TrimSpace(strings.TrimPrefix(line, "id:"))
			}
		}
	}()

	httpClient := &http.Client{Timeout: 10 * time.Second}

	for {
		select {
		case event := <-dataChan:
			if event.Event == "shallow-checked" || event.Event == "whois-cache-checked" {
				domain, _ := event.Data["domain"].(string)
				meta, ok := event.Data["meta"].(map[string]interface{})
				if !ok {
					fmt.Printf("%-20s %-10s (元数据错误)\n", domain, "错误")
					resultCount++
					continue
				}

				existed, _ := meta["existed"].(string)
				status := "未注册"
				regDate := "-"
				expDate := "-"
				registrar := "-"
				domainDetailedStatus := "-"

				if existed == "yes" {
					status = "已注册"
					// 从SSE事件中获取初步日期
					if event.Event == "shallow-checked" {
						if dates, ok := meta["dates"].(map[string]interface{}); ok {
							if created, ok := dates["created"].(string); ok && created != "" {
								regDate = created
							}
						}
					} else if event.Event == "whois-cache-checked" {
						if rDate, ok := meta["registered"].(string); ok && rDate != "" {
							regDate = rDate
						}
						if eDate, ok := meta["expires"].(string); ok && eDate != "" {
							expDate = eDate
						}
					}

					// 调用 /api/v1/whois 获取更详细信息
					urlWhois := fmt.Sprintf("https://instant.who.sb/api/v1/whois?domain=%s&cache=true&return-prices=false", domain)
					reqWhois, _ := http.NewRequest("GET", urlWhois, nil)
					respWhois, err := httpClient.Do(reqWhois)
					if err == nil && respWhois.StatusCode == http.StatusOK {
						bodyBytes, _ := ioutil.ReadAll(respWhois.Body)
						var whoisData WhoisAPIResponse
						if json.Unmarshal(bodyBytes, &whoisData) == nil && whoisData.Parsed.ID != "" {
							registrar = whoisData.Parsed.Registrar
							if whoisData.Parsed.Registered != "" {
								regDate = whoisData.Parsed.Registered
							}
							if whoisData.Parsed.Expires != "" {
								expDate = whoisData.Parsed.Expires
							}
							// domainDetailedStatus = whoisData.Parsed.Status
						}
						respWhois.Body.Close()
					}
				}
				fmt.Printf("%-20s %-10s %-30s %-25s %-25s %s\n", domain, status, registrar, regDate, expDate, domainDetailedStatus)
				resultCount++
				if resultCount >= len(domains) {
					fmt.Println("---------------------------------------------------------------------------------------------------------------------------------------------")
					fmt.Println("所有域名检查完成")
					return
				}
			}
		case err := <-errChan:
			fmt.Printf("发生错误: %v\n", err)
			return
		case <-timeout:
			fmt.Println("\n查询超时")
			if resultCount < len(domains) {
				fmt.Printf("已完成 %d/%d 个域名的查询。\n", resultCount, len(domains))
			}
			return
		}
	}
}
