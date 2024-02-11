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
	for _, line := range strings.Split(string(content), "\r\n") {
		n, _ := strconv.Atoi(line)
		numbers = append(numbers, n)
	}
	return numbers
}

func fuel_calc(x int) int {
	return max(x/3-2, 0)
}

func main() {
	weights := read_input()

	// Part 1
	fuel := 0
	for _, n := range weights {
		fuel += fuel_calc(n)
	}
	fmt.Printf("Part 1: %d\n", fuel)

	// Part 2
	fuel = 0
	for _, n := range weights {
		for inc := fuel_calc(n); inc > 0; inc = fuel_calc(inc) {
			fuel += inc
		}
	}
	fmt.Printf("Part 2: %d\n", fuel)
}
