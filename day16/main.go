package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"slices"
	"strings"
)

func main() {
	var fileName string

	if len(os.Args) == 2 {
		fileName = os.Args[1]
	} else {
		fileName = "input.txt"
	}

	//displayMap := !strings.HasSuffix(fileName, "input.txt")

	// part 1
	maze := readInput(fileName)

	// Fill dead ends until there are no more
	for {
		deadEndsFound := 0
		for y := 0; y < len(maze); y++ {
			for x := 0; x < len(maze[0]); x++ {
				if maze.isDeadEnd(Coord{x: x, y: y}) {
					deadEndsFound += 1
					maze.fillDeadEnd(Coord{x: x, y: y})
				}
			}
		}
		if deadEndsFound == 0 {
			fmt.Println("No dead ends detected")
			break
		}

		for _, row := range maze {
			for _, cell := range row {
				fmt.Print(cell)
			}
			fmt.Println()
		}

		fmt.Println("Filled dead ends")
	}

	mazeStart := maze.find(START)
	fakePrev := Coord{x: mazeStart.x - 1, y: mazeStart.y}
	cost := maze.getPath([]Coord{fakePrev, maze.find(START)})
	fmt.Println(cost)

}

func readInput(filePath string) Grid {
	file, _ := os.Open(filePath)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	grid := Grid{}

	for scanner.Scan() {
		line := scanner.Text()

		gridRow := strings.Split(line, "")
		grid = append(grid, gridRow)
	}

	return grid
}

const FILL = "="
const WALL = "#"
const EMPTY = "."
const START = "S"
const END = "E"

type Coord struct {
	x int
	y int
}

func (c Coord) neighbors() []Coord {
	return []Coord{
		{x: c.x - 1, y: c.y},
		{x: c.x + 1, y: c.y},
		{x: c.x, y: c.y - 1},
		{x: c.x, y: c.y + 1},
	}
}

type Grid [][]string

func (g Grid) isPassable(loc Coord) bool {
	return g[loc.y][loc.x] != WALL && g[loc.y][loc.x] != FILL
}

func (g Grid) getPath(trail []Coord) int { // this segment costs 1 point because we moved to start here
	cost := 1

	var visitNext []Coord
	loc := trail[len(trail)-1]

	// we made it!
	if g[loc.y][loc.x] == END {
		return 0
	}

	// if our trail is longer than 2, check to see if we
	// turned, which costs 1000 points
	if len(trail) > 2 {
		curDir := direction(
			trail[len(trail)-2],
			trail[len(trail)-1],
		)

		prevDir := direction(
			trail[len(trail)-3],
			trail[len(trail)-2],
		)

		if curDir != prevDir {
			fmt.Println("Turned")
			cost += 1000
		}
	}

	visitNext = []Coord{}

	for _, neighbor := range loc.neighbors() {
		if !slices.Contains(trail, neighbor) &&
			g[neighbor.y][neighbor.x] != WALL &&
			g[neighbor.y][neighbor.x] != FILL {
			visitNext = append(visitNext, neighbor)
		}
	}

	// if we can't find a way forward, return a large number
	// so that this path is not chosen
	if len(visitNext) == 0 {
		return math.MaxInt32
	}

	minCost := math.MaxInt32

	for i := range visitNext {
		cost := g.getPath(append(trail, visitNext[i]))
		if cost < minCost {
			minCost = cost
		}
	}

	return cost + minCost
}

func (g Grid) find(value string) Coord {

	for y, row := range g {
		for x := range row {
			if g[y][x] == value {

				return Coord{x: x, y: y}
			}
		}
	}

	return Coord{x: -1, y: -1}
}

func direction(start Coord, end Coord) string {
	if start.x > end.x {
		return "E"
	} else if start.x < end.x {
		return "W"
	} else if start.y > end.y {
		return "N"
	} else if start.y < end.y {
		return "S"
	}

	panic("Invalid direction")
}

func (g *Grid) fillDeadEnd(loc Coord) {
	for {
		if (*g)[loc.y][loc.x] == START || (*g)[loc.y][loc.x] == END {
			return
		}

		toVisit := []Coord{}

		for _, neighbor := range loc.neighbors() {
			if (*g).isPassable(neighbor) {
				toVisit = append(toVisit, neighbor)
			}
		}

		if len(toVisit) != 1 {
			break
		}

		(*g)[loc.y][loc.x] = FILL
		loc = toVisit[0]
	}
}

func (g Grid) isDeadEnd(loc Coord) bool {
	if !g.isPassable(loc) ||
		g[loc.y][loc.x] == START ||
		g[loc.y][loc.x] == END {
		return false
	}

	options := 0
	for _, neighbor := range loc.neighbors() {
		if g.isPassable(neighbor) {
			options++
		}
	}
	return options == 1
}
