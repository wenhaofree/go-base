// 02-write-file.go 展示了Go中写入文件的多种方法
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

func main() {
	fmt.Println("=== Go文件写入示例 ===")

	// 获取临时目录用于存放示例文件
	tempDir := os.TempDir()

	// 示例1: 使用os.WriteFile一次性写入文件
	fmt.Println("\n1. 使用 os.WriteFile 一次性写入文件")
	writeFileAt := filepath.Join(tempDir, "example1.txt")
	writeFileAtOnce(writeFileAt)
	// 读取并显示文件内容
	printFileContent("示例1结果:", writeFileAt)
	defer os.Remove(writeFileAt) // 程序结束时删除临时文件

	// 示例2: 使用os.OpenFile和Write方法写入文件
	fmt.Println("\n2. 使用 os.OpenFile 和 Write 方法直接写入文件")
	writeFileManually := filepath.Join(tempDir, "example2.txt")
	writeFileWithOpenFile(writeFileManually)
	// 读取并显示文件内容
	printFileContent("示例2结果:", writeFileManually)
	defer os.Remove(writeFileManually)

	// 示例3: 使用bufio.Writer进行缓冲写入
	fmt.Println("\n3. 使用 bufio.Writer 进行缓冲写入")
	writeFileBuffered := filepath.Join(tempDir, "example3.txt")
	writeFileWithBuffer(writeFileBuffered)
	// 读取并显示文件内容
	printFileContent("示例3结果:", writeFileBuffered)
	defer os.Remove(writeFileBuffered)

	// 示例4: 追加内容到现有文件
	fmt.Println("\n4. 追加内容到现有文件")
	appendToFile(writeFileBuffered)
	// 读取并显示文件内容
	printFileContent("示例4结果 (追加后):", writeFileBuffered)

	// 示例5: 使用io.Copy从一个文件复制到另一个文件
	fmt.Println("\n5. 使用 io.Copy 复制文件内容")
	copyDestFile := filepath.Join(tempDir, "example5-copy.txt")
	copyFile(writeFileBuffered, copyDestFile)
	// 读取并显示文件内容
	printFileContent("示例5结果 (复制的文件):", copyDestFile)
	defer os.Remove(copyDestFile)
}

// writeFileAtOnce 使用os.WriteFile一次性写入文件
// 适合简单情况和小文件
func writeFileAtOnce(filePath string) {
	// 准备写入的内容
	content := []byte("这是使用os.WriteFile写入的内容。\n这是第二行。\n")

	// 使用0644权限写入文件
	// (0644表示: 所有者可读写，其他人只读)
	err := os.WriteFile(filePath, content, 0644)
	if err != nil {
		fmt.Printf("写入文件 %s 失败: %v\n", filePath, err)
		return
	}

	fmt.Printf("成功写入 %d 字节到文件: %s\n", len(content), filePath)
}

// writeFileWithOpenFile 使用os.OpenFile和Write方法写入文件
// 提供更多控制选项，适合需要特定文件标志和权限的情况
func writeFileWithOpenFile(filePath string) {
	// 打开文件，如果不存在则创建，如果存在则截断
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Printf("打开文件 %s 失败: %v\n", filePath, err)
		return
	}
	// 确保文件会被关闭
	defer file.Close()

	// 准备要写入的文本
	text := "这是第一行，使用OpenFile写入。\n"
	text += "这是第二行。\n"
	text += fmt.Sprintf("当前时间: %s\n", time.Now().Format("2006-01-02 15:04:05"))

	// 将字符串转换为字节数组并写入
	bytesWritten, err := file.Write([]byte(text))
	if err != nil {
		fmt.Printf("写入文件失败: %v\n", err)
		return
	}

	fmt.Printf("成功写入 %d 字节到文件: %s\n", bytesWritten, filePath)
}

// writeFileWithBuffer 使用bufio.Writer进行缓冲写入
// 适合需要频繁小写入的情况，可以提高性能
func writeFileWithBuffer(filePath string) {
	// 打开文件用于写入
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("创建文件 %s 失败: %v\n", filePath, err)
		return
	}
	defer file.Close()

	// 创建一个带缓冲的writer
	writer := bufio.NewWriter(file)

	// 写入多行内容
	lines := []string{
		"这是使用bufio.Writer写入的第一行。",
		"这是第二行，使用缓冲写入可以提高性能。",
		"适合需要频繁小写入的情况。",
		fmt.Sprintf("当前时间: %s", time.Now().Format("2006-01-02 15:04:05")),
	}

	var totalBytes int
	for _, line := range lines {
		// 写入一行文本和换行符
		bytesWritten, err := writer.WriteString(line + "\n")
		if err != nil {
			fmt.Printf("写入行失败: %v\n", err)
			return
		}
		totalBytes += bytesWritten
	}

	// 将缓冲区内容刷新到文件
	// 这一步非常重要，否则数据可能只存在于缓冲区而没有写入文件
	err = writer.Flush()
	if err != nil {
		fmt.Printf("刷新缓冲区失败: %v\n", err)
		return
	}

	fmt.Printf("成功写入 %d 字节到文件: %s\n", totalBytes, filePath)
}

// appendToFile 追加内容到现有文件
func appendToFile(filePath string) {
	// 以追加模式打开文件
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("打开文件 %s 失败: %v\n", filePath, err)
		return
	}
	defer file.Close()

	// 要追加的内容
	appendText := "\n这是追加的内容。\n"
	appendText += fmt.Sprintf("追加时间: %s\n", time.Now().Format("2006-01-02 15:04:05"))

	// 写入追加内容
	bytesWritten, err := file.WriteString(appendText)
	if err != nil {
		fmt.Printf("追加内容失败: %v\n", err)
		return
	}

	fmt.Printf("成功追加 %d 字节到文件: %s\n", bytesWritten, filePath)
}

// copyFile 将一个文件的内容复制到另一个文件
func copyFile(srcPath, dstPath string) {
	// 打开源文件
	srcFile, err := os.Open(srcPath)
	if err != nil {
		fmt.Printf("打开源文件 %s 失败: %v\n", srcPath, err)
		return
	}
	defer srcFile.Close()

	// 创建目标文件
	dstFile, err := os.Create(dstPath)
	if err != nil {
		fmt.Printf("创建目标文件 %s 失败: %v\n", dstPath, err)
		return
	}
	defer dstFile.Close()

	// 使用io.Copy复制内容
	bytesWritten, err := io.Copy(dstFile, srcFile)
	if err != nil {
		fmt.Printf("复制文件内容失败: %v\n", err)
		return
	}

	fmt.Printf("成功从 %s 复制 %d 字节到 %s\n", srcPath, bytesWritten, dstPath)
}

// printFileContent 读取并显示文件内容
func printFileContent(message, filePath string) {
	fmt.Println(message)
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("读取文件失败: %v\n", err)
		return
	}

	fmt.Println("---------------------")
	fmt.Println(string(content))
	fmt.Println("---------------------")
} 