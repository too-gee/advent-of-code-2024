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

	// part 1
	regions := readInput(fileName).regions()

	totalPrice := 0
	for _, region := range regions {
		area, perimeter := region.grid.measure()
		totalPrice += area * perimeter
		fmt.Printf("A region of %s plants with price %d * %d = %d.\n", region.plantType, area, perimeter, area*perimeter)
	}

	fmt.Printf("Total price is %d.\n", totalPrice)

	// part 2
	regions = readInput(fileName).regions()

	bulkPrice := 0
	for _, region := range regions {
		area, _ := region.grid.measure()
		sides := region.grid.sideCount()

		for _, void := range region.voids {
			sides += void.sideCount()
		}

		bulkPrice += area * sides
		fmt.Printf("A region of %s plants with price %d * %d = %d.\n", region.plantType, area, sides, area*sides)
	}

	fmt.Printf("Bulk price is %d.\n", bulkPrice)
}

func readInput(filePath string) Grid {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening %s", filePath)
		return nil
	}
	defer file.Close()

	grid := Grid{}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		row := strings.Split(line, "")

		grid = append(grid, row)
	}

	return grid
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

type Coord struct {
	x int
	y int
}

type Grid [][]string

type Region struct {
	plantType string
	grid      Grid
	voids     []Grid
}

func (g Grid) equalTo(other Grid) bool {
	if g.width() != other.width() || g.height() != other.height() {
		return false
	}

	for y := range g {
		for x := range g[y] {
			if g[y][x] != other[y][x] {
				return false
			}
		}
	}

	return true
}

func (g Grid) blank() bool {
	for y := range g {
		for x := range g[y] {
			if g[y][x] != "." {
				return false
			}
		}
	}

	return true
}

func (g Grid) full() bool {
	for y := range g {
		for x := range g[y] {
			if g[y][x] == "." {
				return false
			}
		}
	}

	return true
}

func (g Grid) invert() Grid {
	tmp := makeGrid(g.width(), g.height())

	for y := range g {
		for x := range g[y] {
			if g[y][x] == "x" {
				tmp[y][x] = "."
			} else {
				tmp[y][x] = "x"
			}
		}
	}

	return tmp
}

func (g Grid) draw() {
	g.drawWithMarkers(map[Coord]string{})
}

func (g Grid) drawWithMarkers(markers map[Coord]string) {
	output := ""

	for y := -2; y < g.height()+2; y++ {
		for x := -2; x < g.width()+2; x++ {
			loc := Coord{x, y}

			// print the border
			if x == -2 || x == g.width()+1 || y == -2 || y == g.height()+1 {
				output += "██"
				continue
			}

			// print a buffer
			if x == -1 || x == g.width() || y == -1 || y == g.height() {
				output += "　"
				continue
			}

			// print a loc marker
			if marker, ok := markers[loc]; ok {
				output += marker
				continue
			}

			// print the rest of the grid
			if g.isInGrid(loc) {
				switch g[y][x] {
				case ".":
					output += "　"
				case "@":
					output += "⏺️ "
				case "x":
					output += "🟦"
				case "A":
					output += "❤️"
				case "B":
					output += "💛"
				case "C":
					output += "💚"
				case "D":
					output += "💙"
				case "E":
					output += "🤎"
				case "F":
					output += "🩶"
				case "G":
					output += "🟠"
				case "H":
					output += "🟣"
				case "I":
					output += "⚫"
				case "J":
					output += "⚪"
				case "K":
					output += "🟥"
				case "L":
					output += "🟨"
				case "M":
					output += "🟩"
				case "N":
					output += "🟦"
				case "O":
					output += "🟫"
				default:
					output += g[y][x] + g[y][x]
				}
			}
		}
		output += "\n"
	}

	fmt.Print(output)

}

func (g Grid) width() int {
	return len(g[0])
}

func (g Grid) height() int {
	return len(g)
}

func (g Grid) copy() Grid {
	tmp := makeGrid(g.width(), g.height())

	for y := range g {
		copy(tmp[y], g[y])
	}

	return tmp
}

func (g Grid) measure() (int, int) {
	area := 0
	perimeter := 0

	visited := makeGrid(g.width(), g.height())
	toVisit := []Coord{g.firstCell("x")}

	for i := 0; i < len(toVisit); i++ {
		loc := toVisit[i]

		if visited[loc.y][loc.x] == "x" {
			continue
		}

		neighbors := g.neighborCoords(loc)

		visited[loc.y][loc.x] = "x"
		toVisit = append(toVisit, neighbors...)

		area += 1
		perimeter += (4 - len(neighbors))
	}

	return area, perimeter
}

func (g Grid) firstCell(value string) Coord {

	for y, row := range g {
		for x := range row {
			if g[y][x] == value {

				return Coord{x: x, y: y}
			}
		}
	}

	return Coord{x: -1, y: -1}
}

func (g Grid) regions() []Region {
	regions := []Region{}

	floods := g.floods()

	for _, flood := range floods {
		plantTypeCoord := flood.firstCell("x")
		plantType := g[plantTypeCoord.y][plantTypeCoord.x]

		regions = append(regions, Region{plantType: plantType, grid: flood, voids: flood.voids()})
	}

	return regions
}

func (g Grid) floods() []Grid {
	coveredGrid := makeGrid(len(g[0]), len(g))

	floods := []Grid{}

	for y, row := range coveredGrid {
		for x, covered := range row {
			if covered == "x" {
				continue
			}

			flood := g.flood(Coord{x: x, y: y})

			area, _ := flood.measure()
			if area > 0 {
				floods = append(floods, flood)
			}

			for i, tmpRow := range flood {
				for j, tmpCell := range tmpRow {
					if tmpCell == "x" {
						coveredGrid[i][j] = "x"
					}
				}
			}
		}
	}

	return floods
}

func (g Grid) neighborCoords(loc Coord) []Coord {
	neighbors := []Coord{}
	value := g[loc.y][loc.x]

	offsets := []Coord{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}

	for _, offset := range offsets {
		candidate := Coord{x: loc.x + offset.x, y: loc.y + offset.y}
		if g.isInGrid(candidate) && g[candidate.y][candidate.x] == value {
			neighbors = append(neighbors, candidate)
		}
	}

	return neighbors
}

func (g Grid) isInGrid(l Coord) bool {
	return l.x >= 0 &&
		l.x < len(g[0]) &&
		l.y >= 0 &&
		l.y < len(g)
}

func (g Grid) flood(loc Coord) Grid {

	plantType := g[loc.y][loc.x]

	regionGrid := makeGrid(g.width(), g.height())
	toVisit := []Coord{loc}

	for {
		newVisits := false

		for i := 0; i < len(toVisit); i++ {
			start := toVisit[i]

			// skipping because we've been here
			if regionGrid[start.y][start.x] == "x" {
				continue
			}

			// skipping because this is another region
			if g[start.y][start.x] != plantType {
				continue
			}

			// mark this as part of the region and let the loop continue
			// one more time
			regionGrid[start.y][start.x] = "x"
			newVisits = true

			// try visiting neighbors
			toVisit = append(toVisit, g.neighborCoords(start)...)
		}

		if !newVisits {
			break
		}
	}

	return regionGrid
}

func (g *Grid) forward(loc *Coord) {
	if oldLoc := (*g).firstCell("@"); g.isInGrid(oldLoc) {
		(*g)[oldLoc.y][oldLoc.x] = "x"
	}

	(*g)[(*loc).y][(*loc).x] = "@"

	(*loc).y--
}

func (g *Grid) rotate(loc *Coord, dir string) {
	// mark the pivot point
	prevValue := ""

	if (*loc).x > -1 && (*loc).y > -1 {
		prevValue = (*g)[(*loc).y][(*loc).x]
		(*g)[(*loc).y][(*loc).x] = "!"
	}

	// rotate the grid
	width, height := (*g).width(), (*g).height()
	result := makeGrid(width, height)

	for y := 0; y < width; y++ {
		for x := 0; x < height; x++ {
			if dir == "L" {
				result[height-1-x][y] = (*g)[y][x]
			}

			if dir == "R" {
				result[x][height-1-y] = (*g)[y][x]
			}
		}
	}

	// find the new pivot point and unmark it
	if (*loc).x > -1 && (*loc).y > -1 {
		*loc = result.firstCell("!")

		result[(*loc).y][(*loc).x] = prevValue
	}

	*g = result
}

func (g Grid) currentView(loc Coord) Grid {
	window := makeGrid(3, 3)

	for y := loc.y - 1; y <= loc.y+1; y++ {
		for x := loc.x - 1; x <= loc.x+1; x++ {
			if g.isInGrid(Coord{x: x, y: y}) {
				window[y-loc.y+1][x-loc.x+1] = g[y][x]
			}
		}
	}

	return window
}

func (g Grid) sideCount() int {
	// track our progress around the perimeter
	moves := []string{}

	// starting conditions
	grid := g.copy()
	moveId := ""
	direction := Direction("N")
	loc := g.firstCell("x")
	prevLoc := loc

	// go until we get a repeat
	for {
		view := grid.currentView(loc)
		move := view.nextMove()

		moveId = fmt.Sprintf("loc: %d,%d  dir: %s  move: %s  view: %v", loc.x, loc.y, direction, move, view)

		// we got a repeat, stop here
		if slices.Contains(moves, moveId) {
			break
		}

		moves = append(moves, moveId)

		if move == "F" {
			grid.forward(&loc)
		}

		if move == "R" {
			grid.rotate(&loc, "R")
			direction.rotate("R")
		}

		if move == "L" {
			grid.rotate(&loc, "L")
			direction.rotate("L")
		}

		if prevLoc != loc {
			prevLoc = loc
		}
	}

	// start from the first occurrence of the repeated move and count the turns
	first := slices.Index(moves, moveId)
	totalTurns := 0

	for i := first; i < len(moves); i++ {
		if strings.Contains(moves[i], "R") || strings.Contains(moves[i], "L") {
			totalTurns++
		}
	}

	return totalTurns
}

func (g Grid) nextMove() string {
	if g[0][0] == "." && g[0][1] == "." && g[0][2] == "." {
		if g[1][0] == "." {
			return "L"
		} else {
			return "R"
		}
	}

	if g[0][0] == "." && g[0][1] == "." && g[1][0] == "." {
		return "L"
	}

	if g[1][0] == "@" && g[0][1] == "x" {
		return "F"
	}

	if g[1][0] != "." && g[2][0] == "." && g[2][1] == "@" {
		return "R"
	}

	return "F"
}

type Direction string

func (d *Direction) rotate(dir string) {
	dirs := []string{"N", "E", "S", "W"}

	i := 0

	if dir == "L" {
		i = 1
	} else {
		i = -1
	}

	currentIndex := slices.Index(dirs, string(*d))
	newIndex := int(math.Abs(float64((currentIndex + i) % 4)))

	*d = Direction(dirs[newIndex])
}

func (g Grid) voids() []Grid {
	voidCells1 := g.voidMiddles()
	g.rotate(&Coord{-1, -1}, "R")
	voidCells2 := g.voidMiddles()
	voidCells2.rotate(&Coord{-1, -1}, "L")

	voidCells := makeGrid(g.width(), g.height())

	for y := range voidCells {
		for x := range voidCells[y] {
			if voidCells1[y][x] == "x" && voidCells2[y][x] == "x" {
				voidCells[y][x] = "x"
			}
		}
	}

	tmp := voidCells.floods()
	grids := []Grid{}

	for _, flood := range tmp {
		if !flood.equalTo(voidCells.invert()) && !flood.full() && !flood.blank() {
			grids = append(grids, flood)
		}
	}

	return grids
}

func (g Grid) voidMiddles() Grid {
	voidCells := makeGrid(g.width(), g.height())

	for y, row := range g {
		first := slices.Index(row, "x")
		slices.Reverse(row)
		last := len(row) - slices.Index(row, "x")
		slices.Reverse(row)

		if first == -1 || last == -1 {
			continue
		}

		for x := first + 1; x < last; x++ {
			if g[y][x] == "." {
				voidCells[y][x] = "x"
			}
		}
	}

	return voidCells
}
