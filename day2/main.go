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

	reports := readInput(fileName)

	// part 1
	safeReports := 0
	unsafeReports := [][]int{}

	for _, report := range reports {
		if reportIsSafe(report, false) {
			safeReports++
		} else {
			unsafeReports = append(unsafeReports, report)
		}
	}

	fmt.Printf("Safe reports: %d\n", safeReports)

	// part 2

	safeDamped := 0

	for _, report := range unsafeReports {
		for i := range report {
			dampedReport := excludeIndex(report, i)
			if reportIsSafe(dampedReport, false) {
				safeDamped += 1
				break
			}
		}
	}

	fmt.Printf("Reports made safe by damping: %d\n", safeDamped)
	fmt.Printf("Total safe reports after damping: %d\n", safeReports+safeDamped)
}

func compare(first int, second int) int {
	if second > first {
		return 1
	} else if first > second {
		return -1
	} else {
		return 0
	}
}

func reportIsSafe(report []int, debug bool) bool {
	var currentDirection int

	reportIsSafe := true
	debugOutput := "debug ->"

	for j, currentValue := range report {
		debugOutput += " " + strconv.Itoa(currentValue)

		if j == 0 {
			// set up in first loop
			currentDirection = compare(report[0], report[1])
			continue
		}

		previousValue := report[j-1]

		previousDirection := currentDirection
		currentDirection = compare(previousValue, currentValue)

		distance := (currentValue - previousValue) * currentDirection

		if distance < 1 || distance > 3 {
			if debug {
				fmt.Printf("%s -- difference too big\n", debugOutput)
			}
			reportIsSafe = false
			break
		}

		if previousDirection != currentDirection {
			if debug {
				fmt.Printf("%s -- direction changed\n", debugOutput)
			}
			reportIsSafe = false
			break
		}
	}

	return reportIsSafe
}

func excludeIndex(slice []int, index int) []int {
	newSlice := make([]int, 0)
	newSlice = append(newSlice, slice[:index]...)
	newSlice = append(newSlice, slice[index+1:]...)
	return newSlice
}

func readInput(filePath string) [][]int {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var reports [][]int

	for scanner.Scan() {
		line := scanner.Text()
		rawValues := strings.Fields(line)

		var values []int
		for _, rawValue := range rawValues {
			value, _ := strconv.Atoi(rawValue)
			values = append(values, value)
		}

		reports = append(reports, values)
	}

	return reports
}
