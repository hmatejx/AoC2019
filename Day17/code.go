package main

import (
	"AoC2019/Day17/code/intcode"
	"fmt"
)

func output_to_scaffold(output []int) map[[2]int]int {
	pos := [2]int{0, 0}
	scaffold := map[[2]int]int{}
	for _, o := range output {
		if o == 35 {
			scaffold[pos] = 1
		}
		fmt.Print(string(rune(o)))
		if o != 10 {
			pos[0]++
		} else {
			pos[1]++
			pos[0] = 0
		}
	}
	return scaffold
}

func from_ascii(line string) []int {
	output := make([]int, len(line)+1)
	for i, c := range line {
		output[i] = int(c)
	}
	output[len(line)] = 10
	return output
}

func to_ascii(output []int) string {
	output_str := ""
	for _, o := range output {
		output_str += string(rune(o))
	}
	return output_str
}

func main() {
	robot := intcode.NewIntcodeCPU()
	robot.Load("input.txt")

	// Part 1
	output, _ := robot.Run([]int{})
	scaffold := output_to_scaffold(output)
	alignment := 0
	for k := range scaffold {
		if scaffold[[2]int{k[0] - 1, k[1]}] == 1 &&
			scaffold[[2]int{k[0] + 1, k[1]}] == 1 &&
			scaffold[[2]int{k[0], k[1] - 1}] == 1 &&
			scaffold[[2]int{k[0], k[1] + 1}] == 1 {
			alignment += k[0] * k[1]
		}
	}
	fmt.Printf("Part 1: %d\n", alignment)

	// Part 2 (path compression by hand)
	prog := [5]string{
		"A,A,B,C,A,C,B,C,A,B",  // Main
		"L,4,L,10,L,6",         // A
		"L,6,L,4,R,8,R,8",      // B
		"L,6,R,8,L,10,L,8,L,8", // C
		"n",
	}
	robot.Reset()
	robot.Poke(0, 2)
	output, _ = robot.Run([]int{})
	fmt.Printf("%s", to_ascii(output))
	for _, line := range prog {
		fmt.Printf("%s\n", line)
		output, _ = robot.Run(from_ascii(line))
		fmt.Printf("%s", to_ascii(output))
	}
	fmt.Printf("Part 2: %d\n", output[len(output)-1])
}
