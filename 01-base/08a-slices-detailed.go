// 08a-slices-detailed.go 更详细地展示了Go中切片的特性
package main

import "fmt"

func main() {
	// ---- 1. 基本创建和操作 ----
	fmt.Println("---- 1. 基本创建和操作 ----")
	// 使用make创建切片：类型、长度、容量（可选）
	// s1 长度为3，容量为3（默认等于长度）
	s1 := make([]string, 3)
	fmt.Printf("s1: %v, 长度: %d, 容量: %d\n", s1, len(s1), cap(s1))
	s1[0] = "a"
	s1[1] = "b"
	s1[2] = "c"
	fmt.Printf("s1赋值后: %v, 长度: %d, 容量: %d\n", s1, len(s1), cap(s1))

	// ---- 2. Append 操作和容量变化 ----
	fmt.Println("\n---- 2. Append 操作和容量变化 ----")
	// 当append时，如果容量足够，直接在原底层数组上添加
	// 如果容量不足，会分配新的更大的底层数组，并将元素复制过去
	s1 = append(s1, "d") // s1容量为3，append一个元素后，长度变为4，容量可能会翻倍
	fmt.Printf("s1 append 'd': %v, 长度: %d, 容量: %d\n", s1, len(s1), cap(s1))

	s1 = append(s1, "e", "f")
	fmt.Printf("s1 append 'e','f': %v, 长度: %d, 容量: %d\n", s1, len(s1), cap(s1))

	// ---- 3. 切片的切片 (Slicing a slice) ----
	fmt.Println("\n---- 3. 切片的切片 ----")
	// 切片操作 s[low:high] 创建一个新切片，共享底层数组
	// 新切片的长度是 high - low
	// 新切片的容量是从原切片的low索引到底层数组的末尾
	s2 := []string{"g", "o", "l", "a", "n", "g"}
	fmt.Printf("s2: %v, 长度: %d, 容量: %d\n", s2, len(s2), cap(s2))

	slice1 := s2[1:4] // 元素 o, l, a
	fmt.Printf("slice1 (s2[1:4]): %v, 长度: %d, 容量: %d\n", slice1, len(slice1), cap(slice1))

	slice2 := s2[:3] // 元素 g, o, l
	fmt.Printf("slice2 (s2[:3]): %v, 长度: %d, 容量: %d\n", slice2, len(slice2), cap(slice2))

	slice3 := s2[3:] // 元素 a, n, g
	fmt.Printf("slice3 (s2[3:]): %v, 长度: %d, 容量: %d\n", slice3, len(slice3), cap(slice3))

	// ---- 4. 修改切片会影响底层数组和其他共享切片 ----
	fmt.Println("\n---- 4. 修改切片会影响底层数组和其他共享切片 ----")
	fmt.Printf("修改前 s2: %v\n", s2)
	fmt.Printf("修改前 slice1: %v\n", slice1)
	slice1[0] = "CHANGED_O" // 修改slice1的第一个元素，对应s2的第二个元素
	fmt.Printf("修改后 s2: %v\n", s2)
	fmt.Printf("修改后 slice1: %v\n", slice1)

	// ---- 5. 复制切片 (copy) ----
	fmt.Println("\n---- 5. 复制切片 (copy) ----")
	// copy函数用于将源切片的元素复制到目标切片
	// 复制的元素数量是两个切片长度的较小值
	// copy是深拷贝元素，但如果元素是引用类型，则拷贝的是引用
	original := []string{"apple", "banana", "cherry"}
	destination := make([]string, len(original))
	numCopied := copy(destination, original)
	fmt.Printf("original: %v\n", original)
	fmt.Printf("destination: %v, 复制了 %d 个元素\n", destination, numCopied)
	destination[0] = "apricot" // 修改destination不影响original，因为它们有不同的底层数组
	fmt.Printf("修改destination后 original: %v\n", original)
	fmt.Printf("修改destination后 destination: %v\n", destination)

	// ---- 6. 多维切片 ----
	fmt.Println("\n---- 6. 多维切片 ----")
	// 与08-slices.go中的示例类似
	board := [][]string{
		[]string{"_", "_", "_"},
		[]string{"_", "_", "_"},
		[]string{"_", "_", "_"},
	}
	board[0][0] = "X"
	board[2][2] = "O"
	board[1][2] = "X"
	board[1][0] = "O"
	board[0][2] = "X"

	for i := 0; i < len(board); i++ {
		fmt.Printf("%s\n", board[i])
	}

	// ---- 7. 切片作为函数参数 ----
	fmt.Println("\n---- 7. 切片作为函数参数 ----")
	numbers := []int{1, 2, 3, 4, 5}
	fmt.Printf("调用modifySlice前: %v\n", numbers)
	modifySlice(numbers) // 切片是引用传递（传递的是切片头的副本，但指向同一底层数组）
	fmt.Printf("调用modifySlice后: %v\n", numbers)

	fmt.Printf("调用appendInFunction前: %v, 长度: %d, 容量: %d\n", numbers, len(numbers), cap(numbers))
	newNumbers := appendInFunction(numbers)
	fmt.Printf("调用appendInFunction后 numbers: %v, 长度: %d, 容量: %d\n", numbers, len(numbers), cap(numbers))
	fmt.Printf("调用appendInFunction后 newNumbers: %v, 长度: %d, 容量: %d\n", newNumbers, len(newNumbers), cap(newNumbers))
}

// 修改切片内容的函数
// 由于切片是引用类型，函数内对切片元素的修改会影响到调用者
func modifySlice(s []int) {
	if len(s) > 0 {
		s[0] = 999
	}
}

// 在函数内对切片进行append操作
// 如果append导致底层数组重新分配，原切片不会感知到这个新数组
// 所以通常append操作的函数需要返回新的切片
func appendInFunction(s []int) []int {
	s = append(s, 66)
	fmt.Printf("在appendInFunction内部 s: %v, 长度: %d, 容量: %d\n", s, len(s), cap(s))
	return s
} 