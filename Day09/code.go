package main

import (
	"AoC2019/Day09/code/intcode"
	"fmt"
)

func main() {
	icpu := intcode.NewIntcodeCPU()
	icpu.Load("input.txt")

	// Part 1
	output, _ := icpu.Run([]int{1})
	fmt.Printf("Part 1: %d\n", output[0])

	// Part 2
	icpu.Reset()
	output, state := icpu.Run([]int{2})
	fmt.Printf("%v %v\n", output, state)
}
