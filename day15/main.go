package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"slices"
	"sort"
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

	// part 1
	grid, moves := readInput(fileName)
	gpsValue := PartOne(grid, moves)
	fmt.Printf("The GPS coordinate sum is %d\n", gpsValue)

	// part 2
	grid, moves = readInput(fileName)
	gpsValue = PartTwo(grid, moves)
	fmt.Printf("The GPS coordinate sum is %d\n", gpsValue)
}

func PartOne(grid shared.Grid, moves []string) int {
	warehouse := Warehouse{Grid: grid, direction: N}

	for _, move := range moves {
		warehouse.moveRobot(move)
	}

	warehouse.draw()

	return warehouse.gpsValue()
}

func PartTwo(grid shared.Grid, moves []string) int {
	warehouse := Warehouse{Grid: grid}
	warehouse.widen()

	for _, move := range moves {
		warehouse.wideMoveRobot(move)
	}
	warehouse.draw()

	return warehouse.gpsValue()
}

func readInput(filePath string) (shared.Grid, []string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening %s", filePath)
		return nil, nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	grid := shared.Grid{}
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

const WALL = "#"
const EMPTY = "."
const BOX = "O"
const ROBOT = "@"

const DBL_WALL = "##"
const DBL_BOX = "[]"
const DBL_EMPTY = ".."

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

type Warehouse struct {
	shared.Grid
	direction int
	wide      bool
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
		(*w).rotate(turn)
	}

	(*w).direction = dir
}

func (w *Warehouse) widen() {
	wideGrid := shared.Grid{}

	for y := range (*w).Height() {
		wideRow := []string{}
		for x := range (*w).Width() {
			var newItem string

			switch (*w).Grid[y][x] {
			case WALL:
				newItem = DBL_WALL
			case BOX:
				newItem = DBL_BOX
			case EMPTY:
				newItem = DBL_EMPTY
			case ROBOT:
				newItem = ROBOT + EMPTY
			}

			wideRow = append(wideRow, strings.Split(newItem, "")...)
		}
		wideGrid = append(wideGrid, wideRow)
	}

	(*w).Grid = wideGrid
	(*w).wide = true
}

func (w *Warehouse) rotate(dir int) {
	// rotate the grid
	width, height := (*w).Width(), (*w).Height()
	result := shared.MakeGrid(width, height)

	for y := 0; y < width; y++ {
		for x := 0; x < height; x++ {
			if dir == CCW {
				result[height-1-x][y] = (*w).Grid[y][x]
			}

			if dir == CW {
				result[x][height-1-y] = (*w).Grid[y][x]
			}
		}
	}

	(*w).Grid = result
}

func (w Warehouse) draw() {
	w.turnToFace(N)

	yMin, yMax := w.Height()-1, 0
	xMin, xMax := w.Width()-1, 0

	for y := range w.Height() {
		for x := range w.Width() {
			if w.Grid[y][x] != "." {
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

	var xIncr int

	if w.wide {
		xIncr = 2
	} else {
		xIncr = 1
	}

	for y := yMin - borderThickness - bufferThickness; y <= yMax+borderThickness+bufferThickness; y++ {
		for x := xMin - ((borderThickness + bufferThickness) * xIncr); x <= xMax+((borderThickness+bufferThickness)*xIncr); x += xIncr {
			loc := shared.Coord{X: x, Y: y}

			// print the border
			if x >= xMin-((borderThickness+bufferThickness)*xIncr) && x < xMin-(bufferThickness*xIncr) ||
				(y >= yMin-borderThickness-bufferThickness && y < yMin-bufferThickness) ||
				(x > xMax+(bufferThickness*xIncr) && x <= xMax+(borderThickness+bufferThickness)*xIncr) ||
				(y > yMax+bufferThickness && y <= yMax+borderThickness+bufferThickness) {
				output += "██"
				continue
			}

			// print a buffer
			if (x >= xMin-(borderThickness*xIncr) && x < xMin) ||
				(y >= yMin-borderThickness && y < yMin) ||
				(x > xMax && x <= xMax+(bufferThickness*xIncr)) ||
				(y > yMax && y <= yMax+bufferThickness) {
				output += "　"
				continue
			}

			// print the grid
			if w.Contains(loc) {
				var cell string
				if w.wide {
					cell = w.Grid[y][x] + w.Grid[y][x+1]
				} else {
					cell = w.Grid[y][x]
				}

				switch cell {
				case ROBOT:
					output += "🤖"
				case WALL:
					output += "⬛"
				case DBL_WALL:
					output += "⬛"
				case BOX:
					output += "⏺️ "
				case EMPTY:
					output += "　"
				case DBL_EMPTY:
					output += "　"
				default:
					output += strings.ReplaceAll(cell, ".", " ")
				}
			}
		}
		output += "\n"
	}

	fmt.Print(output)
}

func reversed(s []string) []string {
	newSlice := make([]string, len(s))

	for i := 0; i < len(s); i++ {
		newSlice[i] = s[len(s)-i-1]
	}

	return newSlice
}

func (w *Warehouse) moveRobot(move string) {
	switch move {
	case UP:
		(*w).turnToFace(E)
	case RIGHT:
		(*w).turnToFace(N)
	case DOWN:
		(*w).turnToFace(E)
	case LEFT:
		(*w).turnToFace(N)
	}

	var tmp []string
	robot := w.Grid.LocationOf(ROBOT)
	var start int

	if move == LEFT || move == DOWN {
		tmp = reversed(w.Grid[robot.Y])
		start = w.Grid.Width() - robot.X - 1
	} else {
		tmp = w.Grid[robot.Y]
		start = robot.X
	}

	space := float64(slices.Index(tmp[start:], EMPTY))
	wall := float64(slices.Index(tmp[start:], WALL))
	end := int(math.Max(math.Min(space, wall), 0)) + start + 1

	row := tmp[start:end]
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

	if move == LEFT || move == DOWN {
		newRow = reversed(newRow)
		tmpStart := start
		start = w.Grid.Width() - end
		end = w.Grid.Width() - tmpStart
	}

	for x := start; x < end; x++ {
		(*w).Grid[robot.Y][x] = newRow[x-start]
	}
}

func (w Warehouse) canMove(currentLoc shared.Coord, move string) []shared.Coord {
	var moveAmt shared.Coord

	switch move {
	case LEFT:
		moveAmt = shared.Coord{X: -1, Y: 0}
	case RIGHT:
		moveAmt = shared.Coord{X: 1, Y: 0}
	case UP:
		moveAmt = shared.Coord{X: 0, Y: -1}
	case DOWN:
		moveAmt = shared.Coord{X: 0, Y: 1}
	}

	movers := []shared.Coord{}
	nextLoc := shared.Coord{X: currentLoc.X + moveAmt.X, Y: currentLoc.Y + moveAmt.Y}

	// If we're moving left or right, pretend the split boxes are just normal boxes
	nextLocStr := w.Grid[nextLoc.Y][nextLoc.X]
	if (move == LEFT || move == RIGHT) && (nextLocStr == "[" || nextLocStr == "]") {
		nextLocStr = "O"
	}

	switch nextLocStr {
	// Empty space, can move
	case ".":
		return []shared.Coord{currentLoc}
	// Just a wall, must stop
	case "#":
		return nil
	// Normal width box, check the next space
	case "O":
		newMovers := w.canMove(nextLoc, move)
		if len(newMovers) == 0 {
			return nil
		}
		movers = []shared.Coord{currentLoc}
		movers = append(movers, newMovers...)
	// Double width box, check both spaces
	case "[":
		newMovers := w.canMove(nextLoc, move)
		if len(newMovers) == 0 {
			return nil
		}
		movers = []shared.Coord{currentLoc}
		movers = append(movers, newMovers...)

		rightNextLoc := shared.Coord{X: nextLoc.X + 1, Y: nextLoc.Y}
		newMovers = w.canMove(rightNextLoc, move)
		if len(newMovers) == 0 {
			return nil
		}
		movers = append(movers, newMovers...)
	case "]":
		newMovers := w.canMove(nextLoc, move)
		if len(newMovers) == 0 {
			return nil
		}
		movers = []shared.Coord{currentLoc}
		movers = append(movers, newMovers...)

		leftNextLoc := shared.Coord{X: nextLoc.X - 1, Y: nextLoc.Y}
		newMovers = w.canMove(leftNextLoc, move)
		if len(newMovers) == 0 {
			return nil
		}
		movers = append(movers, newMovers...)
	}

	// sort by y then x
	sort.Slice(movers, func(i, j int) bool {
		if movers[i].Y == movers[j].Y {
			return movers[i].X < movers[j].X
		} else {
			return movers[i].Y < movers[j].Y
		}
	})

	// only keep unique values (otherwise we get double-moves)
	movers = slices.Compact(movers)

	// reverse the order if we're moving down or right so that moves get
	// executed in the correct order
	if move == DOWN || move == RIGHT {
		for i, j := 0, len(movers)-1; i < j; i, j = i+1, j-1 {
			movers[i], movers[j] = movers[j], movers[i]
		}
	}

	return movers
}

func (w *Warehouse) wideMoveRobot(move string) {
	robot := w.Grid.LocationOf(ROBOT)
	movers := w.canMove(robot, move)

	if len(movers) > 0 {
		var moveAmt shared.Coord

		switch move {
		case LEFT:
			moveAmt = shared.Coord{X: -1, Y: 0}
		case RIGHT:
			moveAmt = shared.Coord{X: 1, Y: 0}
		case UP:
			moveAmt = shared.Coord{X: 0, Y: -1}
		case DOWN:
			moveAmt = shared.Coord{X: 0, Y: 1}
		}

		for _, a := range movers {
			b := shared.Coord{X: a.X + moveAmt.X, Y: a.Y + moveAmt.Y}

			(*w).Grid[b.Y][b.X] = (*w).Grid[a.Y][a.X]
			(*w).Grid[a.Y][a.X] = "."
		}
	}
}

func (w Warehouse) gpsValue() int {
	w.turnToFace(N)

	coordSum := 0
	for y := range w.Grid.Height() {
		for x := range w.Grid.Width() {
			if w.Grid[y][x] == "[" || w.Grid[y][x] == BOX {
				coordSum += (100 * y) + x
			}
		}
	}

	return coordSum
}
