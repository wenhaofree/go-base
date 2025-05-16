// 09-maps.go 展示了Go中映射(maps)的使用
package main

import "fmt"

func main() {
	// 要创建一个空map，使用内建的make函数：make(map[key-type]val-type)
	m := make(map[string]int)

	// 使用典型的 name[key] = val 语法来设置键值对
	m["k1"] = 7
	m["k2"] = 13
	
	// 打印map。例如，使用fmt.Println打印一个map，会输出所有键值对
	fmt.Println("映射:", m)

	// 使用name[key]来获取一个键的值
	v1 := m["k1"]
	fmt.Println("v1:", v1)

	// 内建函数len可以返回一个映射中键值对的数量
	fmt.Println("长度:", len(m))

	// 内建的delete函数从映射中移除键值对
	delete(m, "k2")
	fmt.Println("删除后:", m)

	// 当从一个映射中取值时，可选的第二个返回值表示该键是否在映射中
	// 这可以用来消除缺少该键与值为零值的歧义，例如0或""
	_, prs := m["k2"]
	fmt.Println("存在:", prs)

	// 你也可以通过这种语法在同一行声明和初始化一个新的映射
	n := map[string]int{"foo": 1, "bar": 2}
	fmt.Println("声明:", n)
} 