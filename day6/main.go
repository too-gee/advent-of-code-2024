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

	lab := readInput(fileName)

	// part 1
	visitedLocations := PartOne(lab)

	fmt.Printf("Guard will visit %d positions before leaving the area\n", visitedLocations)

	// part 2
	waysToLoop := PartTwo(lab)

	fmt.Printf("There are %d useful positions for the new obstruction.\n", waysToLoop)
}

func PartOne(input area) int {
	lab := input.copy()

	for {
		nextLocation := lab.guard.nextLocation()

		if !lab.Contains(nextLocation) {
			break
		}

		if lab.isObstructed(nextLocation.X, nextLocation.Y) {
			lab.guard.turn()
		} else {
			lab.guard.move()
			lab.markVisited()
		}
	}

	return lab.visitedLocationCount()
}

func PartTwo(input area) int {
	waysToLoop := 0

	for x := 0; x < len(input.Grid[0]); x++ {
		for y := 0; y < len(input.Grid); y++ {
			lab := input.copy()

			if lab.isObstructed(x, y) {
				continue
			} else {
				lab.Grid[y][x] = "O"
			}

			causesLoop := false

			for {
				nextLocation := lab.guard.nextLocation()

				if !lab.Contains(nextLocation) {
					break
				}

				if lab.isObstructed(nextLocation.X, nextLocation.Y) {
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

	return waysToLoop
}

func readInput(filePath string) area {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening %s", filePath)
		panic(err)
	}
	defer file.Close()

	grid := shared.Grid{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		row := strings.Split(line, "")

		grid = append(grid, row)
	}

	mapArea := area{Grid: grid}
	mapArea.locateGuard()
	mapArea.Grid[mapArea.guard.Y][mapArea.guard.X] = "h"
	return mapArea
}

type Entity struct {
	shared.Coord
	direction string
}

func (e *Entity) draw() string {
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

	panic("Invalid Entity direction")
}

func (e *Entity) turn() {
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

func (e Entity) nextLocation() Entity {
	switch e.direction {
	case "N":
		e.Y -= 1
	case "E":
		e.X += 1
	case "S":
		e.Y += 1
	case "W":
		e.X -= 1
	}

	return e
}

func (e *Entity) move() {
	switch e.direction {
	case "N":
		e.Y -= 1
	case "E":
		e.X += 1
	case "S":
		e.Y += 1
	case "W":
		e.X -= 1
	}
}

type area struct {
	shared.Grid
	guard Entity
}

func (m area) copy() area {
	newGrid := make(shared.Grid, m.Grid.Height())
	for y, row := range m.Grid {
		newRow := make([]string, m.Grid.Width())
		copy(newRow, row)
		newGrid[y] = newRow
	}

	newGuard := Entity{Coord: shared.Coord{X: m.guard.X, Y: m.guard.Y}, direction: m.guard.direction}

	return area{Grid: newGrid, guard: newGuard}
}

func (m *area) draw() {
	for y, row := range m.Grid {
		for x, cell := range row {
			if x == m.guard.X && y == m.guard.Y {
				fmt.Print(m.guard.draw())
			} else {
				fmt.Print(cell)
			}
		}
		fmt.Println()
	}
}

func (m *area) locateGuard() {
	for y, row := range m.Grid {
		for x, cell := range row {
			switch cell {
			case "^":
				m.guard = Entity{Coord: shared.Coord{X: x, Y: y}, direction: "N"}
			case "v":
				m.guard = Entity{Coord: shared.Coord{X: x, Y: y}, direction: "S"}
			case "<":
				m.guard = Entity{Coord: shared.Coord{X: x, Y: y}, direction: "W"}
			case ">":
				m.guard = Entity{Coord: shared.Coord{X: x, Y: y}, direction: "E"}
			}
		}
	}

	m.markVisited()
}

func (m *area) isObstructed(x int, y int) bool {
	cell := m.Grid[y][x]
	return cell == "#" || cell == "O"
}

func (m *area) markVisited() bool {
	prevChar := m.Grid[m.guard.Y][m.guard.X]

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
	m.Grid[m.guard.Y][m.guard.X] = directionsToChar(prevDirections)

	return true
}

func (m *area) visitedLocationCount() int {
	count := 0

	for _, row := range m.Grid {
		for _, cell := range row {
			if cell != "." && cell != "#" {
				count += 1
			}
		}
	}

	return count
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
