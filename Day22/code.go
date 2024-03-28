package main

import (
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
)

func load_input(filename string) []string {
	content, _ := os.ReadFile(filename)
	return strings.Split(string(content), "\r\n")
}

func LCF(x, a, b, m int64) int64 {
	return (a*x + b) % m
}

func shuffle(pos int64, ncards int64, instructions []string) int64 {
	var a, b int64
	for _, inst := range instructions {
		if inst[0:3] == "cut" {
			d, _ := strconv.Atoi(inst[4:])
			a, b = 1, ncards-int64(d)
		} else if inst[0:9] == "deal with" {
			d, _ := strconv.Atoi(inst[20:])
			a, b = int64(d), 0
		} else {
			a, b = -1, -1
		}
		pos = LCF(pos, a, b, ncards)
	}
	return pos
}

func combine_two_LCF(a1, b1, a2, b2, m int64) [2]int64 {
	A, B, A1, M := big.NewInt(0), big.NewInt(0), big.NewInt(a1), big.NewInt(m)
	A.Mul(A1, big.NewInt(a2))
	B.Add(B.Mul(A1, big.NewInt(b2)), big.NewInt(b1))
	return [2]int64{A.Mod(A, M).Int64(), B.Mod(B, M).Int64()}
}

func LCF_to_power(power, a, b, m int64) [2]int64 {
	if power == 0 {
		return [2]int64{1, 0}
	}
	if power%2 == 1 {
		res := LCF_to_power(power-1, a, b, m)
		return combine_two_LCF(res[0], res[1], a, b, m)
	}
	res := LCF_to_power(power/2, a, b, m)
	return combine_two_LCF(res[0], res[1], res[0], res[1], m)
}

func main() {
	instructions := load_input("input.txt")

	// Part 1
	pos := shuffle(int64(2019), int64(10007), instructions)
	fmt.Printf("Part 1: %d\n", pos)

	// Part 2
	ncards := int64(119315717514047)
	times := int64(101741582076661)
	// first, lets find the LCF parameters a & b for a single round
	pos0 := shuffle(int64(0), ncards, instructions)
	pos1 := shuffle(int64(1), ncards, instructions)
	a := pos1 - pos0
	b := pos0
	// calculate the parameters of n rounds
	res := LCF_to_power(times, a, b, ncards)
	a, b = res[0], res[1]
	// invert the LCF
	m, invA, invB, target := big.NewInt(ncards), big.NewInt(0), big.NewInt(0), big.NewInt(2020)
	invA.ModInverse(big.NewInt(a), m)
	invB.Mod(invB.Mul(big.NewInt(-b), invA), m)
	target.Mul(target, invA)
	target.Add(target, invB)
	target.Mod(target, m)
	fmt.Printf("%v (started at %v)\n", target, 2019)
}
