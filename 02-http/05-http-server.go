// 05-http-server.go 展示了Go中HTTP服务器的基本用法
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// 定义一个简单的内存数据存储
type DataStore struct {
	sync.RWMutex
	items map[string]interface{}
}

// 创建新的数据存储
func NewDataStore() *DataStore {
	return &DataStore{
		items: make(map[string]interface{}),
	}
}

// 设置值
func (ds *DataStore) Set(key string, value interface{}) {
	ds.Lock()
	defer ds.Unlock()
	ds.items[key] = value
}

// 获取值
func (ds *DataStore) Get(key string) (interface{}, bool) {
	ds.RLock()
	defer ds.RUnlock()
	value, exists := ds.items[key]
	return value, exists
}

// 定义响应结构体
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Time    time.Time   `json:"time"`
}

func main() {
	fmt.Println("启动HTTP服务器示例...")
	
	// 创建一个数据存储
	store := NewDataStore()
	// 预填充一些数据
	store.Set("greeting", "你好，世界!")
	store.Set("count", 42)
	
	// 创建一个简单的中间件记录请求
	loggingMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			
			// 记录请求信息
			log.Printf("收到请求: %s %s %s", r.Method, r.URL.Path, r.RemoteAddr)
			
			// 调用下一个处理器
			next.ServeHTTP(w, r)
			
			// 记录请求处理时间
			log.Printf("请求处理完成: %s %s %s (耗时: %v)", 
				r.Method, r.URL.Path, r.RemoteAddr, time.Since(start))
		})
	}
	
	// 创建路由器 (可以使用 "github.com/gorilla/mux" 获得更强大的路由功能)
	mux := http.NewServeMux()
	
	// 静态文件服务
	// 如果有index.html等文件，可以通过 http://localhost:8080/ 访问
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	
	// 首页路由
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Go HTTP 服务器示例</title>
			<style>
				body { font-family: Arial, sans-serif; margin: 40px; line-height: 1.6; }
				h1 { color: #333; }
				.endpoint { background: #f4f4f4; padding: 10px; border-radius: 5px; margin-bottom: 10px; }
				code { background: #e0e0e0; padding: 2px 4px; border-radius: 3px; }
			</style>
		</head>
		<body>
			<h1>Go HTTP 服务器示例</h1>
			<p>这是一个简单的HTTP服务器示例，提供以下API端点:</p>
			
			<div class="endpoint">
				<strong>GET /api/data/:key</strong>
				<p>获取指定键的值。例如: <code>/api/data/greeting</code></p>
			</div>
			
			<div class="endpoint">
				<strong>POST /api/data/:key</strong>
				<p>设置指定键的值。需要发送JSON请求体，如: <code>{"value": "新的值"}</code></p>
			</div>
			
			<div class="endpoint">
				<strong>GET /api/time</strong>
				<p>获取当前服务器时间</p>
			</div>
		</body>
		</html>
		`)
	})
	
	// API路由 - 获取时间
	mux.HandleFunc("/api/time", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "只支持GET方法", http.StatusMethodNotAllowed)
			return
		}
		
		response := Response{
			Status: "success",
			Data:   time.Now(),
			Time:   time.Now(),
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
	
	// API路由 - 数据存储
	mux.HandleFunc("/api/data/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		
		// 解析路径中的键
		key := r.URL.Path[len("/api/data/"):]
		if key == "" {
			http.Error(w, "需要指定键", http.StatusBadRequest)
			return
		}
		
		switch r.Method {
		case http.MethodGet:
			// 获取数据
			value, exists := store.Get(key)
			if !exists {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(Response{
					Status:  "error",
					Message: "未找到指定的键",
					Time:    time.Now(),
				})
				return
			}
			
			json.NewEncoder(w).Encode(Response{
				Status: "success",
				Data:   value,
				Time:   time.Now(),
			})
			
		case http.MethodPost:
			// 设置数据
			var requestData struct {
				Value interface{} `json:"value"`
			}
			
			err := json.NewDecoder(r.Body).Decode(&requestData)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(Response{
					Status:  "error",
					Message: "无效的请求体",
					Time:    time.Now(),
				})
				return
			}
			
			store.Set(key, requestData.Value)
			
			json.NewEncoder(w).Encode(Response{
				Status:  "success",
				Message: "数据已更新",
				Time:    time.Now(),
			})
			
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(Response{
				Status:  "error",
				Message: "只支持GET和POST方法",
				Time:    time.Now(),
			})
		}
	})
	
	// 将中间件应用到路由器
	handler := loggingMiddleware(mux)
	
	// 启动服务器
	port := 8080
	serverAddr := ":" + strconv.Itoa(port)
	fmt.Printf("服务器正在监听 http://localhost%s\n", serverAddr)
	fmt.Println("按Ctrl+C停止服务器")
	
	// 设置服务器配置
	server := &http.Server{
		Addr:           serverAddr,
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}
	
	// 启动服务器
	log.Fatal(server.ListenAndServe())
} 