package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func load_input(filename string) []int {
	content, _ := os.ReadFile(filename)
	lines := strings.Split(string(content), "\r\n")
	str := strings.Split(lines[0], "")
	numbers := make([]int, len(str))
	for i, c := range str {
		numbers[i], _ = strconv.Atoi(c)
	}
	return numbers
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func sum(array []int, i0 int, i1 int) (res int) {
	for i := i0; i <= i1; i++ {
		res += array[i]
	}
	return
}

func FFT(input []int, phases int, offset int) []int {
	length := len(input)
	output := make([]int, length)
	copy(output, input)
	for phase := 0; phase < phases; phase++ {
		for fast, slow, digit := 0, 0, 1+offset; digit <= length; digit++ {
			if digit > length/2 {
				// fast path for digits in second half
				if fast == 0 {
					fast = sum(output, digit-1, length-1)
				}
				x := output[digit-1]
				output[digit-1] = fast % 10
				fast -= x
			} else {
				// regular path
				for p, i0 := 0, digit-1; i0 < length; {
					i1 := min(i0+digit-1, length-1)
					if p%2 == 0 {
						slow += sum(output, i0, i1)
					} else {
						slow -= sum(output, i0, i1)
					}
					p++
					i0 += 2 * digit
				}
				output[digit-1] = abs(slow) % 10
				slow = 0
			}
		}
	}
	return output
}

func collapse(list []int) int {
	str := strings.Replace(fmt.Sprint(list), " ", "", -1)
	val, _ := strconv.Atoi(str[1 : len(str)-1])
	return val
}

func main() {
	input := load_input("input.txt")

	// Part 1
	output := FFT(input, 100, 0)
	fmt.Printf("Part 1: %d\n", collapse(output[:8]))

	// Part 2
	expand := 10000
	offset := collapse(input[:7])
	expanded_input := make([]int, expand*len(input))
	for i := 0; i < expand; i++ {
		copy(expanded_input[i*len(input):(i+1)*len(input)], input)
	}
	output = FFT(expanded_input, 100, offset)
	fmt.Printf("Part 2: %d\n", collapse(output[offset:offset+8]))
}
