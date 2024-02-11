package main

import (
	"AoC2019/Day11/code/intcode"
	"fmt"
)

func paint(robot *intcode.IntcodeCPU, initial_panel int) map[[2]int]int {
	hull := map[[2]int]int{}
	position := [2]int{0, 0}
	hull[position] = initial_panel
	direction := [2]int{0, -1}
	for {
		output, state := robot.Run([]int{hull[position]})
		hull[position] = output[0]
		if output[1] == 1 {
			direction = [2]int{-direction[1], direction[0]}
		} else {
			direction = [2]int{direction[1], -direction[0]}
		}
		position[0] += direction[0]
		position[1] += direction[1]
		if state == 99 {
			break
		}
	}
	return hull
}

func display(hull map[[2]int]int) {
	x0, x1, y0, y1 := 1<<63-1, -1<<63, 1<<63-1, -1<<63
	for k, _ := range hull {
		if k[0] < x0 {
			x0 = k[0]
		}
		if k[0] > x1 {
			x1 = k[0]
		}
		if k[1] < y0 {
			y0 = k[1]
		}
		if k[1] > y1 {
			y1 = k[1]
		}
	}
	fmt.Println("Part 2:")
	for y := y0; y <= y1; y++ {
		for x := x0; x <= x1; x++ {
			if hull[[2]int{x, y}] == 1 {
				fmt.Print("█")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}

func main() {
	robot := intcode.NewIntcodeCPU()
	robot.Load("input.txt")

	// Part 1
	hull := paint(robot, 0)
	fmt.Printf("Part 1: %d\n", len(hull))

	// Part 2
	robot.Reset()
	hull = paint(robot, 1)
	display(hull)
}
