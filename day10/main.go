package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/too-gee/advent-of-code-2024/shared"
)

func main() {
	var fileName string

	if len(os.Args) == 2 {
		fileName = os.Args[1]
	} else {
		fileName = "input.txt"
	}

	topoMap := readInput(fileName)

	// part 1
	totalScore := PartOne(topoMap)
	fmt.Printf("The sum of the scores of all trailheads is %d\n", totalScore)

	// part 2
	totalRating := PartTwo(topoMap)
	fmt.Printf("The sum of the ratings of all trailheads is %d\n", totalRating)
}

func PartOne(topoMap Map) int {
	totalScore := 0

	trailHeads := topoMap.trailHeads()

	for _, trailHead := range trailHeads {
		totalScore += len(topoMap.trailEnds(trailHead))
	}

	return totalScore
}

func PartTwo(topoMap Map) int {
	totalRating := 0

	trailHeads := topoMap.trailHeads()

	for _, trailHead := range trailHeads {
		totalRating += topoMap.trailPaths(trailHead)
	}

	return totalRating
}

func readInput(filePath string) Map {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening %s", filePath)
		return nil
	}
	defer file.Close()

	topoMap := Map{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		row := strings.Split(line, "")

		rowInts := []int{}
		for _, value := range row {
			intValue, _ := strconv.Atoi(value)
			rowInts = append(rowInts, intValue)
		}

		topoMap = append(topoMap, rowInts)
	}

	return topoMap
}

type Map [][]int

func (m Map) value(loc shared.Coord) int {
	if loc.X < 0 || loc.X >= len(m[0]) || loc.Y < 0 || loc.Y >= len(m) {
		return -1
	}

	return m[loc.Y][loc.X]
}

func (m Map) draw() {
	for _, row := range m {
		for _, cell := range row {
			fmt.Print(cell)
		}
		fmt.Println()
	}
}

func (m Map) trailHeads() []shared.Coord {
	var trailHeads []shared.Coord

	for y, row := range m {
		for x, cell := range row {
			if cell == 0 {
				trailHeads = append(trailHeads, shared.Coord{X: x, Y: y})
			}
		}
	}

	return trailHeads
}

func (m Map) neighbors(loc shared.Coord) []shared.Coord {
	neighbors := []shared.Coord{}

	offsets := []shared.Coord{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}

	for _, offset := range offsets {
		candidate := shared.Coord{X: loc.X + offset.X, Y: loc.Y + offset.Y}
		if m.value(candidate) == m.value(loc)+1 {
			neighbors = append(neighbors, candidate)
		}
	}

	return neighbors
}

func (m Map) trailEnds(loc shared.Coord) []shared.Coord {
	if m.value(loc) == -1 {
		return nil
	}

	if m.value(loc) == 9 {
		return []shared.Coord{loc}
	}

	trailEnds := []shared.Coord{}

	for _, nextLoc := range m.neighbors(loc) {
		newEnds := m.trailEnds(nextLoc)

		for _, newEnd := range newEnds {
			if !slices.Contains(trailEnds, newEnd) {
				trailEnds = append(trailEnds, newEnd)
			}
		}
	}

	return trailEnds
}

func (m Map) trailPaths(loc shared.Coord) int {
	if m.value(loc) == -1 {
		return 0
	}

	if m.value(loc) == 9 {
		return 1
	}

	trailRating := 0

	for _, nextLoc := range m.neighbors(loc) {
		trailRating += m.trailPaths(nextLoc)
	}

	return trailRating
}
