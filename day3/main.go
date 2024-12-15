package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func main() {
	var fileName string

	if len(os.Args) == 2 {
		fileName = os.Args[1]
	} else {
		fileName = "input.txt"
	}

	allMatches := readInput(fileName)

	// part 1 && part 2

	runningTotal := 0
	switchedRunningTotal := 0
	enabled := true

	for _, match := range allMatches {
		if match[0] == "do()" {
			enabled = true
		} else if match[0] == "don't()" {
			enabled = false
		} else {
			a, _ := strconv.Atoi(match[1])
			b, _ := strconv.Atoi(match[2])

			runningTotal += a * b

			if enabled {
				switchedRunningTotal += a * b
			}
		}
	}

	fmt.Printf("Running total: %d\n", runningTotal)
	fmt.Printf("Switched running total: %d\n", switchedRunningTotal)
}

func readInput(filePath string) [][]string {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	allMatches := [][]string{}
	re := regexp.MustCompile(`mul\((\d{1,3}),(\d{1,3})\)|do\(\)|don't\(\)`)

	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindAllStringSubmatch(line, -1)

		allMatches = append(allMatches, matches...)
	}

	return allMatches
}
