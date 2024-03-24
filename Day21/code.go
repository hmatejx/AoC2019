package main

import (
	"AoC2019/Day21/code/intcode"
	"fmt"
)

func main() {
	robot := intcode.NewIntcodeCPU()
	robot.Load("input.txt")

	// part 1
	prog :=
		"NOT C J\n" +
			"AND D J\n" + // IF C is a hole AND D is not, jump
			"NOT A T\n" +
			"OR T J\n" + // OR if A is a hole, jump
			"WALK\n"
	res, _ := robot.Run(intcode.From_ASCII(prog))
	fmt.Printf("Part 1: %d\n", res[len(res)-1])

	// part 2
	prog =
		"NOT C J\n" +
			"AND D J\n" +
			"AND H J\n" + // IF C is a hole AND D is not and H is not, jump
			"NOT B T\n" +
			"AND D T\n" +
			"OR T J\n" + // OR IF B is a hole AND D is not, jump
			"NOT A T\n" +
			"OR T J\n" + // OR IF A is a hole, jump
			"RUN\n"
	robot.Reset()
	res, _ = robot.Run(intcode.From_ASCII(prog))
	fmt.Printf("Part 2: %d\n", res[len(res)-1])
}
