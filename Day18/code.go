package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

type Graph map[byte]map[byte]int

func flood_fill(startx, starty int, maze [][]byte) map[byte]int {
	front := [][3]int{{startx, starty, 0}}
	visited := map[[2]int]int{}
	reachable := map[byte]int{}
	var state [3]int
	for len(front) > 0 {
		state, front = front[0], front[1:]
		x, y := state[0], state[1]
		visited[[2]int{x, y}] = state[2]
		for d := 1; d <= 4; d++ {
			newx, newy := x, y
			switch d {
			case 1:
				newy--
			case 2:
				newy++
			case 3:
				newx--
			case 4:
				newx++
			}
			if _, ok := visited[[2]int{newx, newy}]; !ok {
				dest := maze[newy][newx]
				if dest == '.' {
					front = append(front, [3]int{newx, newy, state[2] + 1})
				} else if (dest >= '@' && dest <= 'Z') || (dest >= 'a' && dest <= 'z') {
					reachable[dest] = state[2] + 1
				}
			}
		}
	}
	return reachable
}

func load_input(filename string) Graph {
	content, _ := os.ReadFile(filename)
	lines := strings.Split(string(content), "\r\n")
	maze := make([][]byte, len(lines))
	for y := range maze {
		maze[y] = make([]byte, len(lines[0]))
	}
	graph := Graph{}
	for y, line := range lines {
		for x := 0; x < len(line); x++ {
			maze[y][x] = line[x]
		}
	}
	for y := range maze {
		for x, v := range maze[y] {
			//s := string(v)
			if (v >= '@' && v <= 'Z') || (v >= 'a' && v <= 'z') {
				graph[v] = flood_fill(x, y, maze)
			}
		}
	}
	return graph
}

func copy_graph(graph Graph) Graph {
	out := Graph{}
	for k1, v1 := range graph {
		out[k1] = map[byte]int{}
		for k2, v2 := range v1 {
			out[k1][k2] = v2
		}
	}
	return out
}

func remove_node(graph Graph, node byte) {
	neighbors := []byte{}
	for k := range graph[node] {
		neighbors = append(neighbors, k)
	}
	nn := len(neighbors)
	// connect neighbors together
	for i := 0; i < nn; i++ {
		n1 := neighbors[i]
		a := graph[node][n1]
		for j := i + 1; j < nn; j++ {
			n2 := neighbors[j]
			b, c := graph[node][n2], graph[n1][n2]
			d := a + b
			if c == 0 || d < c {
				graph[n1][n2] = d
				graph[n2][n1] = d
			}
		}
	}
	// remove the node
	for k := range graph[node] {
		delete(graph[k], node)
	}
	delete(graph, node)
}

type State struct {
	node  byte
	moves int
	keys  string
	graph Graph
}

type StateHeap []*State

func (h StateHeap) Len() int {
	return len(h)
}

func (h StateHeap) Less(i, j int) bool {
	return h[i].moves < h[j].moves
}

func (h StateHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *StateHeap) Push(x interface{}) {
	item := x.(*State)
	*h = append(*h, item)
}

func (h *StateHeap) Pop() any {
	old := *h
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*h = old[0 : n-1]
	return item
}

func collect_keys(graph Graph) int {
	visited := map[string]int{}
	front := StateHeap{}
	front.Push(&State{node: '@', moves: 0, graph: graph})
	best := 1<<63 - 1
	for len(front) > 0 {
		state := front.Pop().(*State)
		// check if we are at the last node
		if len(state.graph) <= 1 {
			if state.moves < best {
				best = state.moves
				fmt.Printf("New best: %d\n", best)
			}
			continue
		}
		state_str := state.keys + string(state.node)
		// we've been here before faster
		seen := visited[state_str]
		if seen > 0 && seen <= state.moves {
			continue
		}
		visited[state_str] = state.moves
		// check which key we can pick up next
		for next, dist := range state.graph[state.node] {
			if next >= 'a' && next <= 'z' {
				// update the graph, remove the current node and the corresponding door if it exists
				new_graph := copy_graph(state.graph)
				remove_node(new_graph, state.node)
				door_node := strings.ToUpper(string(next))[0]
				if _, ok := new_graph[door_node]; ok {
					remove_node(new_graph, door_node)
				}
				// push new state to heap
				new_keys := append([]byte(state.keys), next)
				slices.Sort(new_keys)
				new_state := State{
					node:  next,
					moves: state.moves + dist,
					keys:  string(new_keys),
					graph: new_graph,
				}
				front.Push(&new_state)
			}
		}
	}
	return best
}

func main() {
	graph := load_input("input.txt")

	// Part 1
	best := collect_keys(graph)
	fmt.Printf("Part 1: %d\n", best)
}
