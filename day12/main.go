package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
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
		sides := region.grid.dumbSideCount()

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
}

func (g Grid) draw() {
	g.drawWithMarkers(map[Coord]string{})
}

func (g Grid) drawWithMarkers(markers map[Coord]string) {
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

		regions = append(regions, Region{plantType: plantType, grid: flood})
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

func (g *Grid) rotate(dir string) {
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

	*g = result
}

func (g Grid) dumbSideCount() int {
	runs := g.dumbRunCount()
	g.rotate("R")
	runs += g.dumbRunCount()

	return runs
}

func (g Grid) dumbRunCount() int {
	runs := 0

	for y := 0; y < g.height()+1; y++ {
		// collapse each row into unique sections
		row := []string{}

		for x := 0; x < g.width(); x++ {
			upper := ""
			lower := ""

			if g.isInGrid(Coord{x: x, y: y - 1}) {
				upper = g[y-1][x]
			} else {
				upper = "."
			}

			if g.isInGrid(Coord{x: x, y: y}) {
				lower = g[y][x]
			} else {
				lower = "."
			}

			currentCells := upper + lower
			if len(row) == 0 || row[len(row)-1] != currentCells {
				row = append(row, currentCells)
			}
		}

		// count the number of "edges"
		for _, section := range row {
			if section[0] != section[1] {
				runs++
			}
		}
	}

	return runs
}
