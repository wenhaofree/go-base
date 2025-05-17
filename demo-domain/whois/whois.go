package whois

import (
	"fmt"
	"io"
	"net"
	"regexp"
	"strings"
	"time"
)

// WhoisResult 包含WHOIS查询的结果
type WhoisResult struct {
	Domain         string
	IsRegistered   bool
	CreationDate   string
	ExpirationDate string
	Registrant     string
	Registrar      string
	RawText        string
}

// WHOIS服务器映射表
var whoisServers = map[string]string{
	".com":  "whois.verisign-grs.com",
	".net":  "whois.verisign-grs.com",
	".org":  "whois.pir.org",
	".info": "whois.afilias.net",
	".cn":   "whois.cnnic.cn",
	".io":   "whois.nic.io",
	".co":   "whois.nic.co",
	".ai":   "whois.nic.ai",
	".app":  "whois.identitydark.cloud",
	".xyz":  "whois.nic.xyz",
	".run":  "whois.donuts.co",
	".me":   "whois.nic.me",
	".pro":  "whois.afilias.net",
	".top":  "whois.nic.top",
	".club": "whois.nic.club",
	".so":   "whois.nic.so",
}

// Query 查询域名的WHOIS信息
func Query(domain string) (*WhoisResult, error) {
	// 确定WHOIS服务器
	var server string
	for tld, srv := range whoisServers {
		if strings.HasSuffix(domain, tld) {
			server = srv
			break
		}
	}

	if server == "" {
		return nil, fmt.Errorf("不支持的域名后缀")
	}

	// 连接WHOIS服务器
	conn, err := net.DialTimeout("tcp", server+":43", 5*time.Second)
	if err != nil {
		return nil, fmt.Errorf("连接WHOIS服务器失败: %w", err)
	}
	defer conn.Close()

	// 发送查询请求
	_, err = conn.Write([]byte(domain + "\r\n"))
	if err != nil {
		return nil, fmt.Errorf("发送查询请求失败: %w", err)
	}

	// 读取响应
	buffer := make([]byte, 0, 4096)
	tmp := make([]byte, 1024)
	for {
		n, err := conn.Read(tmp)
		if err != nil {
			if err != io.EOF {
				return nil, fmt.Errorf("读取响应失败: %w", err)
			}
			break
		}
		buffer = append(buffer, tmp[:n]...)
	}

	rawText := string(buffer)
	result := &WhoisResult{
		Domain:  domain,
		RawText: rawText,
	}

	// 解析响应
	parseResult(result)

	return result, nil
}

// parseResult 解析WHOIS响应文本
func parseResult(result *WhoisResult) {
	// 检查是否已注册
	noMatchPatterns := []string{
		"No match for",
		"NOT FOUND",
		"No Data Found",
		"Domain not found",
		"The queried object does not exist",
	}

	result.IsRegistered = true
	for _, pattern := range noMatchPatterns {
		if strings.Contains(result.RawText, pattern) {
			result.IsRegistered = false
			return
		}
	}

	// 解析创建日期
	creationDatePatterns := []string{
		`(?i)Creation Date: (.+)`,
		`(?i)Created on: (.+)`,
		`(?i)Registration Time: (.+)`,
	}
	for _, pattern := range creationDatePatterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(result.RawText)
		if len(matches) > 1 {
			result.CreationDate = strings.TrimSpace(matches[1])
			break
		}
	}

	// 解析过期日期
	expirationDatePatterns := []string{
		`(?i)Expiration Date: (.+)`,
		`(?i)Registry Expiry Date: (.+)`,
		`(?i)Expiry date: (.+)`,
	}
	for _, pattern := range expirationDatePatterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(result.RawText)
		if len(matches) > 1 {
			result.ExpirationDate = strings.TrimSpace(matches[1])
			break
		}
	}

	// 解析注册人
	registrantPatterns := []string{
		`(?i)Registrant Name: (.+)`,
		`(?i)Registrant: (.+)`,
		`(?i)Registrant Organization: (.+)`,
	}
	for _, pattern := range registrantPatterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(result.RawText)
		if len(matches) > 1 {
			result.Registrant = strings.TrimSpace(matches[1])
			break
		}
	}

	// 解析注册商
	registrarPatterns := []string{
		`(?i)Registrar: (.+)`,
		`(?i)Sponsoring Registrar: (.+)`,
		`(?i)Registrar Name: (.+)`,
	}
	for _, pattern := range registrarPatterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(result.RawText)
		if len(matches) > 1 {
			result.Registrar = strings.TrimSpace(matches[1])
			break
		}
	}
} 