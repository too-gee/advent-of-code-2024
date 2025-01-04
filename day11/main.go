package main

import (
	"bufio"
	"fmt"
	"math"
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

	stones := readInput(fileName)

	// part 1
	stoneCount := PartOne(stones)
	fmt.Printf("Stone count: %d\n", stoneCount)

	// part 2
	stoneCount = PartTwo(stones)
	fmt.Printf("Stone count after a while: %d\n", stoneCount)
}

func PartOne(input Stones) int {
	stones := input.copy()
	maxBlinks := 25

	for blink := 1; blink <= maxBlinks; blink++ {
		for i := 0; i < len(stones); i++ {
			if stones[i] == 0 {
				stones[i] = 1
				continue
			}

			if stones[i] > 0 && digits(stones[i])%2 == 0 {
				left, right := split(stones[i])

				stones[i] = left
				stones = slices.Insert(stones, i+1, right)

				// move forward one spot to account for the split
				i += 1
				continue
			}

			stones[i] = stones[i] * 2024
		}

		fmt.Printf("Stone count after %d blinks: %d\n", blink, len(stones))
	}

	return len(stones)
}

func PartTwo(input Stones) int {
	stones := input.copy()
	maxBlinks := 75

	metaStones := changeSet{}

	for _, stone := range stones {
		metaStones.add(stone, 1)
	}

	for blink := 1; blink <= maxBlinks; blink++ {
		blinkChanges := changeSet{}

		for value, count := range metaStones {
			if value == 0 {
				blinkChanges.add(0, -count)
				blinkChanges.add(1, count)
				continue
			}

			if digits(value)%2 == 0 {
				left, right := split(value)

				blinkChanges.add(value, -count)
				blinkChanges.add(left, count)
				blinkChanges.add(right, count)
				continue
			}

			blinkChanges.add(value, -count)
			blinkChanges.add(value*2024, count)
		}

		metaStones.combine(blinkChanges)

		fmt.Printf("Stone count after %d optimized blinks: %d\n", blink, metaStones.count())
	}

	return metaStones.count()
}

func readInput(filePath string) []int {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening %s", filePath)
		return nil
	}
	defer file.Close()

	stones := []int{}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		row := strings.Split(line, " ")

		rowInts := []int{}
		for _, value := range row {
			intValue, _ := strconv.Atoi(value)
			rowInts = append(rowInts, intValue)
		}

		stones = append(stones, rowInts...)
	}

	return stones
}

type Stones []int

func (s Stones) copy() []int {
	newStones := make([]int, len(s))
	copy(newStones, s)
	return newStones
}

func digits(num int) int {
	return int(math.Log10(float64(num)) + 1)
}

func split(num int) (int, int) {
	factor := int(math.Pow(10, float64(digits(num)/2)))

	left := num / factor
	right := num % factor

	return left, right
}

type changeSet map[int]int

func (c *changeSet) add(stone int, count int) {
	if _, ok := (*c)[stone]; ok {
		(*c)[stone] += count
	} else {
		(*c)[stone] = count
	}
}

func (c *changeSet) combine(other changeSet) {
	for stone, count := range other {
		(*c).add(stone, count)
	}

	for stone, count := range *c {
		if count == 0 {
			delete(*c, stone)
		}
	}
}

func (c changeSet) count() int {
	totalCount := 0

	for _, count := range c {
		totalCount += count
	}

	return totalCount
}
