package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strings"
)

func load_asteroids(filename string) ([][2]int, int, int) {
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalln(err)
	}
	asteroids := [][2]int{}
	lines := strings.Split(string(content), "\r\n")
	for i, line := range lines {
		for j, c := range line {
			if c == '#' {
				asteroids = append(asteroids, [2]int{j, i})
			}
		}
	}
	return asteroids, len(lines[0]), len(lines)
}

func distance2(a1, a2 [2]int) int {
	delta_x := a1[0] - a2[0]
	delta_y := a1[1] - a2[1]
	return delta_x*delta_x + delta_y*delta_y
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

func all_rays(a [2]int, asteroids [][2]int) map[[2]int][]int {
	// calculate all rays defined by asteroid a and other asteroids
	rays := map[[2]int][]int{}
	for i := 0; i < len(asteroids); i++ {
		if asteroids[i] == a {
			continue
		}
		dx := asteroids[i][0] - a[0]
		dy := asteroids[i][1] - a[1]
		gcd := abs(GCD(dy, dx))
		key := [2]int{dy / gcd, dx / gcd}
		if subset, ok := rays[key]; ok {
			rays[key] = append(subset, i)
		} else {
			rays[key] = []int{i}
		}
	}
	return rays
}

func main() {
	asteroids, _, _ := load_asteroids("input.txt")

	// Part 1
	max := 0
	var rays_map map[[2]int][]int
	var station [2]int
	for _, a := range asteroids {
		rays_map_a := all_rays(a, asteroids)
		if len(rays_map_a) > max {
			max = len(rays_map_a)
			station = a
			rays_map = rays_map_a
		}
	}
	fmt.Printf("Part 1: %d\n", max)

	// Part 2
	// sort rays of asteroids by angle and asteroids within a ray by distance to station
	type Asteroid struct {
		dist2 int
		coord [2]int
	}
	type LaserRayTargets struct {
		angle     float64
		asteroids []Asteroid
	}
	destroy := func(lt *LaserRayTargets) (int, [2]int) {
		if len(lt.asteroids) == 0 {
			return 0, [2]int{}
		}
		var coord [2]int
		coord, lt.asteroids = lt.asteroids[0].coord, lt.asteroids[1:]
		return 1, coord
	}
	laser_targets := []LaserRayTargets{}
	for k, v := range rays_map {
		phi := (-math.Atan2(float64(k[1]), float64(k[0])) + math.Pi) * 180 / math.Pi
		if phi < 0 {
			phi += 360
		}
		a := []Asteroid{}
		for _, a_idx := range v {
			a = append(a, Asteroid{dist2: distance2(station, asteroids[a_idx]), coord: asteroids[a_idx]})
		}
		sort.Slice(a, func(i, j int) bool {
			return a[i].dist2 < a[j].dist2
		})
		laser_targets = append(laser_targets, LaserRayTargets{angle: phi, asteroids: a})
	}
	sort.Slice(laser_targets, func(i, j int) bool {
		return laser_targets[i].angle < laser_targets[j].angle
	})
	// now hit them one by one while rotating the laser beam clockwise
	i := 0
	to_destroy := 200
	destroyed := 0
	var coord [2]int
	for {
		var d int
		d, coord = destroy(&laser_targets[i])
		destroyed += d
		if destroyed == to_destroy {
			break
		}
		// rotate to the next ray angle
		i = (i + 1) % len(laser_targets)
	}
	fmt.Printf("Part 2: %d\n", 100*coord[0]+coord[1])
}
