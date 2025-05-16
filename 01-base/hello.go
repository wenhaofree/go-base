// Package main is the entry point for Go applications
package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// calculateSum is a function that takes a slice of integers and returns their sum
func calculateSum(numbers []int) int {
	sum := 0
	for _, num := range numbers {
		sum += num
	}
	return sum
}

func main() {
	// Print a welcome message
	fmt.Println("Welcome to Go Programming Examples!")
	
	// Variable declarations with different types
	var name string = "Gopher"      // Explicit type declaration
	age := 10                       // Type inference
	var isActive = true             // Type inference with var
	var height float64              // Zero-value initialization
	height = 1.73
	
	// Print variable values
	fmt.Println("\n--- Variable Examples ---")
	fmt.Printf("Name: %s\n", name)
	fmt.Printf("Age: %d\n", age)
	fmt.Printf("Active: %v\n", isActive)
	fmt.Printf("Height: %.2f meters\n", height)
	
	// Working with slices
	fmt.Println("\n--- Slice Examples ---")
	
	// Initialize random number generator with current time as seed
	rand.Seed(time.Now().UnixNano())
	
	// Create and populate a slice of random numbers
	numbers := make([]int, 5)
	for i := range numbers {
		numbers[i] = rand.Intn(100) // Random number between 0 and 99
	}
	fmt.Printf("Generated numbers: %v\n", numbers)
	
	// Call a function and store its return value
	sum := calculateSum(numbers)
	fmt.Printf("Sum of numbers: %d\n", sum)
	
	// Demonstrating slice operations
	fruits := []string{"Apple", "Banana", "Cherry", "Durian", "Elderberry"}
	fmt.Printf("All fruits: %v\n", fruits)
	fmt.Printf("First two fruits: %v\n", fruits[:2])
	fmt.Printf("Last three fruits: %v\n", fruits[2:])
	
	// Append to slice
	fruits = append(fruits, "Fig", "Grape")
	fmt.Printf("After appending: %v\n", fruits)
	
	// Control structures example: if-else
	fmt.Println("\n--- Control Structures ---")
	if sum > 200 {
		fmt.Println("Sum is greater than 200")
	} else if sum > 100 {
		fmt.Println("Sum is between 101 and 200")
	} else {
		fmt.Println("Sum is 100 or less")
	}
	
	// For loop examples
	fmt.Println("\n--- Loop Examples ---")
	fmt.Println("Counting from 1 to 5:")
	for i := 1; i <= 5; i++ {
		fmt.Printf("%d ", i)
	}
	fmt.Println()
	
	// Range-based for loop with slice
	fmt.Println("Iterating through fruits:")
	for index, fruit := range fruits {
		fmt.Printf("%d: %s\n", index, fruit)
	}
	
	// String manipulation example
	fmt.Println("\n--- String Manipulation ---")
	message := "  Hello, Go Programming World!  "
	fmt.Printf("Original: %q\n", message)
	fmt.Printf("Trimmed: %q\n", strings.TrimSpace(message))
	fmt.Printf("Uppercase: %q\n", strings.ToUpper(message))
	fmt.Printf("Contains 'Go': %v\n", strings.Contains(message, "Go"))
	fmt.Printf("Replace 'World' with 'Universe': %q\n", 
		strings.Replace(message, "World", "Universe", 1))
	
	// Time example
	fmt.Println("\n--- Current Time ---")
	currentTime := time.Now()
	fmt.Printf("Current time: %s\n", currentTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("Unix timestamp: %d\n", currentTime.Unix())
}

