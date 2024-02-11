package main

import (
	"fmt"
	"strconv"
)

func is_valid(n int, rule2 bool) bool {
	str := strconv.Itoa(n)
	same_adjacent := false
	increasing := true
	for i := 0; i < len(str)-1; i++ {
		if (str[i] == str[i+1] && !rule2) || ((str[i] == str[i+1]) &&
			(i == len(str)-2 || str[i] != str[i+2]) &&
			(i == 0 || str[i] != str[i-1])) {
			same_adjacent = true
		}
		if str[i+1] < str[i] {
			increasing = false
		}
	}
	return increasing && same_adjacent
}

func main() {
	input := [2]int{402328, 864247}

	valid1, valid2 := 0, 0
	for n := input[0]; n <= input[1]; n++ {
		if is_valid(n, false) {
			valid1++
		}
		if is_valid(n, true) {
			valid2++
		}
	}
	fmt.Printf("Part 1: %d\n", valid1)
	fmt.Printf("Part 2: %d\n", valid2)
}
