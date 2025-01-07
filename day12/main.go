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
	regions := readInput(fileName)

	totalPrice := PartOne(regions)
	fmt.Printf("Total price is %d.\n", totalPrice)

	// part 2
	regions = readInput(fileName)

	bulkPrice := PartTwo(regions)

	fmt.Printf("Bulk price is %d.\n", bulkPrice)
}

func PartOne(regions []Region) int {
	totalPrice := 0
	for _, region := range regions {
		area, perimeter := region.measure()
		totalPrice += area * perimeter
		fmt.Printf("A region of %s plants with price %d * %d = %d.\n", region.plantType, area, perimeter, area*perimeter)
	}
	return totalPrice
}

func PartTwo(regions []Region) int {
	bulkPrice := 0
	for _, region := range regions {
		area, _ := region.measure()
		sides := region.dumbSideCount()

		bulkPrice += area * sides
		fmt.Printf("A region of %s plants with price %d * %d = %d.\n", region.plantType, area, sides, area*sides)
	}
	return bulkPrice
}

func readInput(filePath string) []Region {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening %s", filePath)
		return nil
	}
	defer file.Close()

	grid := shared.Grid{}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		row := strings.Split(line, "")

		grid = append(grid, row)
	}

	tmpRegion := Region{plantType: "?", Grid: grid}

	return tmpRegion.regions()
}

func makeGrid(width int, height int) shared.Grid {
	tmp := make(shared.Grid, height)

	for i := range tmp {
		tmp[i] = make([]string, width)

		for j := range tmp[i] {
			tmp[i][j] = "."
		}
	}

	return tmp
}

type Region struct {
	plantType string
	shared.Grid
}

func (r Region) draw() {
	r.drawWithMarkers(map[shared.Coord]string{})
}

func (r Region) drawWithMarkers(markers map[shared.Coord]string) {
	yMin, yMax := r.Height()-1, 0
	xMin, xMax := r.Width()-1, 0

	for y := range r.Height() {
		for x := range r.Width() {
			if r.Grid[y][x] != "." {
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
			loc := shared.Coord{X: x, Y: y}

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
			if r.Contains(loc) {
				switch r.Grid[y][x] {
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
					output += r.Grid[y][x] + r.Grid[y][x]
				}
			}
		}
		output += "\n"
	}

	fmt.Print(output)
}

func (r Region) measure() (int, int) {
	area := 0
	perimeter := 0

	visited := makeGrid(r.Width(), r.Height())
	toVisit := []shared.Coord{r.LocationOf("x")}

	for i := 0; i < len(toVisit); i++ {
		loc := toVisit[i]

		if visited[loc.Y][loc.X] == "x" {
			continue
		}

		neighbors := r.neighborCoords(loc)

		visited[loc.Y][loc.X] = "x"
		toVisit = append(toVisit, neighbors...)

		area += 1
		perimeter += (4 - len(neighbors))
	}

	return area, perimeter
}

func (r Region) regions() []Region {
	coveredGrid := makeGrid(r.Width(), r.Height())

	floods := []Region{}

	for y, row := range coveredGrid {
		for x, covered := range row {
			if covered == "x" {
				continue
			}

			flood := r.flood(shared.Coord{X: x, Y: y})

			area, _ := flood.measure()
			if area > 0 {
				floods = append(floods, flood)
			}

			for i, tmpRow := range flood.Grid {
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

func (r Region) neighborCoords(loc shared.Coord) []shared.Coord {
	neighbors := []shared.Coord{}
	value := r.Grid[loc.Y][loc.X]

	for _, candidate := range loc.Neighbors() {
		if r.Contains(candidate) && r.Grid[candidate.Y][candidate.X] == value {
			neighbors = append(neighbors, candidate)
		}
	}

	return neighbors
}

func (r Region) flood(loc shared.Coord) Region {
	plantType := r.Grid[loc.Y][loc.X]

	regionGrid := makeGrid(r.Width(), r.Height())
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
			if r.Grid[start.Y][start.X] != plantType {
				continue
			}

			// mark this as part of the region and let the loop continue
			// one more time
			regionGrid[start.Y][start.X] = "x"
			newVisits = true

			// try visiting neighbors
			toVisit = append(toVisit, r.neighborCoords(start)...)
		}

		if !newVisits {
			break
		}
	}

	return Region{plantType: plantType, Grid: regionGrid}
}

func (r *Region) rotate(dir string) {
	// rotate the grid
	width, height := (*r).Grid.Width(), (*r).Grid.Height()
	result := makeGrid(width, height)

	for y := 0; y < width; y++ {
		for x := 0; x < height; x++ {
			if dir == "L" {
				result[height-1-x][y] = (*r).Grid[y][x]
			}

			if dir == "R" {
				result[x][height-1-y] = (*r).Grid[y][x]
			}
		}
	}

	(*r).Grid = result
}

func (r Region) dumbSideCount() int {
	runs := r.dumbRunCount()
	r.rotate("R")
	runs += r.dumbRunCount()

	return runs
}

func (r Region) dumbRunCount() int {
	runs := 0

	for y := 0; y < r.Height()+1; y++ {
		// collapse each row into unique sections
		row := []string{}

		for x := 0; x < r.Width(); x++ {
			upper := ""
			lower := ""

			if r.Contains(shared.Coord{X: x, Y: y - 1}) {
				upper = r.Grid[y-1][x]
			} else {
				upper = "."
			}

			if r.Contains(shared.Coord{X: x, Y: y}) {
				lower = r.Grid[y][x]
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
