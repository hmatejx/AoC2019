package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func load_image(filename string, dimx int, dimy int) [][][]int {
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalln(err)
	}

	iptr := 0
	j := 0
	k := 0
	image := make([][][]int, len(content)/dimx/dimy)
	for l := 0; l < len(content)/dimx/dimy; l++ {
		image[l] = make([][]int, dimy)
		for r := 0; r < dimy; r++ {
			image[l][r] = make([]int, dimx)
		}
	}
	for iptr < len(content) {
		for i, c := range strings.Split(string(content)[iptr:iptr+dimx], "") {
			image[k][j][i], _ = strconv.Atoi(c)
		}
		iptr += dimx
		j += 1
		if j == dimy {
			k += 1
			j = 0
		}
	}
	return image
}

func Count(layer [][]int, digit int) int {
	res := 0
	for i := 0; i < len(layer); i++ {
		for j := 0; j < len(layer[0]); j++ {
			if layer[i][j] == digit {
				res += 1
			}
		}
	}
	return res
}

func Decode(image [][][]int) [][]int {
	dimx := len(image[0][0])
	dimy := len(image[0])
	decoded_image := make([][]int, dimy)
	for j := 0; j < dimy; j++ {
		decoded_image[j] = make([]int, dimx)
		for i := 0; i < dimx; i++ {
			var k int
			for k = 0; k < len(image); k++ {
				if image[k][j][i] != 2 {
					break
				}
			}
			decoded_image[j][i] = image[k][j][i]
		}
	}
	return decoded_image
}

func main() {
	dimx := 25
	dimy := 6
	image := load_image("input.txt", dimx, dimy)

	// Part 1
	min := 1<<63 - 1
	var layer [][]int
	for l := 0; l < len(image); l++ {
		res := Count(image[l], 0)
		if res < min {
			min = res
			layer = image[l]
		}
	}
	ones := Count(layer, 1)
	twos := Count(layer, 2)
	fmt.Printf("Part 1: %d\n", ones*twos)

	// Part 2
	fmt.Println("Part 2:")
	decoded := Decode(image)
	for j := 0; j < dimy; j++ {
		for i := 0; i < dimx; i++ {
			if decoded[j][i] == 1 {
				fmt.Printf("█")
			} else {
				fmt.Printf(" ")
			}
		}
		fmt.Println()
	}
}
