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
	// Open the file
	list1, list2 := readFileToLists(os.Args[1])

	// part 1
	differences := 0

	sort.Ints(list1)
	sort.Ints(list2)

	for i := range list1 {
		differences += int(math.Abs(float64(list1[i] - list2[i])))
	}

	fmt.Printf("The total difference is %d\n", differences)

	// part 2
	similarity := 0

	for _, needle := range list1 {
		for _, haystack := range list2 {
			if needle == haystack {
				similarity += needle
			}
		}
	}

	fmt.Printf("The similarity score is %d\n", similarity)
}

func readFileToLists(filePath string) ([]int, []int) {
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
