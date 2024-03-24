package main

import (
	"fmt"
	"os"
	"strings"
)

type grid [5][5]byte

func load_input(filename string) (eris grid) {
	content, _ := os.ReadFile(filename)
	lines := strings.Split(string(content), "\r\n")
	for i, l := range lines {
		for j, c := range []byte(l) {
			eris[i][j] = c
		}
	}
	return
}

func evolve(eris *grid) {
	adj := grid{}
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			if j > 0 && eris[i][j-1] == '#' {
				adj[i][j]++
			}
			if j < 4 && eris[i][j+1] == '#' {
				adj[i][j]++
			}
			if i > 0 && eris[i-1][j] == '#' {
				adj[i][j]++
			}
			if i < 4 && eris[i+1][j] == '#' {
				adj[i][j]++
			}
		}
	}
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			if eris[i][j] == '#' {
				if adj[i][j] != 1 {
					eris[i][j] = '.'
				}
			} else {
				if adj[i][j] == 1 || adj[i][j] == 2 {
					eris[i][j] = '#'
				}
			}
		}
	}
}

func display(eris *grid) {
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			fmt.Printf("%c", eris[i][j])
		}
		fmt.Println()
	}
	fmt.Println()
}

func biodiversity(eris *grid) (res uint32) {
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			if eris[i][j] == '#' {
				res |= 1 << (5*uint32(i) + uint32(j))
			}
		}
	}
	return
}

func main() {
	eris := load_input("input.txt")

	states := map[uint32]byte{}
	for {
		evolve(&eris)
		bd := biodiversity(&eris)
		if states[bd] > 0 {
			fmt.Printf("Part 1: %d\n", bd)
			break
		}
		states[bd]++
	}
}
