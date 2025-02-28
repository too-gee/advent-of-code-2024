package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"slices"
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

	maze := readInput(fileName)

	// Part 1
	if strings.HasSuffix(fileName, "small.txt") {
		fmt.Printf("Part 1: %d cheats save over 64us with 2us cheats\n", Solve(maze, 64, 2))
		fmt.Printf("Part 2: %d cheats save over 76us with 20us cheats\n", Solve(maze, 76, 20))
	} else {
		fmt.Printf("Part 1: %d cheats save over 100us with 2us cheats\n", Solve(maze, 100, 2))
		fmt.Printf("Part 2: %d cheats save over 100us with 20us cheats\n", Solve(maze, 100, 20))
	}

}

func readInput(filePath string) shared.Grid {
	file, _ := os.Open(filePath)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	grid := shared.Grid{}

	for scanner.Scan() {
		line := scanner.Text()

		gridRow := strings.Split(line, "")
		grid = append(grid, gridRow)
	}

	return grid
}

func Solve(maze shared.Grid, minSavings int, cheatLength int) int {
	remaining := GetRemainingLengths(maze)

	cheats := 0
	passable := []string{EMPTY, END}

	for start, value := range remaining {
		for dx := -1 * cheatLength; dx <= cheatLength; dx++ {
			for dy := -1 * cheatLength; dy <= cheatLength; dy++ {
				end := shared.Coord{X: start.X + dx, Y: start.Y + dy}
				duringCheat := int(math.Abs(float64(dx)) + math.Abs(float64(dy)))

				if !(dx == 0 && dy == 0) && // is not the start
					duringCheat <= cheatLength && // is within cheat range
					maze.Contains(end) && // is within maze
					slices.Contains(passable, maze.At(end)) && // is passable
					value-remaining[end]-duringCheat >= minSavings { // is a savings
					cheats++
				}
			}
		}
	}

	return cheats
}

func MazeLength(maze shared.Grid) int {
	totalLength := 1

	for _, row := range maze {
		for _, cell := range row {
			if cell == EMPTY {
				totalLength++
			}
		}
	}

	return totalLength
}

func GetRemainingLengths(maze shared.Grid) map[shared.Coord]int {
	remaining := map[shared.Coord]int{}
	crawler := []shared.Coord{maze.LocationOf(START)}

	for i := MazeLength(maze); i >= 0; i-- {
		cur := crawler[len(crawler)-1]
		remaining[cur] = i

		for _, neighbor := range maze.Neighbors(cur, []string{WALL}) {
			if slices.Contains(crawler, neighbor) {
				continue
			}

			crawler = append(crawler, neighbor)
			break
		}
	}

	return remaining
}

const WALL = "#"
const EMPTY = "."
const START = "S"
const END = "E"
