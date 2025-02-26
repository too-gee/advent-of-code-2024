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
	maze.Draw(map[string]string{"S": "🟢", "E": "🔴"}, nil)
	result := Part1(maze, 100)
	fmt.Println(result)

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

func Part1(g shared.Grid, savings int) int {
	return len(Solve(g, savings))
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

	fmt.Printf("Non-cheat length: %d, minSavings: %d, maxLength: %d\n", totalLength, savings, maxLength)

	start := g.LocationOf(START)
	end := g.LocationOf(END)

	queue := DumbQueue{}
	queue.push(State{
		cheatStart: -1,
		path:       []shared.Coord{start},
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

		//fmt.Printf("Current: %v, len: %d, cheatsLeft: %d\n", current.Loc(), len(current.path), current.CheatsLeft())
		//g.Draw(map[string]string{"S": "🟢", "E": "🔴"}, map[string][]shared.Coord{"⏺️ ": current.path})
		cheatsLeft := current.CheatsLeft()

		for _, neighbor := range g.Neighbors(current.Loc(), nil) {
			if slices.Contains(current.path, neighbor) {
				continue
			}

			newCheatStart := current.cheatStart

			if g.At(neighbor) == WALL {
				if cheatsLeft == 0 {
					continue
				}

				newCheatStart = len(current.path)
			}

			queue.push(State{
				cheatStart: newCheatStart,
				path:       CopyAppend(current.path, neighbor),
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

type State struct {
	cheatStart int
	path       []shared.Coord
}

func (s State) Loc() shared.Coord {
	return s.path[len(s.path)-1]
}

func (s State) CheatsLeft() int {
	if s.cheatStart == -1 {
		return 2
	}

	return int(math.Max(0, float64(s.cheatStart-len(s.path)+1)))
}

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
