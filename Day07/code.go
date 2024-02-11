package main

import (
	"AoC2019/Day07/code/intcode"
	"fmt"
)

func permutations(data []int) [][]int {
	permutation := make([]int, len(data))
	indexInUse := make([]bool, len(data))

	var ret [][]int
	var f func(idx int)

	f = func(idx int) {
		if idx >= len(data) {
			arr := make([]int, len(data))
			copy(arr, permutation)
			ret = append(ret, arr)
			return
		}
		for i := 0; i < len(data); i++ {
			if !indexInUse[i] {
				indexInUse[i] = true
				permutation[idx] = data[i]
				f(idx + 1)
				indexInUse[i] = false
			}
		}
	}
	f(0)

	return ret
}

func runSequence(modules []*intcode.IntcodeCPU, phasing []int, reset bool) int {
	output := []int{0}
	for i := 0; i < 5; i++ {
		output, _ = modules[i].Run([]int{phasing[i], output[0]})
		if reset {
			modules[i].Reset()
		}
	}
	return output[0]
}

func runSequenceWithFeedback(modules []*intcode.IntcodeCPU, phasing []int, reset bool) int {
	var ret int
	output := []int{0}
	first_time := true
	// feedback loop
	for {
		for i := 0; i < 5; i++ {
			if first_time {
				output, ret = modules[i].Run(append([]int{phasing[i]}, output...))
			} else {
				output, ret = modules[i].Run(output)
			}
		}
		first_time = false
		if ret == 99 {
			break
		}
	}
	// reset all modules back to the initial state
	if reset {
		for i := 0; i < 5; i++ {
			modules[i].Reset()
		}
	}
	return output[0]
}

func main() {
	// load the intcode into each module
	var modules []*intcode.IntcodeCPU
	for i := 0; i < 5; i++ {
		modules = append(modules, intcode.NewIntcodeCPU())
		modules[i].Load("input.txt")
	}

	// Part 1
	phase := permutations([]int{0, 1, 2, 3, 4})
	max_out := -1 << 63
	for _, phasing := range phase {
		output := runSequence(modules, phasing, true)
		if output > max_out {
			max_out = output
		}
	}
	fmt.Printf("Part 1: %d\n", max_out)

	// Part 2
	phase = permutations([]int{5, 6, 7, 8, 9})
	max_out = -1 << 63
	for _, phasing := range phase {
		output := runSequenceWithFeedback(modules, phasing, true)
		if output > max_out {
			max_out = output
		}
	}
	fmt.Printf("Part 2: %d", max_out)
}
