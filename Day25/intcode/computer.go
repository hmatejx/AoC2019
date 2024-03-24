package intcode

import (
	"log"
	"os"
	"strconv"
	"strings"
)

type IntcodeCPU struct {
	code   []int
	memory map[int]int
	rbase  int
	ip     int
}

func NewIntcodeCPU() *IntcodeCPU {
	return &IntcodeCPU{code: nil, memory: map[int]int{}, rbase: 0, ip: 0}
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
	icpu.memory = map[int]int{}
	for i := 0; i < len(icpu.code); i++ {
		icpu.memory[i] = icpu.code[i]
	}
	icpu.ip = 0
	icpu.rbase = 0
}

func (icpu *IntcodeCPU) Copy() *IntcodeCPU {
	new_cpu := NewIntcodeCPU()
	new_cpu.code = make([]int, len(icpu.code))
	copy(new_cpu.code, icpu.code)
	for k, v := range icpu.memory {
		new_cpu.memory[k] = v
	}
	new_cpu.ip, new_cpu.rbase = icpu.ip, icpu.rbase
	return new_cpu
}

func (icpu *IntcodeCPU) Run(input []int) ([]int, int) {
	// auxiliary function
	_addr_ := func(mode, i int) int {
		switch mode {
		case 1:
			return i
		case 2:
			return icpu.memory[i] + icpu.rbase
		default:
			return icpu.memory[i]
		}
	}
	output := []int{}
loop:
	for icpu.ip < len(icpu.memory) {
		// decode the opcode
		opcode := icpu.memory[icpu.ip] % 100
		// and parameter modes
		mode_a := icpu.memory[icpu.ip] / 10000
		mode_b := (icpu.memory[icpu.ip] % 10000) / 1000
		mode_c := (icpu.memory[icpu.ip] % 1000) / 100
		var p1, p2 int
		switch opcode {
		// addition, multiplication
		case 1, 2:
			dest_addr := _addr_(mode_a, icpu.ip+3)
			p1, p2 = icpu.memory[_addr_(mode_c, icpu.ip+1)], icpu.memory[_addr_(mode_b, icpu.ip+2)]
			if opcode == 1 {
				icpu.memory[dest_addr] = p1 + p2
			} else {
				icpu.memory[dest_addr] = p1 * p2
			}
			icpu.ip += 4
		// input
		case 3:
			// return when waiting for input
			if len(input) == 0 {
				return output, opcode
			}
			// pop the first element from the input
			current_input := input[0]
			input = input[1:]
			icpu.memory[_addr_(mode_c, icpu.ip+1)] = current_input
			icpu.ip += 2
		// output
		case 4:
			output = append(output, icpu.memory[_addr_(mode_c, icpu.ip+1)])
			icpu.ip += 2
		// jump-if-true, jump-if-false
		case 5, 6:
			p1 = icpu.memory[_addr_(mode_c, icpu.ip+1)]
			if (p1 != 0 && opcode == 5) || (p1 == 0 && opcode == 6) {
				icpu.ip = icpu.memory[_addr_(mode_b, icpu.ip+2)]
			} else {
				icpu.ip += 3
			}
		// less-than, equals
		case 7, 8:
			dest_addr := _addr_(mode_a, icpu.ip+3)
			p1, p2 = icpu.memory[_addr_(mode_c, icpu.ip+1)], icpu.memory[_addr_(mode_b, icpu.ip+2)]
			if (p1 < p2 && opcode == 7) || (p1 == p2 && opcode == 8) {
				icpu.memory[dest_addr] = 1
			} else {
				icpu.memory[dest_addr] = 0
			}
			icpu.ip += 4
		case 9:
			icpu.rbase += icpu.memory[_addr_(mode_c, icpu.ip+1)]
			icpu.ip += 2
		case 99:
			break loop
		}
	}
	return output, 99
}

func From_ASCII(line string) []int {
	output := make([]int, len(line)+1)
	for i, c := range line {
		output[i] = int(c)
	}
	output[len(line)] = 10
	return output
}

func To_ASCII(output []int) string {
	output_str := ""
	for _, o := range output {
		output_str += string(rune(o))
	}
	return output_str
}
