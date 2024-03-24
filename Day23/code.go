package main

import (
	"AoC2019/Day23/code/intcode"
	"fmt"
)

func simulate_network(include_nat bool) int {
	// boot up the computers
	ncomp := 50
	nodes := make([]*intcode.IntcodeCPU, ncomp)
	for addr := 0; addr < ncomp; addr++ {
		nodes[addr] = intcode.NewIntcodeCPU()
		nodes[addr].Load("input.txt")
		nodes[addr].Run([]int{addr})
	}
	// send packets around until something happens
	queues := make([][]int, ncomp)
	nat_packet := []int{}
	nat_y := map[int]bool{}
	for {
	not_idle:
		for addr := 0; addr < ncomp; addr++ {
			// consume first packet from the queue
			input := []int{-1}
			if len(queues[addr]) > 0 {
				input = queues[addr][0:2]
				queues[addr] = queues[addr][2:]
			}
			out, _ := nodes[addr].Run(input)
			// push output packets to the respective queues
			for i := 0; i < len(out); i += 3 {
				// check for packet to address 255 (NAT)
				if out[i] == 255 {
					if !include_nat {
						return out[i+2]
					}
					nat_packet = out[i+1 : i+3]
				} else {
					queues[out[i]] = append(queues[out[i]], out[i+1:i+3]...)
				}
			}
		}
		// check if all computers are idle
		if include_nat {
			for addr := 0; addr < ncomp; addr++ {
				if len(queues[addr]) > 0 {
					goto not_idle
				}
			}
			// check if Y value of the NAT packet repeats
			if nat_y[nat_packet[1]] {
				return nat_packet[1]
			}
			// remember the NAT packet y value and send the NAT packet to node 0
			nat_y[nat_packet[1]] = true
			queues[0] = append(queues[0], nat_packet...)
		}
	}
}

func main() {
	fmt.Printf("Part 1: %d\n", simulate_network(false))
	fmt.Printf("Part 2: %d\n", simulate_network(true))
}
