// 09a-maps-detailed.go 更详细地展示了Go中映射(maps)的特性
package main

import (
	"fmt"
	"sort" // 用于演示按键排序遍历
)

func main() {
	// ---- 1. 基本创建、赋值和访问 ----
	fmt.Println("---- 1. 基本创建、赋值和访问 ----")
	// 使用make创建映射
	ages := make(map[string]int)
	ages["Alice"] = 30
	ages["Bob"] = 25
	ages["Charlie"] = 35
	fmt.Printf("ages 映射: %v\n", ages)
	fmt.Printf("Alice的年龄: %d\n", ages["Alice"])

	// 声明并初始化映射
	balances := map[string]float64{
		"USD": 123.45,
		"EUR": 234.56,
	}
	fmt.Printf("balances 映射: %v\n", balances)

	// ---- 2. 获取值和检查键是否存在 ----
	fmt.Println("\n---- 2. 获取值和检查键是否存在 ----")
	// 获取一个不存在的键，会返回该值类型的零值
	davidAge := ages["David"] // David不存在，返回int的零值0
	fmt.Printf("David的年龄 (不存在时返回零值): %d\n", davidAge)

	// 安全地检查键是否存在：value, exists := m[key]
	age, exists := ages["Bob"]
	if exists {
		fmt.Printf("Bob存在，年龄: %d\n", age)
	} else {
		fmt.Println("Bob不存在")
	}

	age, exists = ages["David"]
	if exists {
		fmt.Printf("David存在，年龄: %d\n", age)
	} else {
		fmt.Println("David不存在")
	}

	// ---- 3. 删除键值对 (delete) ----
	fmt.Println("\n---- 3. 删除键值对 ----")
	fmt.Printf("删除前 ages: %v\n", ages)
	delete(ages, "Alice") // 删除Alice
	fmt.Printf("删除Alice后 ages: %v\n", ages)
	delete(ages, "Unknown") // 删除一个不存在的键不会报错
	fmt.Printf("删除Unknown后 ages: %v\n", ages)

	// ---- 4. 遍历映射 (range) ----
	fmt.Println("\n---- 4. 遍历映射 (迭代顺序不确定) ----")
	// 映射的迭代顺序是不确定的
	for name, age := range ages {
		fmt.Printf("%s 的年龄是 %d\n", name, age)
	}

	// ---- 5. 按键排序遍历映射 ----
	fmt.Println("\n---- 5. 按键排序遍历映射 ----")
	// 如果需要按特定顺序遍历，先提取键到切片，排序切片，再遍历
	names := make([]string, 0, len(ages)) // 创建一个切片用于存储键
	for name := range ages {
		names = append(names, name)
	}
	sort.Strings(names) // 对键进行排序
	fmt.Println("按名字字母顺序:")
	for _, name := range names {
		fmt.Printf("%s 的年龄是 %d\n", name, ages[name])
	}

	// ---- 6. 映射作为函数参数 ----
	fmt.Println("\n---- 6. 映射作为函数参数 ----")
	fmt.Printf("调用modifyMap前 balances: %v\n", balances)
	modifyMap(balances)
	fmt.Printf("调用modifyMap后 balances: %v\n", balances)

	// ---- 7. 键的类型 ----
	fmt.Println("\n---- 7. 键的类型 ----")
	// 键必须是可比较的类型。切片、映射、函数不能作为键
	// m := make(map[[]int]string) // 编译错误: invalid map key type []int
	type Person struct {
		Name string
		Age  int
	}
	// 结构体可以作为键，前提是其所有字段都是可比较的
	personMap := make(map[Person]string)
	p1 := Person{"Eve", 28}
	p2 := Person{"Frank", 40}
	personMap[p1] = "Engineer"
	personMap[p2] = "Doctor"
	fmt.Printf("personMap: %v\n", personMap)

}

// 修改映射内容的函数
// 由于映射是引用类型，函数内对映射的修改会影响到调用者
func modifyMap(m map[string]float64) {
	m["JPY"] = 15000.75
	delete(m, "EUR")
} 