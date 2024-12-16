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

	// part 1
	lab := readInput(fileName)

	for {
		nextLocation := lab.guard.nextLocation()

		if !lab.isInArea(nextLocation.x, nextLocation.y) {
			break
		}

		if lab.isObstructed(nextLocation.x, nextLocation.y) {
			lab.guard.turn()
		} else {
			lab.guard.move()
			lab.markVisited()
		}
	}

	fmt.Printf("Guard will visit %d positions before leaving the area\n", lab.visitedLocationCount())

	// part 2
	waysToLoop := 0

	for x := 0; x < len(lab.grid[0]); x++ {
		for y := 0; y < len(lab.grid); y++ {
			lab := readInput(fileName)

			if lab.isObstructed(x, y) {
				continue
			} else {
				lab.grid[y][x] = "O"
			}

			causesLoop := false

			for {
				nextLocation := lab.guard.nextLocation()

				if !lab.isInArea(nextLocation.x, nextLocation.y) {
					break
				}

				if lab.isObstructed(nextLocation.x, nextLocation.y) {
					lab.guard.turn()
				} else {
					lab.guard.move()
					isNewPath := lab.markVisited()

					if !isNewPath {
						causesLoop = true
						break
					}
				}
			}

			if causesLoop {
				waysToLoop += 1
			}
		}
	}

	fmt.Printf("There are %d useful positions for the new obstruction.\n", waysToLoop)
}

func readInput(filePath string) area {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening %s", filePath)
		panic(err)
	}
	defer file.Close()

	grid := [][]string{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		row := strings.Split(line, "")

		grid = append(grid, row)
	}

	mapArea := area{grid: grid}
	mapArea.locateGuard()
	mapArea.grid[mapArea.guard.y][mapArea.guard.x] = "h"
	return mapArea
}

type entity struct {
	x         int
	y         int
	direction string
}

func (e *entity) draw() string {
	switch e.direction {
	case "N":
		return "^"
	case "E":
		return ">"
	case "S":
		return "v"
	case "W":
		return "<"
	}

	panic("Invalid entity direction")
}

func (e *entity) turn() {
	switch e.direction {
	case "N":
		e.direction = "E"
		break
	case "E":
		e.direction = "S"
		break
	case "S":
		e.direction = "W"
		break
	case "W":
		e.direction = "N"
		break
	}
}

func (e entity) nextLocation() entity {
	switch e.direction {
	case "N":
		e.y -= 1
	case "E":
		e.x += 1
	case "S":
		e.y += 1
	case "W":
		e.x -= 1
	}

	return e
}

func (e *entity) move() {
	switch e.direction {
	case "N":
		e.y -= 1
	case "E":
		e.x += 1
	case "S":
		e.y += 1
	case "W":
		e.x -= 1
	}
}

type area struct {
	grid  [][]string
	guard entity
}

func (m *area) draw() {
	for y, row := range m.grid {
		for x, cell := range row {
			if x == m.guard.x && y == m.guard.y {
				fmt.Print(m.guard.draw())
			} else {
				fmt.Print(cell)
			}
		}
		fmt.Println()
	}
}

func (m *area) locateGuard() {
	for y, row := range m.grid {
		for x, cell := range row {
			switch cell {
			case "^":
				m.guard = entity{x: x, y: y, direction: "N"}
			case "v":
				m.guard = entity{x: x, y: y, direction: "S"}
			case "<":
				m.guard = entity{x: x, y: y, direction: "W"}
			case ">":
				m.guard = entity{x: x, y: y, direction: "E"}
			}
		}
	}

	m.markVisited()
}

func (m *area) isObstructed(x int, y int) bool {
	cell := m.grid[y][x]
	return cell == "#" || cell == "O"
}

func (m *area) markVisited() bool {
	prevChar := m.grid[m.guard.y][m.guard.x]

	var prevDirections []string

	if prevChar == "." {
		prevDirections = []string{}
	} else {
		prevDirections = charToDirections(prevChar)
	}

	if slices.Contains(prevDirections, m.guard.direction) {
		return false
	}

	prevDirections = append(prevDirections, m.guard.direction)
	m.grid[m.guard.y][m.guard.x] = directionsToChar(prevDirections)

	return true
}

func (m *area) visitedLocationCount() int {
	count := 0

	for _, row := range m.grid {
		for _, cell := range row {
			if cell != "." && cell != "#" {
				count += 1
			}
		}
	}

	return count
}

func (m area) isInArea(x int, y int) bool {
	return x >= 0 && x < len(m.grid) && y >= 0 && y < len(m.grid[0])
}

func charToDirections(char string) []string {
	charCode := int(char[0])
	binary := strconv.FormatInt(int64(charCode-96), 2)
	binaryStr := fmt.Sprintf("%04s", binary)
	binarySlice := strings.Split(binaryStr, "")

	directions := []string{}

	if binarySlice[0] == "1" {
		directions = append(directions, "N")
	}

	if binarySlice[1] == "1" {
		directions = append(directions, "E")
	}

	if binarySlice[2] == "1" {
		directions = append(directions, "S")
	}

	if binarySlice[3] == "1" {
		directions = append(directions, "W")
	}

	return directions
}

func directionsToChar(directions []string) string {
	directionsMap := map[string]int{
		"N": 8,
		"E": 4,
		"S": 2,
		"W": 1,
	}

	charCode := 0

	for _, direction := range directions {
		charCode += directionsMap[direction]
	}

	return string(rune(charCode + 96))
}
