// 05-file-special.go 展示了Go中文件权限设置和一些特殊文件操作
package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

func main() {
	fmt.Println("=== Go文件权限和特殊文件操作示例 ===")

	// 为了安全，先创建一个临时目录
	tempDir, err := os.MkdirTemp("", "go-file-special")
	if err != nil {
		fmt.Printf("创建临时目录失败: %v\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tempDir) // 程序结束时删除临时目录

	fmt.Printf("已创建临时目录: %s\n", tempDir)

	// 1. 设置文件权限
	fmt.Println("\n1. 设置和检查文件权限")
	permissionsFile := filepath.Join(tempDir, "permissions.txt")
	filePermissions(permissionsFile)

	// 2. 创建并使用临时文件
	fmt.Println("\n2. 创建和使用临时文件")
	tempFileOperations(tempDir)

	// 3. 文件锁定示例(仅限非Windows)
	fmt.Println("\n3. 文件锁定示例")
	if runtime.GOOS != "windows" {
		fileLockingExample(tempDir)
	} else {
		fmt.Println("文件锁定示例在Windows上不完全支持，已跳过")
	}

	// 4. 使用符号链接 (仅限非Windows或Windows管理员)
	fmt.Println("\n4. 创建和使用符号链接")
	symlinkOperations(tempDir)

	// 5. 创建内存映射文件
	fmt.Println("\n5. 文件副本和备份")
	fileCopyAndBackup(tempDir)
}

// filePermissions 演示如何设置和检查文件权限
func filePermissions(filePath string) {
	// 创建文件时设置初始权限
	initialContent := []byte("这是一个用于演示权限设置的文件。\n")
	err := os.WriteFile(filePath, initialContent, 0644)
	if err != nil {
		fmt.Printf("创建文件失败: %v\n", err)
		return
	}
	fmt.Printf("已创建文件: %s\n", filePath)

	// 读取并显示初始权限
	info, err := os.Stat(filePath)
	if err != nil {
		fmt.Printf("获取文件信息失败: %v\n", err)
		return
	}
	mode := info.Mode()
	fmt.Printf("初始文件权限: %s (八进制: %o)\n", mode, mode.Perm())

	// 详细解释权限位
	explainPermissions(mode)

	// 修改文件权限为只读(对所有用户)
	fmt.Println("\n将文件设置为只读...")
	err = os.Chmod(filePath, 0444)
	if err != nil {
		fmt.Printf("修改文件权限失败: %v\n", err)
		return
	}

	// 显示新的权限
	info, _ = os.Stat(filePath)
	mode = info.Mode()
	fmt.Printf("修改后的文件权限: %s (八进制: %o)\n", mode, mode.Perm())
	explainPermissions(mode)

	// 尝试写入只读文件
	fmt.Println("\n尝试写入只读文件...")
	err = os.WriteFile(filePath, []byte("新内容"), 0444)
	if err != nil {
		fmt.Printf("预期的错误: %v\n", err)
	} else {
		fmt.Println("警告: 能够写入只读文件")
	}

	// 恢复为可写入权限
	fmt.Println("\n恢复文件为可写入...")
	err = os.Chmod(filePath, 0644)
	if err != nil {
		fmt.Printf("恢复文件权限失败: %v\n", err)
		return
	}

	// 验证可以写入
	err = os.WriteFile(filePath, []byte("已恢复写入权限，这是新内容。\n"), 0644)
	if err != nil {
		fmt.Printf("写入文件失败: %v\n", err)
	} else {
		fmt.Println("成功写入文件")
	}

	// 显示最终权限
	info, _ = os.Stat(filePath)
	mode = info.Mode()
	fmt.Printf("最终文件权限: %s (八进制: %o)\n", mode, mode.Perm())
}

// explainPermissions 详细解释文件权限
func explainPermissions(mode fs.FileMode) {
	fmt.Println("权限解释:")
	fmt.Printf("  类型: %s\n", getFileType(mode))
	fmt.Printf("  所有者权限: %s\n", getPermissionString(mode, 4, 2, 1))
	fmt.Printf("  用户组权限: %s\n", getPermissionString(mode, 0o40, 0o20, 0o10))
	fmt.Printf("  其他用户权限: %s\n", getPermissionString(mode, 4, 2, 1))
}

// getFileType 返回文件类型的字符串描述
func getFileType(mode fs.FileMode) string {
	if mode.IsDir() {
		return "目录"
	}
	if mode&fs.ModeSymlink != 0 {
		return "符号链接"
	}
	if mode&fs.ModeNamedPipe != 0 {
		return "命名管道"
	}
	if mode&fs.ModeSocket != 0 {
		return "Socket"
	}
	if mode&fs.ModeDevice != 0 {
		return "设备文件"
	}
	return "普通文件"
}

// getPermissionString 返回指定级别的权限字符串
func getPermissionString(mode fs.FileMode, r, w, x uint32) string {
	var result string
	perm := mode.Perm()

	// 读权限
	if perm&fs.FileMode(r) != 0 {
		result += "读(r) "
	} else {
		result += "   "
	}

	// 写权限
	if perm&fs.FileMode(w) != 0 {
		result += "写(w) "
	} else {
		result += "   "
	}

	// 执行权限
	if perm&fs.FileMode(x) != 0 {
		result += "执行(x)"
	}

	return result
}

// tempFileOperations 演示临时文件的创建和使用
func tempFileOperations(dirPath string) {
	// 创建一个临时文件
	// "*"会被替换为一个随机字符串
	tempFile, err := os.CreateTemp(dirPath, "example-*.txt")
	if err != nil {
		fmt.Printf("创建临时文件失败: %v\n", err)
		return
	}
	defer os.Remove(tempFile.Name()) // 确保退出时删除临时文件
	fmt.Printf("创建了临时文件: %s\n", tempFile.Name())

	// 写入一些数据
	data := []byte("这是临时文件中的数据。\n这将在程序结束时自动删除。\n")
	_, err = tempFile.Write(data)
	if err != nil {
		fmt.Printf("写入临时文件失败: %v\n", err)
		tempFile.Close()
		return
	}

	// 将文件指针移回文件开头，并读取内容
	_, err = tempFile.Seek(0, io.SeekStart)
	if err != nil {
		fmt.Printf("移动文件指针失败: %v\n", err)
		tempFile.Close()
		return
	}

	content, err := io.ReadAll(tempFile)
	if err != nil {
		fmt.Printf("读取临时文件失败: %v\n", err)
		tempFile.Close()
		return
	}

	// 关闭文件
	tempFile.Close()

	// 显示读取的内容
	fmt.Printf("临时文件内容:\n%s\n", string(content))
}

// fileLockingExample 演示一个简单的文件锁定示例
// 注意：这种方法在不同操作系统上的行为可能不同
func fileLockingExample(dirPath string) {
	lockFile := filepath.Join(dirPath, "lockfile.txt")

	// 创建或打开文件
	file, err := os.OpenFile(lockFile, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		fmt.Printf("打开锁文件失败: %v\n", err)
		return
	}
	defer file.Close()

	// 模拟一个文件锁定过程
	// 写入"锁"以表示文件正在被使用
	_, err = file.WriteString("LOCKED: " + time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		fmt.Printf("写入锁文件失败: %v\n", err)
		return
	}
	fmt.Println("文件已锁定，正在处理...")

	// 在实际应用中，这里会执行需要锁定的操作
	time.Sleep(2 * time.Second)

	// 释放锁
	err = file.Truncate(0)
	if err != nil {
		fmt.Printf("截断锁文件失败: %v\n", err)
		return
	}
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		fmt.Printf("移动文件指针失败: %v\n", err)
		return
	}
	_, err = file.WriteString("UNLOCKED: " + time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		fmt.Printf("写入锁文件失败: %v\n", err)
		return
	}
	fmt.Println("锁已释放")

	// 显示当前锁文件内容
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		fmt.Printf("移动文件指针失败: %v\n", err)
		return
	}
	content, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("读取锁文件失败: %v\n", err)
		return
	}
	fmt.Printf("锁文件最终内容: %s\n", string(content))
}

// symlinkOperations 演示符号链接操作
func symlinkOperations(dirPath string) {
	// 创建一个目标文件
	targetFile := filepath.Join(dirPath, "target.txt")
	err := os.WriteFile(targetFile, []byte("这是目标文件的内容。\n"), 0644)
	if err != nil {
		fmt.Printf("创建目标文件失败: %v\n", err)
		return
	}
	fmt.Printf("创建了目标文件: %s\n", targetFile)

	// 创建指向该文件的符号链接
	linkFile := filepath.Join(dirPath, "link.txt")
	err = os.Symlink(targetFile, linkFile)
	if err != nil {
		fmt.Printf("创建符号链接失败: %v\n", err)
		fmt.Println("注意: 在Windows上，创建符号链接通常需要管理员权限")
		return
	}
	fmt.Printf("创建了符号链接: %s -> %s\n", linkFile, targetFile)

	// 读取符号链接的目标路径
	linkTarget, err := os.Readlink(linkFile)
	if err != nil {
		fmt.Printf("读取符号链接目标失败: %v\n", err)
		return
	}
	fmt.Printf("符号链接的目标: %s\n", linkTarget)

	// 通过符号链接读取文件内容
	content, err := os.ReadFile(linkFile)
	if err != nil {
		fmt.Printf("通过符号链接读取文件失败: %v\n", err)
		return
	}
	fmt.Printf("通过符号链接读取的内容: %s", string(content))

	// 删除符号链接(但保留目标文件)
	err = os.Remove(linkFile)
	if err != nil {
		fmt.Printf("删除符号链接失败: %v\n", err)
	} else {
		fmt.Printf("已删除符号链接: %s\n", linkFile)
	}

	// 确认目标文件仍然存在
	_, err = os.Stat(targetFile)
	if err != nil {
		fmt.Printf("目标文件检查失败: %v\n", err)
	} else {
		fmt.Println("目标文件仍然存在")
	}
}

// fileCopyAndBackup 展示文件复制和备份操作
func fileCopyAndBackup(dirPath string) {
	// 创建一个要复制的文件
	originalFile := filepath.Join(dirPath, "original.txt")
	content := "这是原始文件的内容。\n将在示例中被复制和备份。\n"
	err := os.WriteFile(originalFile, []byte(content), 0644)
	if err != nil {
		fmt.Printf("创建原始文件失败: %v\n", err)
		return
	}
	fmt.Printf("创建了原始文件: %s\n", originalFile)

	// 1. 简单复制 (使用io.Copy)
	copyFile := filepath.Join(dirPath, "copy.txt")
	err = copyFileUsingIOCopy(originalFile, copyFile)
	if err != nil {
		fmt.Printf("复制文件失败: %v\n", err)
	} else {
		fmt.Printf("成功复制文件: %s -> %s\n", originalFile, copyFile)
	}

	// 2. 创建带时间戳的备份
	backupFile := filepath.Join(dirPath, fmt.Sprintf("backup-%s.txt",
		time.Now().Format("20060102-150405")))
	err = copyFileUsingIOCopy(originalFile, backupFile)
	if err != nil {
		fmt.Printf("创建备份失败: %v\n", err)
	} else {
		fmt.Printf("成功创建备份: %s\n", backupFile)
	}

	// 3. 修改原文件以演示备份的用途
	updatedContent := content + "这是添加的新内容。\n更新于 " +
		time.Now().Format("2006-01-02 15:04:05") + "\n"
	err = os.WriteFile(originalFile, []byte(updatedContent), 0644)
	if err != nil {
		fmt.Printf("更新原始文件失败: %v\n", err)
	} else {
		fmt.Println("已更新原始文件")
	}

	// 显示所有文件的内容长度进行比较
	showFileSizes(map[string]string{
		"原始文件": originalFile,
		"复制文件": copyFile,
		"备份文件": backupFile,
	})
}

// copyFileUsingIOCopy 使用io.Copy复制文件
func copyFileUsingIOCopy(src, dst string) error {
	// 打开源文件
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("打开源文件失败: %w", err)
	}
	defer srcFile.Close()

	// 创建目标文件
	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("创建目标文件失败: %w", err)
	}
	defer dstFile.Close()

	// 复制内容
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("复制文件内容失败: %w", err)
	}

	// 同步文件内容到磁盘
	err = dstFile.Sync()
	if err != nil {
		return fmt.Errorf("同步文件到磁盘失败: %w", err)
	}

	return nil
}

// showFileSizes 显示多个文件的大小
func showFileSizes(files map[string]string) {
	fmt.Println("\n文件大小比较:")
	for label, path := range files {
		info, err := os.Stat(path)
		if err != nil {
			fmt.Printf("  %-8s: %s (获取信息失败: %v)\n", label, path, err)
			continue
		}
		fmt.Printf("  %-8s: %s (%d 字节)\n", label, path, info.Size())
	}
} 