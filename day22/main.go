package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	var fileName string

	if len(os.Args) == 2 {
		fileName = os.Args[1]
	} else {
		fileName = "input.txt"
	}

	initNums := readInput(fileName)

	fmt.Println(initNums)
}

func readInput(filePath string) []int {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening %s", filePath)
		return nil
	}
	defer file.Close()

	initNums := []int{}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		lineValue, _ := strconv.Atoi(line)

		initNums = append(initNums, lineValue)
	}

	return initNums
}

func Solve(initNums []int, iters int) int {
	return 0
}
