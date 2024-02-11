package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

type Graph map[string]map[string]int

func flood_fill(start [2]int, maze map[[2]int]byte) map[string]int {
	front := [][3]int{{start[0], start[1], 0}}
	visited := map[[2]int]int{}
	reachable := map[string]int{}
	var state [3]int
	for len(front) > 0 {
		state, front = front[0], front[1:]
		pos := [2]int{state[0], state[1]}
		visited[pos] = state[2]
		for d := 1; d <= 4; d++ {
			trial_pos := pos
			switch d {
			case 1:
				trial_pos[1]--
			case 2:
				trial_pos[1]++
			case 3:
				trial_pos[0]--
			case 4:
				trial_pos[0]++
			}
			if _, ok := visited[trial_pos]; !ok {
				dest := maze[trial_pos]
				if dest == 46 {
					front = append(front, [3]int{trial_pos[0], trial_pos[1], state[2] + 1})
				} else if (dest >= 64 && dest <= 90) || (dest >= 97 && dest <= 122) {
					reachable[string(dest)] = state[2] + 1
				}
			}
		}
	}
	return reachable
}

func is_key(node string) bool {
	return node[0] >= 97 && node[0] <= 122
}

func is_door(node string) bool {
	return node[0] >= 64 && node[0] <= 90
}

func load_input(filename string) Graph {
	content, _ := os.ReadFile(filename)
	lines := strings.Split(string(content), "\r\n")
	maze := map[[2]int]byte{}
	graph := Graph{}
	for y, line := range lines {
		for x := 0; x < len(line); x++ {
			maze[[2]int{x, y}] = line[x]
		}
		for k, v := range maze {
			s := string(v)
			if is_key(s) || is_door(s) {
				graph[s] = flood_fill(k, maze)
			}
		}
	}
	return graph
}

func copy_graph(graph Graph) Graph {
	out := Graph{}
	for k1, v1 := range graph {
		out[k1] = map[string]int{}
		for k2, v2 := range v1 {
			out[k1][k2] = v2
		}
	}
	return out
}

func remove_node(graph Graph, node string) {
	neighbors := []string{}
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
	node  string
	moves int
	path  []string
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

	node_state := func(path []string, current_node string) string {
		p := make([]string, len(path))
		copy(p, path)
		slices.Sort(p)
		pstr := ""
		for _, c := range p {
			pstr += c
		}
		return pstr + current_node
	}

	visited := map[string]bool{}
	front := StateHeap{}
	front.Push(&State{node: "@", moves: 0, graph: graph})
	best := 1<<63 - 1
	for len(front) > 0 {
		state := front.Pop().(*State)
		// check if we are at the last node
		if len(state.graph) <= 1 {
			if state.moves < best {
				best = state.moves
				fmt.Printf("New best: %d %v\n", best, state.path)
			}
			continue
		}
		visited[node_state(state.path, state.node)] = true
		// check which key we can pick up next
		for next, dist := range state.graph[state.node] {
			if is_key(next) {
				if visited[node_state(state.path, next)] {
					continue
				}
				// update the graph, remove the current node and the corresponding door if it exists
				new_graph := copy_graph(state.graph)
				remove_node(new_graph, state.node)
				key_node := strings.ToUpper(next)
				if _, ok := new_graph[key_node]; ok {
					remove_node(new_graph, key_node)
				}
				new_state := State{
					node:  next,
					moves: state.moves + dist,
					graph: new_graph,
				}
				lp := len(state.path)
				new_state.path = make([]string, lp+1)
				copy(new_state.path, state.path)
				new_state.path[lp] = next
				front.Push(&new_state)
			}
		}
	}
	return best
}

func main() {
	graph := load_input("test.txt")

	// Part 1
	best := collect_keys(graph)
	fmt.Printf("%v\n", best)
}
