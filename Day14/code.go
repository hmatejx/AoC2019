package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type component struct {
	name   string
	amount int
}

type reaction struct {
	target      component
	ingredients []component
}

func load_reactions(filename string) map[string]reaction {
	get_component := func(str string) component {
		t := strings.Split(str, " ")
		a := t[1]
		b, _ := strconv.Atoi(t[0])
		return component{name: a, amount: b}
	}
	content, _ := os.ReadFile(filename)
	reactions := map[string]reaction{}
	lines := strings.Split(string(content), "\r\n")
	for _, line := range lines {
		lr := strings.Split(line, " => ")
		target := get_component(lr[1])
		ingredient_list := strings.Split(lr[0], ", ")
		ingredients := []component{}
		for _, ing := range ingredient_list {
			ingredients = append(ingredients, get_component(ing))
		}
		reactions[target.name] = reaction{target: target, ingredients: ingredients}
	}
	return reactions
}

func cascade(fuel int, reactions *map[string]reaction) int {
	bom, surplus := map[string]int{"FUEL": fuel}, map[string]int{}
	for {
		for target := range bom {
			if target == "ORE" {
				continue
			}
			// first spend the surplus
			if surplus[target] >= bom[target] {
				surplus[target] -= bom[target]
				delete(bom, target)
				continue
			}
			bom[target] -= surplus[target]
			surplus[target] = 0
			// then process the reactions
			reaction := (*reactions)[target]
			units := bom[target] / reaction.target.amount
			if bom[target]%reaction.target.amount != 0 {
				units += 1
				surplus[target] = units*reaction.target.amount - bom[target]
			}
			for _, ing := range reaction.ingredients {
				bom[ing.name] += ing.amount * units
			}
			delete(bom, target)
		}
		if len(bom) == 1 {
			break
		}
	}
	return bom["ORE"]
}

func main() {
	reactions := load_reactions("input.txt")

	// Part 1
	res := cascade(1, &reactions)
	fmt.Printf("Part 1: %d\n", res)

	// Part 2
	lo, hi, mid := 1, 10000000, 0
	for {
		mid = (lo + hi) / 2
		res = cascade(mid, &reactions)
		if hi-lo <= 2 {
			break
		}
		if res <= 1000000000000 {
			lo = mid
		} else {
			hi = mid
		}
	}
	fmt.Printf("Part 2: %d", lo)
}
