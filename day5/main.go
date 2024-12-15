package main

import (
	"bufio"
	"fmt"
	"os"
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
	pageSum := 0

	for _, update := range updates {
		isCorrect := true
		for _, rule := range rules {
			if getIndex(update, rule[0]) == -1 || getIndex(update, rule[1]) == -1 {
				continue
			}

			if getIndex(update, rule[0]) >= getIndex(update, rule[1]) {
				isCorrect = false
				break
			}
		}

		if isCorrect {
			middleIndex := int((len(update) - 1) / 2)
			pageSum += update[middleIndex]
		}
	}

	fmt.Printf("The sum of the middle page numbers is %d\n", pageSum)
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

func getIndex(slice []int, value int) int {
	for i, v := range slice {
		if v == value {
			return i
		}
	}
	return -1
}
