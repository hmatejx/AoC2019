package main

import (
	"AoC2019/Day13/code/intcode"
	"fmt"
)

func display(arcade map[[2]int]int) {
	dict := map[int]string{0: " ", 1: "█", 2: "░", 3: "▔", 4: "●"}
	// works only for fixed size 42x24 arcade size
	for y := 0; y <= 23; y++ {
		for x := 0; x <= 41; x++ {
			fmt.Printf(dict[arcade[[2]int{x, y}]])
		}
		fmt.Println()
	}
}

// returns the position of the ball, paddle, the number of blocks, and score
func update_game_state(output []int, arcade *map[[2]int]int, oldscore int) (int, int, int, int) {
	var ball, paddle, blocks int
	score := oldscore
	for i := 0; i < len(output)-2; i += 3 {
		(*arcade)[[2]int{output[i], output[i+1]}] = output[i+2]
		if output[i+2] == 4 {
			ball = output[i]
		} else if output[i+2] == 3 {
			paddle = output[i]
		}
		if output[i] == -1 {
			score = output[i+2]
		}
	}
	for _, v := range *arcade {
		if v == 2 {
			blocks++
		}
	}
	return ball, paddle, blocks, score
}

func main() {
	cabinet := intcode.NewIntcodeCPU()
	cabinet.Load("input.txt")

	// Part 1
	output, _ := cabinet.Run([]int{})
	arcade := map[[2]int]int{}
	ball, paddle, blocks, _ := update_game_state(output, &arcade, 0)
	fmt.Printf("Part 1: %d\n", blocks)

	// Part 2
	cabinet.Reset()
	cabinet.Poke(0, 2) // insert coin
	var score, state int
	for state != 99 {
		joystick := []int{0}
		if ball > paddle {
			joystick[0] = 1
		} else if ball < paddle {
			joystick[0] = -1
		}
		output, state = cabinet.Run(joystick)
		ball, paddle, _, score = update_game_state(output, &arcade, score)
		fmt.Print("\033[H")
		display(arcade)
	}
	fmt.Printf("Part 2: %d\n", score)
}
