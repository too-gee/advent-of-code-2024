package main

import (
	"bufio"
	"fmt"
	"math"
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

	equations := readInput(fileName)

	// part 1
	totalCalibrationResult := 0

	for _, equation := range equations {
		permutations := permutations(len(equation)-2, []string{"+", "*"})

		solvable := false

		for _, permutation := range permutations {
			runningResult := equation[1]

			for i := 2; i < len(equation); i++ {
				operator := permutation[i-2]

				switch operator {
				case "+":
					runningResult += equation[i]
				case "*":
					runningResult *= equation[i]
				}

				if runningResult > equation[0] {
					break
				}
			}

			if runningResult == equation[0] {
				solvable = true
				break
			}
		}

		if solvable {
			totalCalibrationResult += equation[0]
		}
	}

	fmt.Printf("The total calibration result is %d\n", totalCalibrationResult)

	// part 2
	totalCalibrationResult = 0

	for _, equation := range equations {
		permutations := permutations(len(equation)-2, []string{"+", "*", "||"})

		solvable := false

		for _, permutation := range permutations {
			runningResult := equation[1]

			for i := 2; i < len(equation); i++ {
				operator := permutation[i-2]

				switch operator {
				case "+":
					runningResult += equation[i]
				case "*":
					runningResult *= equation[i]
				case "||":
					strA := strconv.Itoa(runningResult)
					strB := strconv.Itoa(equation[i])
					newInt, _ := strconv.Atoi(strA + strB)
					runningResult = newInt
				}

				if runningResult > equation[0] {
					break
				}
			}

			if runningResult == equation[0] {
				solvable = true
				break
			}
		}

		if solvable {
			totalCalibrationResult += equation[0]
		}
	}

	fmt.Printf("The NEW total calibration result is %d\n", totalCalibrationResult)
}

func readInput(filePath string) [][]int {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening %s", filePath)
		return nil
	}
	defer file.Close()

	equations := [][]int{}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.Replace(scanner.Text(), ":", "", 1)
		rowStr := strings.Split(line, " ")

		row := []int{}
		for _, str := range rowStr {
			var num int
			num, _ = strconv.Atoi(str)
			row = append(row, num)
		}

		equations = append(equations, row)
	}

	return equations
}

func permutations(length int, options []string) [][]string {
	permutations := [][]string{}

	for i := 0; i < int(math.Pow(float64(len(options)), float64(length))); i++ {
		binary := strconv.FormatInt(int64(i), len(options))
		binaryStr := fmt.Sprintf("%0"+strconv.Itoa(length)+"s", binary)
		binarySlice := strings.Split(binaryStr, "")
		strSlice := []string{}

		for _, str := range binarySlice {
			intVal, _ := strconv.Atoi(str)
			strSlice = append(strSlice, options[intVal])
		}

		permutations = append(permutations, strSlice)
	}

	return permutations
}
