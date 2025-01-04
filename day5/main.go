package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

func main() {
	var fileName string

	if len(os.Args) == 2 {
		fileName = os.Args[1]
	} else {
		fileName = "input.txt"
	}

	rules, updates := readInput(fileName)

	// part 1
	pageSum := PartOne(rules, updates)
	fmt.Printf("The sum of the middle page numbers is %d\n", pageSum)

	// part 2
	correctedPageSum := PartTwo(rules, updates)
	fmt.Printf("The sum of the corrected middle page numbers is %d\n", correctedPageSum)
}

func PartOne(rules [][]int, updates [][]int) int {
	pageSum := 0
	incorrectUpdates := [][]int{}

	for _, update := range updates {
		if updateIsCorrect(update, rules) {
			pageSum += middlePageValue(update)
		} else {
			incorrectUpdates = append(incorrectUpdates, update)
		}
	}

	return pageSum
}

func PartTwo(rules [][]int, updates [][]int) int {
	incorrectUpdates := [][]int{}

	for _, update := range updates {
		if !updateIsCorrect(update, rules) {
			incorrectUpdates = append(incorrectUpdates, update)
		}
	}

	correctedPageSum := 0
	for _, update := range incorrectUpdates {
		ruleIndex := 0

		for {
			if updateIsCorrect(update, rules) {
				break
			}

			if ruleIndex >= len(rules) {
				ruleIndex = 0
			}

			rule := rules[ruleIndex]

			// if we fail a rule, swap the numbers around
			if !ruleIsFollowed(rule, update) {
				firstIndex := slices.Index(update, rule[0])
				secondIndex := slices.Index(update, rule[1])

				tmp := update[firstIndex]
				update[firstIndex] = update[secondIndex]
				update[secondIndex] = tmp
			}

			ruleIndex += 1
		}

		correctedPageSum += middlePageValue(update)
	}

	return correctedPageSum
}

func updateIsCorrect(update []int, rules [][]int) bool {
	for _, rule := range rules {
		if !ruleIsFollowed(rule, update) {
			return false
		}
	}

	return true
}

func ruleIsFollowed(rule []int, update []int) bool {
	firstIndex := slices.Index(update, rule[0])
	secondIndex := slices.Index(update, rule[1])

	if firstIndex == -1 || secondIndex == -1 {
		return true
	}

	return firstIndex < secondIndex
}

func middlePageValue(update []int) int {
	middleIndex := int((len(update) - 1) / 2)
	return update[middleIndex]
}

func readInput(filePath string) ([][]int, [][]int) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, nil
	}
	defer file.Close()

	var rules [][]int
	var updates [][]int

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		// parse a rule
		if strings.Contains(line, "|") {
			strParts := strings.Split(line, "|")
			rule := []int{}

			for _, strPart := range strParts {
				var part int
				part, _ = strconv.Atoi(strPart)
				rule = append(rule, part)
			}

			rules = append(rules, rule)
		}

		// parse an update
		if strings.Contains(line, ",") {
			strParts := strings.Split(line, ",")
			update := []int{}

			for _, strPart := range strParts {
				var part int
				part, _ = strconv.Atoi(strPart)
				update = append(update, part)
			}

			updates = append(updates, update)
		}
	}

	return rules, updates
}
