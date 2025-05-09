package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	var fileName string

	if len(os.Args) == 2 {
		fileName = os.Args[1]
	} else {
		fileName = "input.txt"
	}

	locks, keys := readInput(fileName)

	fmt.Printf("There are %d lock combos that fit.", Part1(locks, keys))
}

func readInput(filePath string) ([][]int, [][]int) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening %s", filePath)
		return nil, nil
	}
	defer file.Close()

	locks := [][]int{}
	keys := [][]int{}

	scanner := bufio.NewScanner(file)

	pintype := ""
	var pins []int
	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			switch pintype {
			case "lock":
				locks = append(locks, pins)
			case "key":
				keys = append(keys, pins)
			}

			pintype = ""
			continue
		}

		if pintype == "" && line == "....." {
			pintype = "key"
			pins = []int{-1, -1, -1, -1, -1}
			continue
		}

		if pintype == "" && line == "#####" {
			pintype = "lock"
			pins = []int{0, 0, 0, 0, 0}
			continue
		}

		for i, c := range line {
			if c == '#' {
				pins[i]++
			}
		}
	}

	// Commit one more time since the input doesn't end in a blank line
	switch pintype {
	case "lock":
		locks = append(locks, pins)
	case "key":
		keys = append(keys, pins)
	}

	return locks, keys
}

func Part1(locks [][]int, keys [][]int) int {
	fits := 0

	for _, lock := range locks {
		for _, key := range keys {
			if lock[0]+key[0] < 6 &&
				lock[1]+key[1] < 6 &&
				lock[2]+key[2] < 6 &&
				lock[3]+key[3] < 6 &&
				lock[4]+key[4] < 6 {
				fits++
			}
		}
	}

	return fits
}
