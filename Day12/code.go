package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Moon struct {
	position [3]int
	velocity [3]int
}

func load_moons(filename string) []*Moon {
	content, _ := os.ReadFile(filename)
	moons := []*Moon{}
	lines := strings.Split(string(content), "\r\n")
	re := regexp.MustCompile(`([a-z]=|>|<)`)
	for _, line := range lines {
		pos := [3]int{}
		strpos := strings.Split(string(re.ReplaceAll([]byte(line), []byte{})), ", ")
		for i, s := range strpos {
			n, _ := strconv.Atoi(s)
			pos[i] = n
		}
		moons = append(moons, &Moon{position: pos})
	}
	return moons
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

func LCM(a, b int) int {
	return a * b / GCD(a, b)
}

func apply_gravity(moons *[]*Moon) {
	n_moons := len(*moons)
	for i := 0; i < n_moons; i++ {
		m1 := (*moons)[i]
		for j := i + 1; j < n_moons; j++ {
			m2 := (*moons)[j]
			for k := 0; k < 3; k++ {
				if m1.position[k] > m2.position[k] {
					m1.velocity[k] -= 1
					m2.velocity[k] += 1
				} else if m1.position[k] < m2.position[k] {
					m1.velocity[k] += 1
					m2.velocity[k] -= 1
				}
			}
		}
	}
	for i := 0; i < n_moons; i++ {
		for j := 0; j < 3; j++ {
			(*moons)[i].position[j] += (*moons)[i].velocity[j]
		}
	}
}

func display(moons []*Moon) {
	for _, m := range moons {
		fmt.Printf("pos=%v, vel=%v\n", m.position, m.velocity)
	}
}

func main() {
	moons := load_moons("input.txt")

	// Part 1
	for i := 1; i <= 1000; i++ {
		apply_gravity(&moons)
	}
	energy := 0
	for _, m := range moons {
		energy += (abs(m.position[0]) + abs(m.position[1]) + abs(m.position[2])) *
			(abs(m.velocity[0]) + abs(m.velocity[1]) + abs(m.velocity[2]))
	}
	fmt.Printf("Part 1: %d\n", energy)

	// Part 2
	moons = load_moons("input.txt")
	history_x := make([](map[int]int), len(moons))
	history_y := make([](map[int]int), len(moons))
	history_z := make([](map[int]int), len(moons))
	for i, m := range moons {
		history_x[i] = map[int]int{}
		history_y[i] = map[int]int{}
		history_z[i] = map[int]int{}
		history_x[i][m.position[0]] = 0
		history_y[i][m.position[1]] = 0
		history_z[i][m.position[2]] = 0
	}
	iter := 0
	var x_period, y_period, z_period int
	found_x, found_y, found_z := false, false, false
loop:
	for {
		iter += 1
		apply_gravity(&moons)
		all_repeated_x, all_repeated_y, all_repeated_z := true, true, true
		for i := range moons {
			if before, ok := history_x[i][moons[i].position[0]]; ok && !found_x {
				fmt.Printf("Moon %d x-position repeated at iteration %d, delta: %d (seen before at %d)\n", i+1, iter, iter-before, before)
				history_x[i][moons[i].position[0]] = iter
			} else {
				all_repeated_x = false
			}
			if before, ok := history_y[i][moons[i].position[1]]; ok && !found_y {
				fmt.Printf("Moon %d y-position repeated at iteration %d, delta: %d (seen before at %d)\n", i+1, iter, iter-before, before)
				history_y[i][moons[i].position[1]] = iter
			} else {
				all_repeated_y = false
			}
			if before, ok := history_z[i][moons[i].position[2]]; ok && !found_z {
				fmt.Printf("Moon %d z-position repeated at iteration %d, delta: %d (seen before at %d)\n", i+1, iter, iter-before, before)
				history_z[i][moons[i].position[2]] = iter
			} else {
				all_repeated_z = false
			}
		}
		if all_repeated_x && x_period == 0 {
			x_period = iter + 1
			found_x = true
		}
		if all_repeated_y && y_period == 0 {
			y_period = iter + 1
			found_y = true
		}
		if all_repeated_z && z_period == 0 {
			z_period = iter + 1
			found_z = true
		}
		if found_x && found_y && found_z {
			break loop
		}
	}
	fmt.Printf("Part 2: %d\n", LCM(LCM(x_period, y_period), z_period))
}
