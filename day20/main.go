package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

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
	fmt.Println(Part1(maze, 100))

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

func Part1(maze shared.Grid, minSavings int) int {
	MarkPotentialCheats(&maze)
	remaining := GetRemainingLengths(maze)

	cheats := 0

	for x, row := range maze {
		for y := range row {
			cur := shared.Coord{X: x, Y: y}

			if CHEAT == maze.At(cur) {
				// Get all values of neighbors of this cheat location
				vals := []int{}

				for _, neighbor := range maze.Neighbors(cur, []string{WALL, CHEAT}) {
					vals = append(vals, remaining[neighbor])
				}

				// The smallest value must be our start and the other values must be the ends
				start := slices.Min(vals)
				ends := []int{}

				for x := range vals {
					if vals[x] != start {
						ends = append(ends, vals[x])
					}
				}

				// Count how many times the savings are greater than the minimum savings
				for _, end := range ends {
					if end-start > minSavings {
						cheats++
					}
				}
			}
		}
	}

	return cheats
}

func Solve(g shared.Grid, savings int) []State {
	totalLength := 1

	for _, row := range g {
		for _, cell := range row {
			if cell == EMPTY {
				totalLength++
			}
		}
	}

	maxLength := totalLength - savings

	start := g.LocationOf(START)
	end := g.LocationOf(END)

	queue := DumbQueue{}
	queue.push(State{
		cheated: false,
		path:    []shared.Coord{start},
	})

	finished := []State{}

	for len(queue) > 0 {
		current := queue.pop()

		// if we're at the end, save and continue
		if current.Loc() == end {
			finished = append(finished, current)
			continue
		}

		if current.Len() > maxLength {
			continue
		}

		var blockers []string

		if current.cheated {
			blockers = []string{WALL, CHEAT}
		} else {
			blockers = []string{WALL}
		}

		for _, neighbor := range g.Neighbors(current.Loc(), blockers) {
			if slices.Contains(current.path, neighbor) {
				continue
			}

			var newCheated bool

			if g.At(neighbor) == CHEAT {
				newCheated = true
			} else {
				newCheated = current.cheated
			}

			queue.push(State{
				cheated: newCheated,
				path:    CopyAppend(current.path, neighbor),
			})
		}
	}

	return finished
}

func CopyAppend(a []shared.Coord, b shared.Coord) []shared.Coord {
	newSlice := make([]shared.Coord, len(a)+1)
	copy(newSlice, a)
	newSlice[len(a)] = b

	return newSlice
}

const WALL = "#"
const EMPTY = "."
const START = "S"
const END = "E"
const CHEAT = "O"

type State struct {
	cheated bool
	path    []shared.Coord
}

func (s State) Loc() shared.Coord {
	return s.path[len(s.path)-1]
}

func now() int64 {
	currentTime := time.Now()
	return currentTime.UnixNano() / int64(time.Millisecond)
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

func MarkPotentialCheats(maze *shared.Grid) {
	open := []string{EMPTY, START, END}
	rotateDir := []string{"R", "L"}

	for _, dir := range rotateDir {
		for y := 1; y < len(*maze)-1; y++ {
			row := (*maze)[y]
			for x := 1; x < len(row)-3; x++ {
				if row[x+1] == WALL && slices.Contains(open, row[x]) && slices.Contains(open, row[x+2]) {
					(*maze)[y][x+1] = CHEAT
				}
			}
		}

		(*maze).Rotate(dir)
	}
}

func GetRemainingLengths(maze shared.Grid) map[shared.Coord]int {
	remaining := map[shared.Coord]int{}
	crawler := []shared.Coord{maze.LocationOf(START)}

	for i := MazeLength(maze); i > 0; i-- {
		cur := crawler[len(crawler)-1]
		remaining[cur] = i

		for _, neighbor := range maze.Neighbors(cur, []string{WALL, CHEAT}) {
			if slices.Contains(crawler, neighbor) {
				continue
			}

			crawler = append(crawler, neighbor)
			break
		}
	}

	return remaining
}

/*func (s State) CheatsLeft() int {
	if s.cheatStart == -1 {
		return 2
	}

	return int(math.Max(0, float64(s.cheatStart-len(s.path)+1)))
}*/

func (s State) Len() int {
	return len(s.path) - 1
}

type DumbQueue []State

func (q *DumbQueue) push(item State) { *q = append(*q, item) }

func (q *DumbQueue) pop() State {
	item := (*q)[len(*q)-1]
	*q = (*q)[0 : len(*q)-1]
	return item
}
