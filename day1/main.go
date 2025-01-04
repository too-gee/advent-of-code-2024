package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
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

	list1, list2 := readInput(fileName)

	// part 1
	difference := PartOne(list1, list2)
	fmt.Printf("The total difference is %d\n", difference)

	// part 2
	similarity := PartTwo(list1, list2)
	fmt.Printf("The similarity score is %d\n", similarity)
}

func PartOne(list1 []int, list2 []int) int {
	difference := 0

	sort.Ints(list1)
	sort.Ints(list2)

	for i := range list1 {
		difference += int(math.Abs(float64(list1[i] - list2[i])))
	}

	return difference
}

func PartTwo(list1 []int, list2 []int) int {
	similarity := 0

	for _, needle := range list1 {
		for _, haystack := range list2 {
			if needle == haystack {
				similarity += needle
			}
		}
	}

	return similarity
}

func readInput(filePath string) ([]int, []int) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	list1 := make([]int, 0)
	list2 := make([]int, 0)

	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()

		parts := strings.Fields(line)

		num1, _ := strconv.Atoi(parts[0])
		num2, _ := strconv.Atoi(parts[1])

		list1 = append(list1, num1)
		list2 = append(list2, num2)
	}

	return list1, list2
}
