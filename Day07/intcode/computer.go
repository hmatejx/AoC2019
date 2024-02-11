package intcode

import (
	"log"
	"os"
	"strconv"
	"strings"
)

type IntcodeCPU struct {
	code   []int
	memory []int
	input  []int
	ip     int
}

func NewIntcodeCPU() *IntcodeCPU {
	return &IntcodeCPU{code: []int{}, memory: []int{}, input: []int{}, ip: 0}
}

func (icpu *IntcodeCPU) Load(filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalln(err)
	}
	numbers := []int{}
	for _, line := range strings.Split(string(content), ",") {
		n, _ := strconv.Atoi(line)
		numbers = append(numbers, n)
	}
	icpu.code = numbers
	icpu.Reset()
}

func (icpu *IntcodeCPU) Reset() {
	if len(icpu.memory) != len(icpu.code) {
		icpu.memory = make([]int, len(icpu.code))
	}
	copy(icpu.memory, icpu.code)
	icpu.ip = 0
}

func (icpu *IntcodeCPU) Run(input []int) ([]int, int) {
	// setup input & output
	var current_input int
	icpu.input = make([]int, len(input))
	copy(icpu.input, input)
	output := []int{}

	// auxiliary function
	get_par := func(mode, i int) int {
		if mode == 1 {
			return icpu.memory[i]
		}
		return icpu.memory[icpu.memory[i]]
	}

loop:
	for icpu.ip < len(icpu.memory) {
		// decode the opcode
		opcode := icpu.memory[icpu.ip] % 100
		// and parameter modes
		// mode_a := code[ip] / 10000
		mode_b := (icpu.memory[icpu.ip] % 10000) / 1000
		mode_c := (icpu.memory[icpu.ip] % 1000) / 100
		var p1, p2 int
		switch opcode {
		// addition, multiplication
		case 1, 2:
			if opcode == 1 {
				icpu.memory[icpu.memory[icpu.ip+3]] = get_par(mode_c, icpu.ip+1) + get_par(mode_b, icpu.ip+2)
			} else {
				icpu.memory[icpu.memory[icpu.ip+3]] = get_par(mode_c, icpu.ip+1) * get_par(mode_b, icpu.ip+2)
			}
			icpu.ip += 4
		// input
		case 3:
			// return when waiting for input
			if len(icpu.input) == 0 {
				return output, opcode
			}
			// pop the first element from the input
			current_input, icpu.input = icpu.input[0], icpu.input[1:]
			icpu.memory[icpu.memory[icpu.ip+1]] = current_input
			icpu.ip += 2
		// output
		case 4:
			output = append(output, get_par(mode_c, icpu.ip+1))
			icpu.ip += 2
		// jump-if-true, jump-if-false
		case 5, 6:
			p1 = get_par(mode_c, icpu.ip+1)
			if (p1 != 0 && opcode == 5) || (p1 == 0 && opcode == 6) {
				icpu.ip = get_par(mode_b, icpu.ip+2)
			} else {
				icpu.ip += 3
			}
		// less-than, equals
		case 7, 8:
			p1, p2 = get_par(mode_c, icpu.ip+1), get_par(mode_b, icpu.ip+2)
			if (p1 < p2 && opcode == 7) || (p1 == p2 && opcode == 8) {
				icpu.memory[icpu.memory[icpu.ip+3]] = 1
			} else {
				icpu.memory[icpu.memory[icpu.ip+3]] = 0
			}
			icpu.ip += 4
		case 99:
			break loop
		}
	}
	return output, 99
}
