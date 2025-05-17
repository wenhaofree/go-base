// 01-read-file.go 展示了Go中读取文件的多种方法
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	fmt.Println("=== Go文件读取示例 ===")

	// 准备示例文件
	sampleFile := prepareExampleFile()
	defer os.Remove(sampleFile) // 程序结束时删除临时文件

	fmt.Printf("\n1. 使用 os.ReadFile 读取整个文件（推荐用于小文件）\n")
	readEntireFile(sampleFile)

	fmt.Printf("\n2. 使用 bufio.Scanner 按行读取（高效处理大文件）\n")
	readFileByLine(sampleFile)

	fmt.Printf("\n3. 使用 io.ReadAll 通过文件句柄读取\n")
	readFileWithFileHandle(sampleFile)

	fmt.Printf("\n4. 分块读取文件（处理超大文件）\n")
	readFileInChunks(sampleFile)
}

// prepareExampleFile 创建一个临时示例文件用于演示
func prepareExampleFile() string {
	// 获取临时目录
	tempDir := os.TempDir()
	// 创建一个临时文件路径
	tempFile := filepath.Join(tempDir, "go-file-example.txt")

	// 写入一些示例内容
	content := []byte(`这是第一行内容。
这是第二行内容。
这是第三行内容。
这是第四行，包含一些数字：123456789。
这是最后一行内容。`)

	// 写入文件，如果出错则打印错误并结束程序
	err := os.WriteFile(tempFile, content, 0644)
	if err != nil {
		fmt.Printf("创建示例文件失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("已创建示例文件: %s\n", tempFile)
	return tempFile
}

// readEntireFile 使用os.ReadFile一次性读取整个文件内容
// 优点: 简单易用
// 缺点: 对于大文件会占用大量内存
func readEntireFile(filePath string) {
	// 读取整个文件内容到内存
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("读取文件失败: %v\n", err)
		return
	}

	// 打印读取的内容
	fmt.Printf("文件内容 (%d 字节):\n", len(data))
	fmt.Println("---------------------")
	fmt.Println(string(data))
	fmt.Println("---------------------")
}

// readFileByLine 使用bufio.Scanner按行读取文件
// 优点: 内存效率高，适合处理大文件
// 缺点: 稍微复杂一些
func readFileByLine(filePath string) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("打开文件失败: %v\n", err)
		return
	}
	// 确保文件会被关闭
	defer file.Close()

	// 创建一个扫描器
	scanner := bufio.NewScanner(file)

	// 逐行读取
	lineCount := 0
	for scanner.Scan() {
		lineCount++
		line := scanner.Text()
		fmt.Printf("第 %d 行: %s\n", lineCount, line)
	}

	// 检查是否有扫描错误
	if err := scanner.Err(); err != nil {
		fmt.Printf("扫描文件时出错: %v\n", err)
	}
}

// readFileWithFileHandle 使用文件句柄和io.ReadAll读取文件
// 这种方法在需要先打开文件进行其他操作，然后再读取全部内容时很有用
func readFileWithFileHandle(filePath string) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("打开文件失败: %v\n", err)
		return
	}
	defer file.Close()

	// 读取所有内容
	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("读取文件失败: %v\n", err)
		return
	}

	fmt.Printf("文件大小: %d 字节\n", len(data))
	fmt.Printf("文件内容前50个字符: %s\n", string(data[:min(50, len(data))]))
}

// readFileInChunks 以固定大小的块读取文件
// 优点: 可以控制内存使用量，适合处理超大文件
func readFileInChunks(filePath string) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("打开文件失败: %v\n", err)
		return
	}
	defer file.Close()

	// 定义缓冲区大小
	bufferSize := 16 // 小的缓冲区用于演示，实际使用中通常是4KB或更大
	buffer := make([]byte, bufferSize)

	var totalBytes int
	for {
		// 读取一块数据
		bytesRead, err := file.Read(buffer)
		totalBytes += bytesRead

		if bytesRead > 0 {
			fmt.Printf("读取了 %d 字节: %s\n", bytesRead, string(buffer[:bytesRead]))
		}

		// 检查是否到达文件结尾或发生错误
		if err != nil {
			if err == io.EOF {
				fmt.Println("已到达文件结尾")
			} else {
				fmt.Printf("读取文件时出错: %v\n", err)
			}
			break
		}
	}

	fmt.Printf("总共读取: %d 字节\n", totalBytes)
}

// min返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
} 