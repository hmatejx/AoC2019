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
	input  []int
	rbase  int
	ip     int
}

func NewIntcodeCPU() *IntcodeCPU {
	return &IntcodeCPU{code: []int{}, memory: map[int]int{}, input: []int{}, rbase: 0, ip: 0}
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

func (icpu *IntcodeCPU) Peek(address int) int {
	return icpu.memory[address]
}

func (icpu *IntcodeCPU) Poke(address, value int) {
	icpu.memory[address] = value
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
	// setup input & output
	var current_input int
	icpu.input = make([]int, len(input))
	copy(icpu.input, input)
	output := []int{}

	// auxiliary function
	get_par := func(mode, i int) int {
		switch mode {
		case 1:
			return icpu.memory[i]
		case 2:
			return icpu.memory[icpu.memory[i]+icpu.rbase]
		default:
			return icpu.memory[icpu.memory[i]]
		}
	}
	get_write_addr := func(mode, i int) int {
		if mode == 2 {
			return icpu.memory[i] + icpu.rbase
		}
		return icpu.memory[i]
	}

loop:
	for icpu.ip < len(icpu.memory) {
		// decode the opcode
		//fmt.Printf("ip: %d,\trbase: %d, \topcode: %d\t", icpu.ip, icpu.rbase, icpu.memory[icpu.ip])
		opcode := icpu.memory[icpu.ip] % 100
		// and parameter modes
		mode_a := icpu.memory[icpu.ip] / 10000
		mode_b := (icpu.memory[icpu.ip] % 10000) / 1000
		mode_c := (icpu.memory[icpu.ip] % 1000) / 100
		//fmt.Printf("mode A: %d, mode B: %d, mode C: %d\n", mode_a, mode_b, mode_c)
		var p1, p2 int
		switch opcode {
		// addition, multiplication
		case 1, 2:
			addr := get_write_addr(mode_a, icpu.ip+3)
			if opcode == 1 {
				icpu.memory[addr] = get_par(mode_c, icpu.ip+1) + get_par(mode_b, icpu.ip+2)
			} else {
				icpu.memory[addr] = get_par(mode_c, icpu.ip+1) * get_par(mode_b, icpu.ip+2)
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
			icpu.memory[get_write_addr(mode_c, icpu.ip+1)] = current_input
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
			addr := icpu.memory[icpu.ip+3]
			if mode_a == 2 {
				addr += icpu.rbase
			}
			if (p1 < p2 && opcode == 7) || (p1 == p2 && opcode == 8) {
				icpu.memory[addr] = 1
			} else {
				icpu.memory[addr] = 0
			}
			icpu.ip += 4
		case 9:
			icpu.rbase += get_par(mode_c, icpu.ip+1)
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
