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
	possible := 0
	for _, design := range designs {
		if isPossible(towels, design) {
			possible++
		}
	}
	return possible
}

type DumbQueue []string

func (q *DumbQueue) push(item string) { *q = append(*q, item) }

func (q *DumbQueue) pop() string {
	item := (*q)[len(*q)-1]
	*q = (*q)[0 : len(*q)-1]
	return item
}

func isPossible(towels []string, design string) bool {
	workQueue := DumbQueue{}
	workQueue.push("")

	for len(workQueue) > 0 {
		current := workQueue.pop()

		if current == design {
			return true
		}

		for _, towel := range towels {
			if strings.HasPrefix(design, current+towel) {
				workQueue.push(current + towel)
			}
		}
	}

	return false
}
