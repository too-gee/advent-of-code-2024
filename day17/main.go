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

	registers, program := readInput(fileName)

	fmt.Println(registers)
	fmt.Println(program)
}

func readInput(filePath string) ([]int, []int) {
	file, _ := os.Open(filePath)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var a, b, c int
	var program []int

	for scanner.Scan() {
		line := scanner.Text()

		fmt.Sscanf(line, "Register A: %d", &a)
		fmt.Sscanf(line, "Register B: %d", &b)
		fmt.Sscanf(line, "Register C: %d", &c)

		var tmpProgram string
		_, err := fmt.Sscanf(line, "Program: %s", &tmpProgram)

		if err == nil {
			tmpProgram := strings.Split(tmpProgram, ",")
			for _, tmpStr := range tmpProgram {
				tmpInt, _ := strconv.Atoi(tmpStr)
				program = append(program, tmpInt)
			}
		}
	}

	return []int{a, b, c}, program
}

func Solve(registers []int, program []int) ([]int, string) {
	return []int{0, 0, 0}, ""
}
