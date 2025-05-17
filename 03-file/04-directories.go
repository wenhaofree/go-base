// 04-directories.go 展示了Go中目录相关操作的方法
package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

func main() {
	fmt.Println("=== Go目录操作示例 ===")

	// 创建一个临时目录作为基础目录
	tempBaseDir, err := os.MkdirTemp("", "go-dir-example")
	if err != nil {
		fmt.Printf("创建临时基础目录失败: %v\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tempBaseDir) // 程序结束时递归删除所有临时目录

	fmt.Printf("已创建临时基础目录: %s\n", tempBaseDir)

	// 1. 创建目录和子目录
	fmt.Println("\n1. 创建目录和子目录")
	createDirectories(tempBaseDir)

	// 2. 创建一些示例文件
	fmt.Println("\n2. 创建一些示例文件")
	createExampleFiles(tempBaseDir)

	// 3. 获取当前工作目录
	fmt.Println("\n3. 获取当前工作目录")
	getCurrentDirectory()

	// 4. 列出目录内容
	fmt.Println("\n4. 列出目录内容")
	listDirectoryContents(tempBaseDir)

	// 5. 递归遍历目录
	fmt.Println("\n5. 递归遍历目录")
	walkDirectory(tempBaseDir)

	// 6. 通过通配符模式查找文件
	fmt.Println("\n6. 通过通配符模式查找文件")
	findFilesByPattern(tempBaseDir, "*.txt")

	// 7. 创建和删除临时目录
	fmt.Println("\n7. 创建和删除临时目录")
	createAndRemoveTempDir()
}

// createDirectories 创建一个目录结构
func createDirectories(baseDir string) {
	// 创建一个简单的目录
	simpleDir := filepath.Join(baseDir, "simple-dir")
	err := os.Mkdir(simpleDir, 0755)
	if err != nil {
		fmt.Printf("创建简单目录失败: %v\n", err)
	} else {
		fmt.Printf("创建目录: %s\n", simpleDir)
	}

	// 创建嵌套目录 (包括不存在的父目录)
	nestedDir := filepath.Join(baseDir, "parent", "child", "grandchild")
	err = os.MkdirAll(nestedDir, 0755)
	if err != nil {
		fmt.Printf("创建嵌套目录失败: %v\n", err)
	} else {
		fmt.Printf("创建嵌套目录: %s\n", nestedDir)
	}
}

// createExampleFiles 在目录中创建一些示例文件
func createExampleFiles(baseDir string) {
	// 创建文件到根目录
	rootFile := filepath.Join(baseDir, "root-file.txt")
	err := os.WriteFile(rootFile, []byte("这是根目录文件的内容。"), 0644)
	if err != nil {
		fmt.Printf("创建根目录文件失败: %v\n", err)
	} else {
		fmt.Printf("创建文件: %s\n", rootFile)
	}

	// 创建文件到simple-dir
	simpleFile := filepath.Join(baseDir, "simple-dir", "simple-file.txt")
	err = os.WriteFile(simpleFile, []byte("这是simple-dir目录中的文件。"), 0644)
	if err != nil {
		fmt.Printf("创建simple-dir中的文件失败: %v\n", err)
	} else {
		fmt.Printf("创建文件: %s\n", simpleFile)
	}

	// 创建文件到嵌套目录
	nestedFile := filepath.Join(baseDir, "parent", "child", "nested-file.txt")
	err = os.WriteFile(nestedFile, []byte("这是嵌套目录中的文件。"), 0644)
	if err != nil {
		fmt.Printf("创建嵌套目录中的文件失败: %v\n", err)
	} else {
		fmt.Printf("创建文件: %s\n", nestedFile)
	}

	// 创建一个JSON文件
	jsonFile := filepath.Join(baseDir, "config.json")
	err = os.WriteFile(jsonFile, []byte("{\n  \"name\": \"示例配置\",\n  \"version\": 1.0\n}"), 0644)
	if err != nil {
		fmt.Printf("创建JSON文件失败: %v\n", err)
	} else {
		fmt.Printf("创建文件: %s\n", jsonFile)
	}
}

// getCurrentDirectory 获取当前工作目录
func getCurrentDirectory() {
	// 获取当前目录
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("获取当前目录失败: %v\n", err)
		return
	}
	fmt.Printf("当前工作目录: %s\n", currentDir)
}

// listDirectoryContents 列出目录中的所有条目
func listDirectoryContents(dirPath string) {
	// 读取目录
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		fmt.Printf("读取目录 %s 失败: %v\n", dirPath, err)
		return
	}

	// 显示目录内容
	fmt.Printf("目录 %s 中的内容:\n", dirPath)
	for i, entry := range entries {
		// 获取条目的更多信息
		info, err := entry.Info()
		if err != nil {
			fmt.Printf("  %d. %s [获取信息失败: %v]\n", i+1, entry.Name(), err)
			continue
		}

		// 根据条目类型显示不同信息
		if entry.IsDir() {
			fmt.Printf("  %d. [目录] %s\n", i+1, entry.Name())
		} else {
			fmt.Printf("  %d. [文件] %s (大小: %d字节, 修改时间: %s)\n",
				i+1, entry.Name(), info.Size(), info.ModTime().Format("2006-01-02 15:04:05"))
		}
	}

	// 如果这是基础目录，递归列出子目录
	if dirPath == filepath.Join(dirPath, filepath.Base(dirPath)) {
		// 列出简单目录内容
		simpleDir := filepath.Join(dirPath, "simple-dir")
		fmt.Printf("\n简单目录内容:\n")
		listDirectoryContents(simpleDir)

		// 列出嵌套目录内容
		parentDir := filepath.Join(dirPath, "parent")
		if _, err := os.Stat(parentDir); err == nil {
			fmt.Printf("\n嵌套目录内容 (%s):\n", parentDir)
			listDirectoryContents(parentDir)
		}
	}
}

// walkDirectory 递归遍历目录结构
func walkDirectory(dirPath string) {
	fmt.Printf("递归遍历目录 %s:\n", dirPath)

	// 使用filepath.Walk遍历整个目录树
	err := filepath.Walk(dirPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("访问 %s 时出错: %v\n", path, err)
			return nil // 继续遍历
		}

		// 计算相对路径，便于显示
		relPath, err := filepath.Rel(dirPath, path)
		if err != nil {
			relPath = path
		}
		if relPath == "." {
			relPath = "[根目录]"
		}

		// 获取路径深度，用于缩进显示
		depth := len(filepath.SplitList(relPath))
		indent := ""
		for i := 1; i < depth; i++ {
			indent += "  "
		}

		// 显示文件或目录的信息
		if info.IsDir() {
			fmt.Printf("%s📁 %s/\n", indent, filepath.Base(path))
		} else {
			fmt.Printf("%s📄 %s (%d字节)\n", indent, filepath.Base(path), info.Size())
		}

		return nil
	})

	if err != nil {
		fmt.Printf("遍历目录时出错: %v\n", err)
	}
}

// findFilesByPattern 使用通配符模式查找文件
func findFilesByPattern(dirPath, pattern string) {
	fmt.Printf("在 %s 中查找匹配 %s 的文件:\n", dirPath, pattern)

	// 使用filepath.Glob查找匹配的文件
	matches, err := filepath.Glob(filepath.Join(dirPath, "**", pattern))
	if err != nil {
		fmt.Printf("查找文件失败: %v\n", err)
		return
	}

	// 由于filepath.Glob不支持**递归通配符，我们需要手动遍历
	var allMatches []string

	// 首先添加直接匹配的文件
	directMatches, _ := filepath.Glob(filepath.Join(dirPath, pattern))
	allMatches = append(allMatches, directMatches...)

	// 然后通过Walk函数遍历所有子目录查找匹配文件
	_ = filepath.Walk(dirPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() && matchPattern(filepath.Base(path), pattern) {
			allMatches = append(allMatches, path)
		}
		return nil
	})

	// 显示找到的文件
	if len(allMatches) == 0 {
		fmt.Println("  未找到匹配的文件")
	} else {
		for i, match := range allMatches {
			info, err := os.Stat(match)
			if err != nil {
				fmt.Printf("  %d. %s [获取信息失败: %v]\n", i+1, match, err)
				continue
			}

			relPath, _ := filepath.Rel(dirPath, match)
			fmt.Printf("  %d. %s (大小: %d字节, 修改时间: %s)\n",
				i+1, relPath, info.Size(), info.ModTime().Format("2006-01-02 15:04:05"))
		}
	}
}

// matchPattern 检查文件名是否匹配简单的通配符模式
// 只支持*和?通配符
func matchPattern(name, pattern string) bool {
	matched, err := filepath.Match(pattern, name)
	if err != nil {
		return false
	}
	return matched
}

// createAndRemoveTempDir 创建和删除临时目录
func createAndRemoveTempDir() {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "go-temp-dir-*")
	if err != nil {
		fmt.Printf("创建临时目录失败: %v\n", err)
		return
	}
	fmt.Printf("创建临时目录: %s\n", tempDir)

	// 在临时目录中创建一个文件
	tempFile := filepath.Join(tempDir, "temp-file.txt")
	err = os.WriteFile(tempFile, []byte("这是临时文件的内容。"), 0644)
	if err != nil {
		fmt.Printf("创建临时文件失败: %v\n", err)
	} else {
		fmt.Printf("创建临时文件: %s\n", tempFile)
	}

	// 等待一秒后删除临时目录
	fmt.Println("1秒后将删除临时目录...")
	time.Sleep(1 * time.Second)

	// 递归删除目录及其内容
	err = os.RemoveAll(tempDir)
	if err != nil {
		fmt.Printf("删除临时目录失败: %v\n", err)
	} else {
		fmt.Printf("已删除临时目录: %s\n", tempDir)
	}

	// 确认目录已被删除
	_, err = os.Stat(tempDir)
	if os.IsNotExist(err) {
		fmt.Println("已确认临时目录不存在")
	} else if err != nil {
		fmt.Printf("检查临时目录时出错: %v\n", err)
	} else {
		fmt.Println("警告: 临时目录似乎仍然存在")
	}
} 