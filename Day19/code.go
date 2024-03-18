package main

import (
	"AoC2019/Day19/code/intcode"
	"fmt"
)

func check(robot *intcode.IntcodeCPU, x, y int) bool {
	robot.Reset()
	res, _ := robot.Run([]int{x, y})
	return res[0] > 0
}

func main() {
	robot := intcode.NewIntcodeCPU()
	robot.Load("input.txt")
	total := 0
	size := 50
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if check(robot, x, y) {
				total += 1
			}
		}
	}
	fmt.Printf("Part 1: %d\n", total)

	x, y := 0, 0
	for !check(robot, x+99, y) {
		y += 1
		for !check(robot, x, y+99) {
			x += 1
		}
	}
	fmt.Printf("%d\n", x*100*100+y)
}
