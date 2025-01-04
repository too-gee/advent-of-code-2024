package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
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

type Grid [][]string

type Region struct {
	plantType string
	grid      Grid
}

func (g Grid) draw() {
	g.drawWithMarkers(map[shared.Coord]string{})
}

func (g Grid) drawWithMarkers(markers map[shared.Coord]string) {
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
			loc := shared.Coord{x, y}

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
	toVisit := []shared.Coord{g.firstCell("x")}

	for i := 0; i < len(toVisit); i++ {
		loc := toVisit[i]

		if visited[loc.Y][loc.X] == "x" {
			continue
		}

		neighbors := g.neighborCoords(loc)

		visited[loc.Y][loc.X] = "x"
		toVisit = append(toVisit, neighbors...)

		area += 1
		perimeter += (4 - len(neighbors))
	}

	return area, perimeter
}

func (g Grid) firstCell(value string) shared.Coord {

	for y, row := range g {
		for x := range row {
			if g[y][x] == value {

				return shared.Coord{X: x, Y: y}
			}
		}
	}

	return shared.Coord{X: -1, Y: -1}
}

func (g Grid) regions() []Region {
	regions := []Region{}

	floods := g.floods()

	for _, flood := range floods {
		plantTypeCoord := flood.firstCell("x")
		plantType := g[plantTypeCoord.Y][plantTypeCoord.X]

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

			flood := g.flood(shared.Coord{X: x, Y: y})

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

func (g Grid) neighborCoords(loc shared.Coord) []shared.Coord {
	neighbors := []shared.Coord{}
	value := g[loc.Y][loc.X]

	offsets := []shared.Coord{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}

	for _, offset := range offsets {
		candidate := shared.Coord{X: loc.X + offset.X, Y: loc.Y + offset.Y}
		if g.isInGrid(candidate) && g[candidate.Y][candidate.X] == value {
			neighbors = append(neighbors, candidate)
		}
	}

	return neighbors
}

func (g Grid) isInGrid(l shared.Coord) bool {
	return l.X >= 0 &&
		l.X < len(g[0]) &&
		l.Y >= 0 &&
		l.Y < len(g)
}

func (g Grid) flood(loc shared.Coord) Grid {

	plantType := g[loc.Y][loc.X]

	regionGrid := makeGrid(g.width(), g.height())
	toVisit := []shared.Coord{loc}

	for {
		newVisits := false

		for i := 0; i < len(toVisit); i++ {
			start := toVisit[i]

			// skipping because we've been here
			if regionGrid[start.Y][start.X] == "x" {
				continue
			}

			// skipping because this is another region
			if g[start.Y][start.X] != plantType {
				continue
			}

			// mark this as part of the region and let the loop continue
			// one more time
			regionGrid[start.Y][start.X] = "x"
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

			if g.isInGrid(shared.Coord{X: x, Y: y - 1}) {
				upper = g[y-1][x]
			} else {
				upper = "."
			}

			if g.isInGrid(shared.Coord{X: x, Y: y}) {
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
