// 03-file-info.go 展示了Go中获取和修改文件信息的方法
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func main() {
	fmt.Println("=== Go文件信息和元数据操作示例 ===")

	// 创建一个临时文件用于演示
	tempDir := os.TempDir()
	filePath := filepath.Join(tempDir, "file-info-example.txt")

	// 创建示例文件
	err := os.WriteFile(filePath, []byte("这是示例文件内容，用于演示文件信息和元数据操作。"), 0644)
	if err != nil {
		fmt.Printf("创建示例文件失败: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(filePath) // 程序结束时删除临时文件

	fmt.Printf("已创建示例文件: %s\n", filePath)

	// 1. 获取文件信息
	fmt.Println("\n1. 获取文件基本信息 (os.Stat)")
	getFileInfo(filePath)

	// 2. 检查文件是否存在
	fmt.Println("\n2. 检查文件是否存在")
	checkFileExists(filePath)
	checkFileExists(filePath + ".nonexistent")

	// 3. 修改文件权限
	fmt.Println("\n3. 修改文件权限")
	changeFilePermissions(filePath)

	// 4. 修改文件时间戳
	fmt.Println("\n4. 修改文件时间戳")
	changeFileTimes(filePath)

	// 5. 重命名文件
	fmt.Println("\n5. 重命名文件")
	newPath := filepath.Join(tempDir, "renamed-file.txt")
	renameFile(filePath, newPath)
	// 更新文件路径变量
	filePath = newPath

	// 6. 获取绝对路径
	fmt.Println("\n6. 获取文件的绝对路径")
	getAbsolutePath(filePath)

	// 7. 分析文件路径
	fmt.Println("\n7. 分析文件路径组成部分")
	analyzeFilePath(filePath)
}

// getFileInfo 获取文件的基本信息
func getFileInfo(filePath string) {
	// 获取文件信息
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		fmt.Printf("获取文件信息失败: %v\n", err)
		return
	}

	// 显示文件信息
	fmt.Println("文件信息:")
	fmt.Printf("  名称: %s\n", fileInfo.Name())
	fmt.Printf("  大小: %d 字节\n", fileInfo.Size())
	fmt.Printf("  权限: %v\n", fileInfo.Mode())
	fmt.Printf("  修改时间: %v\n", fileInfo.ModTime())
	fmt.Printf("  是否是目录: %t\n", fileInfo.IsDir())
}

// checkFileExists 检查文件是否存在
func checkFileExists(filePath string) {
	// 尝试获取文件信息
	_, err := os.Stat(filePath)

	// 解析错误
	if err == nil {
		fmt.Printf("文件存在: %s\n", filePath)
		return
	}

	// 使用os.IsNotExist检查错误是否表示文件不存在
	if os.IsNotExist(err) {
		fmt.Printf("文件不存在: %s\n", filePath)
		return
	}

	// 其他错误
	fmt.Printf("检查文件是否存在时出错: %v\n", err)
}

// changeFilePermissions 修改文件权限
func changeFilePermissions(filePath string) {
	// 显示原始权限
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		fmt.Printf("获取文件信息失败: %v\n", err)
		return
	}
	fmt.Printf("原始文件权限: %v\n", fileInfo.Mode())

	// 修改文件权限为0600 (只有所有者可读写)
	err = os.Chmod(filePath, 0600)
	if err != nil {
		fmt.Printf("修改文件权限失败: %v\n", err)
		return
	}

	// 显示新权限
	fileInfo, _ = os.Stat(filePath)
	fmt.Printf("修改后文件权限: %v\n", fileInfo.Mode())
}

// changeFileTimes 修改文件的访问和修改时间
func changeFileTimes(filePath string) {
	// 显示原始修改时间
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		fmt.Printf("获取文件信息失败: %v\n", err)
		return
	}
	fmt.Printf("原始文件修改时间: %v\n", fileInfo.ModTime())

	// 设置一个过去的时间
	pastTime := time.Date(2020, 1, 1, 12, 0, 0, 0, time.Local)

	// 修改文件时间
	err = os.Chtimes(filePath, pastTime, pastTime)
	if err != nil {
		fmt.Printf("修改文件时间失败: %v\n", err)
		return
	}

	// 显示新的修改时间
	fileInfo, _ = os.Stat(filePath)
	fmt.Printf("修改后文件修改时间: %v\n", fileInfo.ModTime())
}

// renameFile 重命名/移动文件
func renameFile(oldPath, newPath string) {
	// 尝试删除目标文件(如果存在)，以避免冲突
	_ = os.Remove(newPath)

	// 重命名文件
	err := os.Rename(oldPath, newPath)
	if err != nil {
		fmt.Printf("重命名文件失败: %v\n", err)
		return
	}

	fmt.Printf("文件已重命名: %s -> %s\n", oldPath, newPath)

	// 确认新文件存在
	if _, err := os.Stat(newPath); err == nil {
		fmt.Printf("确认新文件存在: %s\n", newPath)
	}
}

// getAbsolutePath 获取文件的绝对路径
func getAbsolutePath(filePath string) {
	// 获取绝对路径
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		fmt.Printf("获取绝对路径失败: %v\n", err)
		return
	}

	fmt.Printf("原始路径: %s\n", filePath)
	fmt.Printf("绝对路径: %s\n", absPath)
}

// analyzeFilePath 分析文件路径的组成部分
func analyzeFilePath(filePath string) {
	// 获取路径各组成部分
	dir := filepath.Dir(filePath)
	base := filepath.Base(filePath)
	ext := filepath.Ext(filePath)
	nameWithoutExt := base[:len(base)-len(ext)]

	fmt.Println("文件路径分析:")
	fmt.Printf("  完整路径: %s\n", filePath)
	fmt.Printf("  目录部分: %s\n", dir)
	fmt.Printf("  文件名(带扩展名): %s\n", base)
	fmt.Printf("  扩展名: %s\n", ext)
	fmt.Printf("  文件名(不带扩展名): %s\n", nameWithoutExt)

	// 获取分隔符
	fmt.Printf("  路径分隔符: %s\n", string(filepath.Separator))
} 