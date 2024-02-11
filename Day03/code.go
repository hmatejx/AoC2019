package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

func read_input() ([][2]int, [][2]int) {
	content, err := os.ReadFile("input.txt")
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(content), "\r\n")
	p1, p2 := strings.Split(lines[0], ","), strings.Split(lines[1], ",")

	direction := func(d uint8) (int, int) {
		switch d {
		case uint8('R'):
			return 1, 0
		case uint8('L'):
			return -1, 0
		case uint8('D'):
			return 0, 1
		}
		return 0, -1
	}

	make_wire := func(path []string) [][2]int {
		x, y := 0, 0
		wire := [][2]int{}
		wire = append(wire, [2]int{x, y})
		for _, step := range path {
			dx, dy := direction(step[0])
			length, _ := strconv.Atoi(step[1:])
			x, y = x+dx*length, y+dy*length
			wire = append(wire, [2]int{x, y})
		}
		return wire
	}

	return make_wire(p1), make_wire(p2)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	wire1, wire2 := read_input()

	var d, wd int
	distances := []int{}
	wire_distances := []int{}
	l1 := 0
	for i := 0; i < len(wire1)-1; i++ {
		a1, b1 := wire1[i], wire1[i+1]
		l2 := 0
		for j := 0; j < len(wire2)-1; j++ {
			a2, b2 := wire2[j], wire2[j+1]
			if max(a1[0], b1[0]) >= min(a2[0], b2[0]) && min(a1[0], b1[0]) <= max(a2[0], b2[0]) &&
				max(a1[1], b1[1]) >= min(a2[1], b2[1]) && min(a1[1], b1[1]) <= max(a2[1], b2[1]) {
				dx1, dy1 := b1[0]-a1[0], b1[1]-a1[1]
				dx2, dy2 := b2[0]-a2[0], b2[1]-a2[1]
				if dx1 == 0 && dy2 == 0 {
					d = abs(a1[0]) + abs(a2[1])
					wd = l1 + abs(a2[1]-a1[1]) + l2 + abs(a1[0]-a2[0])
				} else if dx2 == 0 && dy1 == 0 {
					d = abs(a2[0]) + abs(a1[1])
					wd = l1 + abs(a2[0]-a1[0]) + l2 + abs(a1[1]-a2[1])
				}
				if d > 0 {
					distances = append(distances, d)
					wire_distances = append(wire_distances, wd)
				}
			}
			l2 += abs(b2[0]-a2[0]) + abs(b2[1]-a2[1])
		}
		l1 += abs(b1[0]-a1[0]) + abs(b1[1]-a1[1])
	}
	fmt.Printf("%v\n", distances)
	fmt.Printf("Part 1: %d\n", slices.Min(distances))
	fmt.Printf("Part 2: %d\n", slices.Min(wire_distances))
}
