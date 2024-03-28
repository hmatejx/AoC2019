package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

type grid [5][5]byte

var grid_storage, grid_map, empty_grid = []grid{}, map[int]*grid{}, grid{}

func load_input(filename string) (eris grid) {
	content, _ := os.ReadFile(filename)
	lines := strings.Split(string(content), "\r\n")
	for i, l := range lines {
		for j, c := range []byte(l) {
			if c == '#' {
				eris[i][j] = 1
			}
		}
	}
	return
}

func map_keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func get_level(level int) *grid {
	if _, ok := grid_map[level]; ok {
		return grid_map[level]
	} else {
		return &empty_grid
	}
}

func count_neighbors(level, i, j int, recursive bool) (neigh byte) {
	// levels: above > current > below
	above := get_level(level - 1)
	current := get_level(level)
	below := get_level(level + 1)
	// always check these
	if j > 0 && current[i][j-1] == 1 {
		neigh++
	}
	if j < 4 && current[i][j+1] == 1 {
		neigh++
	}
	if i > 0 && current[i-1][j] == 1 {
		neigh++
	}
	if i < 4 && current[i+1][j] == 1 {
		neigh++
	}
	// handle the additional neighbors from lower or higher level grids
	if recursive {
		// check higher level for corner and edge tiles
		if i == 0 && above[1][2] == 1 {
			neigh++
		} else if i == 4 && above[3][2] == 1 {
			neigh++
		}
		if j == 0 && above[2][1] == 1 {
			neigh++
		} else if j == 4 && above[2][3] == 1 {
			neigh++
		}
		// check lower levels for the special center cross tiles
		idx := 5*i + j
		for k := 0; k < 5; k++ {
			switch idx {
			case 7:
				neigh += below[0][k]
			case 11:
				neigh += below[k][0]
			case 13:
				neigh += below[k][4]
			case 17:
				neigh += below[4][k]
			}
		}
	}
	return
}

func evolve(recursive bool) {
	// get list of levels
	levels := map_keys(grid_map)
	minlev, maxlev := slices.Min(levels)-1, slices.Max(levels)+1
	// when recursive add an upper and lower level
	if recursive {
		levels = append(levels, minlev, maxlev)
	}
	// count neighbors
	neigh := make([]grid, len(levels))
	for k, l := range levels {
		n := &neigh[k]
		for i := 0; i < 5; i++ {
			for j := 0; j < 5; j++ {
				if i == 2 && j == 2 && recursive {
					continue
				}
				n[i][j] = count_neighbors(l, i, j, recursive)
			}
		}
	}
	// update grids
	for k, l := range levels {
		g := get_level(l)
		for i := 0; i < 5; i++ {
			for j := 0; j < 5; j++ {
				if i == 2 && j == 2 && recursive {
					continue
				}
				if g[i][j] == 1 {
					if neigh[k][i][j] != 1 {
						g[i][j] = 0
					}
				} else {
					if neigh[k][i][j] == 1 || neigh[k][i][j] == 2 {
						// add a new grid to the storage and map if necessary
						if g == &empty_grid && recursive {
							grid_storage = append(grid_storage, empty_grid)
							g = &grid_storage[len(grid_storage)-1]
							grid_map[l] = g
						}
						g[i][j] = 1
					}
				}
			}
		}
	}
}

func display(level int, recursive bool) {
	g := get_level(level)
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			if recursive && i == 2 && j == 2 {
				fmt.Print("?")
			} else if g[i][j] == 1 {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func biodiversity(g *grid) (res uint32) {
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			if g[i][j] > 0 {
				res |= 1 << (5*uint32(i) + uint32(j))
			}
		}
	}
	return
}

func setup(initial grid) {
	// (re-)initialize the grid storage and mapping
	grid_storage = nil
	grid_storage = append(grid_storage, initial)
	clear(grid_map)
	grid_map[0] = &grid_storage[0]
}

func main() {
	eris := load_input("input.txt")

	// part 1
	setup(eris)
	states := map[uint32]byte{}
	for {
		evolve(false)
		bd := biodiversity(grid_map[0])
		if states[bd] > 0 {
			fmt.Printf("Part 1: %d\n", bd)
			break
		}
		states[bd]++
	}

	// part 2
	setup(eris)
	for i := 0; i < 200; i++ {
		evolve(true)
	}
	bugs := 0
	for l := -100; l <= 100; l++ {
		for i := 0; i < 5; i++ {
			for j := 0; j < 5; j++ {
				bugs += int(grid_map[l][i][j])
			}
		}
	}
	fmt.Printf("Part 2: %d\n", bugs)
}
