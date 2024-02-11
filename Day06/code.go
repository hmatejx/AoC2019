package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func read_input(filename string) (map[string]string, map[string][]string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalln(err)
	}
	lines := strings.Split(string(content), "\r\n")

	parent := map[string]string{}
	children := map[string][]string{}
	for _, line := range lines {
		nodes := strings.Split(line, ")")
		p, c := nodes[0], nodes[1]
		// assign child to parent
		if ch, ok := children[p]; ok {
			children[p] = append(ch, c)
		} else {
			children[p] = []string{c}
		}
		// assign parent to child
		parent[c] = p
		// add child node to list of parents
		if _, ok := children[c]; !ok {
			children[c] = []string{}
		}
		parent["COM"] = ""
	}
	return parent, children
}

func get_endnodes(children map[string][]string) []string {
	// get end nodes
	end_nodes := []string{}
	for k, v := range children {
		if len(v) == 0 {
			end_nodes = append(end_nodes, k)
		}
	}
	return end_nodes
}

func get_indirect_orbits(end_nodes []string, parent map[string]string) map[string](map[string]int) {
	// generate all indirect relationships
	indirect := map[string](map[string]int){}
	for _, node := range end_nodes {
		for {
			p := parent[node]
			sep := 1
			// check if we are at the root or if we already walked along this chain
			if p == "COM" || indirect[parent[p]][node] > 0 {
				break
			}
			for {
				p = parent[p]
				sep++
				if p == "" {
					break
				}
				if _, ok := indirect[p]; ok {
					indirect[p][node] = sep
				} else {
					indirect[p] = map[string]int{node: sep}
				}
			}
			node = parent[node]
		}
	}
	return indirect
}

func main() {
	direct, children := read_input("test.txt")
	indirect := get_indirect_orbits(get_endnodes(children), direct)

	// Part 1
	ndirect, nindirect := 0, 0
	for _, p := range direct {
		if p != "" {
			ndirect++
		}
	}
	for _, node := range indirect {
		nindirect += len(node)
	}
	fmt.Printf("Part 1: %d\n", ndirect+nindirect)

	// Part 2
	min_transfers := 1<<63 - 1
	for _, values := range indirect {
		d_san, san := values["SAN"]
		d_you, you := values["YOU"]
		// check if both you and santa are indirectly connected to this node
		if san && you {
			transfers := d_san + d_you - 2
			if transfers < min_transfers {
				min_transfers = transfers
			}
		}
	}
	fmt.Printf("Part 2: %d\n", min_transfers)
}
