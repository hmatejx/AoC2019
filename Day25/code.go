package main

import (
	"AoC2019/Day25/code/intcode"
	"fmt"
	"strings"
)

func main() {
	droid := intcode.NewIntcodeCPU()
	droid.Load("input.txt")

	// interactive exploration of the map revealed the following path
	input := []string{"north\n", "take sand\n", "north\n", "take space heater\n",
		"east\n", "take semiconductor\n", "west\n", "south\n", "south\n",
		"east\n", "take ornament\n", "south\n", "take festive hat\n", "east\n",
		"take asterisk\n", "south\n", "west\n", "take food ration\n", "east\n",
		"east\n", "take cake\n", "west\n", "north\n", "west\n", "north\n",
		"west\n", "west\n", "north\n", "north\n", "inv"}
	out, state := droid.Run(intcode.From_ASCII(strings.Join(input, "")))
	fmt.Printf("%v %v %v\n", intcode.From_ASCII(strings.Join(input, "")), intcode.To_ASCII(out), state)

	// try all combinations of valid items (using gray code for faster picking and dropping)
	items := []string{"asterisk", "ornament", "cake", "space heater", "festive hat", "semiconductor", "food ration", "sand"}
	var i uint8 = 0
	for {
		gray := i ^ (i >> 1)
		input := []int{}
		for b := 0; b < 8; b++ {
			if gray&(1<<b) == 0 {
				input = append(input, intcode.From_ASCII(strings.Join([]string{"drop", items[b]}, " "))...)
			}
		}
		input = append(input, intcode.From_ASCII("inv\nwest")...)
		out, state := droid.Copy().Run(input)
		i++
		if i == 0 || state == 99 {
			fmt.Printf("%v\n", intcode.To_ASCII(out))
			break
		}
	}
}
