package main

import "fmt"

// Map map a integer array to another function
func Map(vs []int, f func(int) int) []int {
	vsm := make([]int, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

func square(x int) int {
	return x * 2
}

func main() {
	a := []int{1, 2, 3, 4, 5, 6, 10}

	square := Map(a, square)
	for _, val := range square {
		fmt.Printf("%d ", val)
	}
	fmt.Println()
}
