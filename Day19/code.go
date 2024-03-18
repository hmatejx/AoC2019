package main

import (
	"AoC2019/Day19/code/intcode"
	"fmt"
)

var robot *intcode.IntcodeCPU

func check(x, y int) bool {
	robot.Reset()
	res, _ := robot.Run([]int{x, y})
	return res[0] > 0
}

func main() {
	robot = intcode.NewIntcodeCPU()
	robot.Load("input.txt")

	// part 1
	total := 0
	for y := 0; y < 50; y++ {
		for x := 0; x < 50; x++ {
			if check(x, y) {
				total += 1
			}
		}
	}
	fmt.Printf("Part 1: %d\n", total)

	// part 2
	x, y := 0, 0
	for !check(x+99, y) {
		y += 1
		for !check(x, y+99) {
			x += 1
		}
	}
	fmt.Printf("Part 2: %d\n", x*100*100+y)
}
