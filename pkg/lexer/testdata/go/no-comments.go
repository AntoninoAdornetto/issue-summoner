package go_test

import "fmt"

func add(a, b int) {
	return a + b
}

func subtract(a, b int) {
	return a - b
}

func main() {
	x := add(60, 9)
	y := subtract(70, 1)

	if x == y {
		fmt.Printf("Hello, World\n")
	}

	fmt.Printf("No Comments in this source file")
}
