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

func run(p1, p2 int) int {
	code := read_input()
	code[1] = p1
	code[2] = p2
loop:
	for ip := 0; ip < len(code); ip += 4 {
		switch code[ip] {
		case 1:
			code[code[ip+3]] = code[code[ip+1]] + code[code[ip+2]]
		case 2:
			code[code[ip+3]] = code[code[ip+1]] * code[code[ip+2]]
		case 99:
			break loop
		}
	}
	return code[0]
}

func main() {
	fmt.Printf("Part 1: %d\n", run(1, 12))
	// searching for 19690720 <- run(49, 25) ... linear function, k*p1 + p2
	fmt.Printf("Part 2: %d\n", 49*100+25)
}
