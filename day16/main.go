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

	// part 1 & 2
	bestCost, optimalTiles := Solve(maze)
	fmt.Printf("Lowest cost: %d\n", bestCost)
	fmt.Printf("Optimal path count: %d\n", optimalTiles)
}

func readInput(filePath string) Maze {
	file, _ := os.Open(filePath)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	grid := Maze{}

	for scanner.Scan() {
		line := scanner.Text()

		gridRow := strings.Split(line, "")
		grid.Grid = append(grid.Grid, gridRow)
	}

	return grid
}

func Solve(m Maze) (int, int) {
	start := m.LocationOf(START)
	end := m.LocationOf(END)

	queue := DumbQueue{}
	queue.push(State{
		dir:  "E",
		cost: 0,
		path: []shared.Coord{start},
	})

	bestCost := math.MaxInt
	visitedCost := map[shared.Coord]int{}
	bestPaths := [][]shared.Coord{}

	for len(queue) > 0 {
		current := queue.pop()

		// fill a cost grid (not checking direction)
		localBestCost, ok := visitedCost[current.Loc()]

		if !ok || current.cost < localBestCost {
			visitedCost[current.Loc()] = current.cost
			localBestCost = current.cost
		}

		// if another path has crossed here with a much lower cost, skip
		if current.cost > bestCost || current.cost > (localBestCost+1000) {
			continue
		}

		// if we're at the end, save and continue
		if current.Loc() == end {
			if current.cost < bestCost {
				bestPaths = [][]shared.Coord{current.path}
				bestCost = current.cost
			} else if current.cost == bestCost {
				bestPaths = append(bestPaths, current.path)
			}
			continue
		}

		for newDir, neighbor := range m.Neighbors(current.Loc()) {
			if slices.Contains(current.path, neighbor) {
				continue
			}

			newCost := current.cost + 1

			if newDir != current.dir {
				newCost += 1000
			}

			queue.push(State{
				dir:  newDir,
				cost: newCost,
				path: CopyAppend(current.path, neighbor),
			})
		}
	}

	bestTiles := []shared.Coord{}
	for _, path := range bestPaths {
		for _, loc := range path {
			if !slices.Contains(bestTiles, loc) {
				bestTiles = append(bestTiles, loc)
			}
		}
	}

	m.Draw(map[string]string{"S": "🟢", "E": "🔴"}, map[string][]shared.Coord{"⏺️ ": bestTiles})

	return bestCost, len(bestTiles)
}

func CopyAppend(a []shared.Coord, b shared.Coord) []shared.Coord {
	newSlice := make([]shared.Coord, len(a)+1)
	copy(newSlice, a)
	newSlice[len(a)] = b

	return newSlice
}

const FILL = "="
const WALL = "#"
const EMPTY = "."
const START = "S"
const END = "E"

type Maze struct {
	shared.Grid
}

func (m Maze) Neighbors(loc shared.Coord) map[string]shared.Coord {
	return m.Grid.Neighbors(loc, []string{WALL, FILL})
}

type State struct {
	dir  string
	cost int
	path []shared.Coord
}

func (s State) Loc() shared.Coord {
	return s.path[len(s.path)-1]
}

type DumbQueue []State

func (q *DumbQueue) push(item State) { *q = append(*q, item) }

func (q *DumbQueue) pop() State {
	item := (*q)[len(*q)-1]
	*q = (*q)[0 : len(*q)-1]
	return item
}
