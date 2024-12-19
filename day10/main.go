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

	topoMap := readInput(fileName)

	// part 1
	trailHeads := topoMap.trailHeads()

	totalScore := 0

	for _, trailHead := range trailHeads {
		totalScore += len(topoMap.trailEnds(trailHead))
	}

	fmt.Printf("The sum of the scores of all trailheads is %d\n", totalScore)

	// part 2
	totalRating := 0

	for _, trailHead := range trailHeads {
		totalRating += topoMap.trailPaths(trailHead)
	}

	fmt.Printf("The sum of the ratings of all trailheads is %d\n", totalRating)
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

type Coord struct {
	x int
	y int
}

func (m Map) value(loc Coord) int {
	if loc.x < 0 || loc.x >= len(m[0]) || loc.y < 0 || loc.y >= len(m) {
		return -1
	}

	return m[loc.y][loc.x]
}

func (m Map) draw() {
	for _, row := range m {
		for _, cell := range row {
			fmt.Print(cell)
		}
		fmt.Println()
	}
}

func (m Map) trailHeads() []Coord {
	var trailHeads []Coord

	for y, row := range m {
		for x, cell := range row {
			if cell == 0 {
				trailHeads = append(trailHeads, Coord{x: x, y: y})
			}
		}
	}

	return trailHeads
}

func (m Map) neighbors(loc Coord) []Coord {
	neighbors := []Coord{}

	offsets := []Coord{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}

	for _, offset := range offsets {
		candidate := Coord{x: loc.x + offset.x, y: loc.y + offset.y}
		if m.value(candidate) == m.value(loc)+1 {
			neighbors = append(neighbors, candidate)
		}
	}

	return neighbors
}

func (m Map) trailEnds(loc Coord) []Coord {
	if m.value(loc) == -1 {
		return nil
	}

	if m.value(loc) == 9 {
		return []Coord{loc}
	}

	trailEnds := []Coord{}

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

func (m Map) trailPaths(loc Coord) int {
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
