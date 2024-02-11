package main

import (
	"AoC2019/Day15/code/intcode"
	"fmt"
)

type Robot struct {
	icpu           *intcode.IntcodeCPU
	position       [2]int
	trial_position [2]int
}

type State struct {
	steps  int
	status int
	robot  *Robot
}

func new_robot(filename string) *Robot {
	new_robot := &Robot{icpu: intcode.NewIntcodeCPU()}
	new_robot.icpu.Load(filename)
	return new_robot
}

func move(position *[2]int, direction int) {
	switch direction {
	case 1:
		position[1]--
	case 2:
		position[1]++
	case 3:
		position[0]--
	case 4:
		position[0]++
	}
}

func (r *Robot) run(input int) int {
	output, _ := r.icpu.Run([]int{input})
	r.trial_position = r.position
	move(&r.trial_position, input)
	if output[0] > 0 {
		r.position = r.trial_position
	}
	return output[0]
}

func (r *Robot) copy() *Robot {
	return &Robot{icpu: r.icpu.Copy(), position: r.position}
}

func explore_maze(robot *Robot) (map[[2]int]int, int, [2]int) {
	maze := map[[2]int]int{}
	visited := map[[2]int]int{}
	front := []State{{steps: 0, status: 3, robot: robot}}
	var min_steps int
	var oxigen [2]int
	for len(front) > 0 {
		cur_state := front[0]
		front = front[1:]
		visited[cur_state.robot.position] = cur_state.steps + 1
		maze[cur_state.robot.position] = cur_state.status
		// record the number of steps we first reach the target
		if cur_state.status == 2 && min_steps == 0 {
			min_steps = cur_state.steps
			oxigen = cur_state.robot.position
		}
		// try to move in every direction
		for d := 1; d <= 4; d++ {
			// check if we've already been here
			trial_position := cur_state.robot.position
			move(&trial_position, d)
			if visited[trial_position] != 0 {
				continue
			}
			// let the try to move
			status := cur_state.robot.run(d)
			if status > 0 {
				// add a copy of the robot to explore further to the front
				front = append(front, State{steps: cur_state.steps + 1, status: status, robot: cur_state.robot.copy()})
				switch d { // move the robot back to where it was
				case 1, 2:
					cur_state.robot.run(3 - d)
				case 3, 4:
					cur_state.robot.run(7 - d)
				}
			} else {
				maze[trial_position] = 0
			}
		}
	}
	return maze, min_steps, oxigen
}

func flood_fill(start [2]int, maze map[[2]int]int) int {
	front := [][3]int{{start[0], start[1], 0}}
	visited := map[[2]int]int{}
	max_steps := 0
	var state [3]int
	for len(front) > 0 {
		state, front = front[0], front[1:]
		pos := [2]int{state[0], state[1]}
		visited[pos] = state[2]
		if state[2] > max_steps {
			max_steps = state[2]
		}
		for d := 1; d <= 4; d++ {
			trial_pos := pos
			move(&trial_pos, d)
			if _, ok := visited[trial_pos]; !ok && maze[trial_pos] > 0 {
				front = append(front, [3]int{trial_pos[0], trial_pos[1], state[2] + 1})
			}
		}
	}
	return max_steps
}

func display(maze map[[2]int]int) {
	dict := map[int]string{3: "*", 0: "█", 1: " ", 2: "O"}
	x0, x1, y0, y1 := 1<<63-1, -1<<63, 1<<63-1, -1<<63
	for k := range maze {
		if k[0] < x0 {
			x0 = k[0]
		}
		if k[0] > x1 {
			x1 = k[0]
		}
		if k[1] < y0 {
			y0 = k[1]
		}
		if k[1] > y1 {
			y1 = k[1]
		}
	}
	for y := y0; y <= y1; y++ {
		for x := x0; x <= x1; x++ {
			fmt.Printf(dict[maze[[2]int{x, y}]])
		}
		fmt.Println()
	}
}

func main() {
	robot := new_robot("input.txt")

	// Part 1
	maze, steps, oxigen := explore_maze(robot)
	display(maze)
	fmt.Printf("Part 1: %d\n", steps)

	// Part 2
	time := flood_fill(oxigen, maze)
	fmt.Printf("Part 2: %d\n", time)
}
