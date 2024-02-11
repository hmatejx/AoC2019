package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func read_input() []int {
	content, err := os.ReadFile("input.txt")
	if err != nil {
		log.Fatalln(err)
	}
	numbers := []int{}
	for _, line := range strings.Split(string(content), ",") {
		n, _ := strconv.Atoi(line)
		numbers = append(numbers, n)
	}
	return numbers
}

func run(code []int, input []int) ([]int, []int) {
	// load program to memory
	memory := make([]int, len(code))
	copy(memory, code)

	// setup input & output
	output := []int{}
	input_idx := 0

	// auxiliary function
	get_par := func(mode, i int) int {
		if mode == 1 {
			return memory[i]
		}
		return memory[memory[i]]
	}

	ip := 0
loop:
	for ip < len(memory) {
		// decode the opcode
		opcode := memory[ip] % 100
		// and parameter modes
		// mode_a := code[ip] / 10000
		mode_b := (memory[ip] % 10000) / 1000
		mode_c := (memory[ip] % 1000) / 100
		var p1, p2 int
		switch opcode {
		// addition, multiplication
		case 1, 2:
			if opcode == 1 {
				memory[memory[ip+3]] = get_par(mode_c, ip+1) + get_par(mode_b, ip+2)
			} else {
				memory[memory[ip+3]] = get_par(mode_c, ip+1) * get_par(mode_b, ip+2)
			}
			ip += 4
		// input
		case 3:
			memory[memory[ip+1]] = input[input_idx]
			input_idx++
			ip += 2
		// output
		case 4:
			output = append(output, get_par(mode_c, ip+1))
			ip += 2
		// jump-if-true, jump-if-false
		case 5, 6:
			p1 = get_par(mode_c, ip+1)
			if (p1 != 0 && opcode == 5) || (p1 == 0 && opcode == 6) {
				ip = get_par(mode_b, ip+2)
			} else {
				ip += 3
			}
		// less-than, equals
		case 7, 8:
			p1, p2 = get_par(mode_c, ip+1), get_par(mode_b, ip+2)
			if (p1 < p2 && opcode == 7) || (p1 == p2 && opcode == 8) {
				memory[memory[ip+3]] = 1
			} else {
				memory[memory[ip+3]] = 0
			}
			ip += 4
		case 99:
			break loop
		}
	}
	return memory, output
}

func main() {
	code := read_input()
	_, output := run(code, []int{1})
	fmt.Printf("Part 1: %v\n", output)
	_, output = run(code, []int{5})
	fmt.Printf("Part 2: %d\n", output[0])
}
