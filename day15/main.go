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

	displayMap := !strings.HasSuffix(fileName, "input.txt")

	grid, moves := readInput(fileName)
	warehouse := Warehouse{grid: grid, direction: N}

	// part 1
	for _, move := range moves {
		if displayMap {
			warehouse.draw()
			fmt.Printf("Move: %v\n", move)
		}

		switch move {
		case UP:
			warehouse.turnToFace(E)
		case RIGHT:
			warehouse.turnToFace(N)
		case DOWN:
			warehouse.turnToFace(W)
		case LEFT:
			warehouse.turnToFace(S)
		}

		robot := warehouse.grid.find(ROBOT)
		start := robot.x
		space := float64(slices.Index(warehouse.grid[robot.y][start:], EMPTY))
		wall := float64(slices.Index(warehouse.grid[robot.y][start:], WALL))
		end := int(math.Max(math.Min(space, wall), 0)) + start + 1

		row := warehouse.grid[robot.y][start:end]
		newRow := []string{}

		for _, cell := range row {
			if cell == EMPTY {
				newRow = append(newRow, EMPTY)
			}
		}

		for _, cell := range row {
			if cell != EMPTY {
				newRow = append(newRow, cell)
			}
		}

		for x := start; x < end; x++ {
			warehouse.grid[robot.y][x] = newRow[x-start]
		}
	}

	warehouse.draw()

	coordSum := 0
	for y := range warehouse.grid.height() {
		for x := range warehouse.grid.width() {
			if warehouse.grid[y][x] == BOX {
				coordSum += (100 * y) + x
			}
		}
	}

	fmt.Printf("The GPS coordinate sum is %d\n", coordSum)
}

func readInput(filePath string) (Grid, []string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening %s", filePath)
		return nil, nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	grid := Grid{}
	moves := []string{}

	readMoves := false

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			readMoves = true
			continue
		}

		if readMoves {
			lineMoves := strings.Split(line, "")
			moves = append(moves, lineMoves...)
		} else {
			lineGrid := strings.Split(line, "")
			grid = append(grid, lineGrid)
		}
	}

	return grid, moves
}

func makeGrid(width int, height int) Grid {
	tmp := make(Grid, height)

	for i := range tmp {
		tmp[i] = make([]string, width)

		for j := range tmp[i] {
			tmp[i][j] = "."
		}
	}

	return tmp
}

const WALL = "#"
const EMPTY = "."
const BOX = "O"
const ROBOT = "@"

const N = 0
const E = 1
const S = 2
const W = 3

const CCW = 0
const CW = 1

const UP = "^"
const RIGHT = ">"
const DOWN = "v"
const LEFT = "<"

type Coord struct {
	x int
	y int
}

type Grid [][]string

type Warehouse struct {
	grid      Grid
	direction int
}

func (w *Warehouse) turnToFace(dir int) {
	if dir == w.direction {
		return
	}

	var turns []int

	if (dir == N && w.direction == E) ||
		(dir == E && w.direction == S) ||
		(dir == S && w.direction == W) ||
		(dir == W && w.direction == N) {
		turns = []int{CCW}
	} else if (dir == N && w.direction == W) ||
		(dir == W && w.direction == S) ||
		(dir == S && w.direction == E) ||
		(dir == E && w.direction == N) {
		turns = []int{CW}
	} else {
		turns = []int{CCW, CCW}
	}

	for _, turn := range turns {
		(*w).grid.rotate(turn)
	}

	(*w).direction = dir
}

func (w *Warehouse) draw() {
	(*w).turnToFace(N)
	(*w).grid.draw()
}

func (g Grid) isInGrid(l Coord) bool {
	return l.x >= 0 &&
		l.x < g.width() &&
		l.y >= 0 &&
		l.y < g.height()
}

func (g Grid) width() int {
	return len(g[0])
}

func (g Grid) height() int {
	return len(g)
}

func (g *Grid) rotate(dir int) {
	// rotate the grid
	width, height := (*g).width(), (*g).height()
	result := makeGrid(width, height)

	for y := 0; y < width; y++ {
		for x := 0; x < height; x++ {
			if dir == CCW {
				result[height-1-x][y] = (*g)[y][x]
			}

			if dir == CW {
				result[x][height-1-y] = (*g)[y][x]
			}
		}
	}

	*g = result
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

func (g Grid) draw() {
	yMin, yMax := g.height()-1, 0
	xMin, xMax := g.width()-1, 0

	for y := range g {
		for x := range g[y] {
			if g[y][x] != "." {
				yMin = int(math.Min(float64(yMin), float64(y)))
				yMax = int(math.Max(float64(yMax), float64(y)))
				xMin = int(math.Min(float64(xMin), float64(x)))
				xMax = int(math.Max(float64(xMax), float64(x)))
			}
		}
	}

	output := ""

	borderThickness := 1
	bufferThickness := 1

	for y := yMin - borderThickness - bufferThickness; y <= yMax+borderThickness+bufferThickness; y++ {
		for x := xMin - borderThickness - bufferThickness; x <= xMax+borderThickness+bufferThickness; x++ {
			loc := Coord{x, y}

			// print the border
			if x >= xMin-borderThickness-bufferThickness && x < xMin-bufferThickness ||
				(y >= yMin-borderThickness-bufferThickness && y < yMin-bufferThickness) ||
				(x > xMax+bufferThickness && x <= xMax+borderThickness+bufferThickness) ||
				(y > yMax+bufferThickness && y <= yMax+borderThickness+bufferThickness) {
				output += "██"
				continue
			}

			// print a buffer
			if (x >= xMin-borderThickness && x < xMin) ||
				(y >= yMin-borderThickness && y < yMin) ||
				(x > xMax && x <= xMax+bufferThickness) ||
				(y > yMax && y <= yMax+bufferThickness) {
				output += "　"
				continue
			}

			// print the grid
			if g.isInGrid(loc) {
				switch g[y][x] {
				case ROBOT:
					output += "🤖"
				case WALL:
					output += "⬛"
				case BOX:
					output += "⏺️ "
				default:
					output += "　"
				}
			}
		}
		output += "\n"
	}

	fmt.Print(output)
}
