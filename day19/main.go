package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	var fileName string

	if len(os.Args) == 2 {
		fileName = os.Args[1]
	} else {
		fileName = "input.txt"
	}

	towels, designs := readInput(fileName)

	// Part 1
	possible := Part1(towels, designs)
	fmt.Printf("Part 1: %d designs possible\n", possible)

	// Part 2
	combos := Part2(towels, designs)
	fmt.Printf("Part 2: %d possible combinations\n", combos)
}

func readInput(filePath string) ([]string, []string) {
	file, _ := os.Open(filePath)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var towels []string
	var designs []string

	for scanner.Scan() {
		line := scanner.Text()

		if towels == nil {
			towels = strings.Split(line, ", ")
			continue
		}

		if len(line) == 0 {
			continue
		}

		designs = append(designs, line)
	}

	return towels, designs
}

func Part1(towels []string, designs []string) int {
	memo := map[string]int{}
	possible := 0

	for _, design := range designs {
		if countCombos(design, towels, memo) > 0 {
			possible++
		}
	}

	return possible
}

func Part2(towels []string, designs []string) int {
	memo := map[string]int{}
	totalCombos := 0

	for _, design := range designs {
		totalCombos += countCombos(design, towels, memo)
	}

	return totalCombos
}

func countCombos(design string, pieces []string, memo map[string]int) int {
	if len(design) == 0 {
		return 1
	}

	if cacheHit, ok := memo[design]; ok {
		return cacheHit
	}

	foundCombos := 0

	for _, piece := range pieces {
		if len(piece) <= len(design) && strings.HasPrefix(design, piece) {
			foundCombos += countCombos(design[len(piece):], pieces, memo)
		}
	}

	memo[design] = foundCombos
	return foundCombos
}
