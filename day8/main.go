package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
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

	grid := readInput(fileName)

	// part 1
	frequencies := getUniqueCharacters(grid)
	antennaLocations := map[string]coordList{}
	antiNodeLocations := coordList{}

	for _, frequency := range frequencies {
		antennaLocations[frequency] = coordList{}
		for y, row := range grid {
			for x, cell := range row {
				if cell == frequency {
					antennaLocations[frequency] = append(antennaLocations[frequency], shared.Coord{X: x, Y: y})
				}
			}
		}
	}

	for _, antennas := range antennaLocations {
		pairs := pairs(makeRange(0, len(antennas)-1))

		for _, pair := range pairs {
			antennaA := antennas[pair[0]]
			antennaB := antennas[pair[1]]

			antinode := shared.Coord{X: (2 * antennaA.X) - antennaB.X, Y: (2 * antennaA.Y) - antennaB.Y}

			if isInGrid(antinode, grid) && !slices.Contains(antiNodeLocations, antinode) {
				antiNodeLocations = append(antiNodeLocations, antinode)
			}
		}
	}

	fmt.Printf("There are %d unique locations that contain an antinode\n", len(antiNodeLocations))

	// part 2
	for _, antennas := range antennaLocations {
		pairs := pairs(makeRange(0, len(antennas)-1))

		for _, pair := range pairs {
			antennaA := antennas[pair[0]]
			antennaB := antennas[pair[1]]

			if !slices.Contains(antiNodeLocations, antennaA) {
				antiNodeLocations = append(antiNodeLocations, antennaA)
			}

			rise := antennaA.Y - antennaB.Y
			run := antennaA.X - antennaB.X
			currentX := antennaA.X
			currentY := antennaA.Y

			for {
				antinode := shared.Coord{X: currentX + run, Y: currentY + rise}

				if !isInGrid(antinode, grid) {
					break
				}

				if !slices.Contains(antiNodeLocations, antinode) {
					antiNodeLocations = append(antiNodeLocations, antinode)
				}

				currentX += run
				currentY += rise
			}

		}
	}

	fmt.Printf("There are %d unique locations that contain a resonant antinode\n", len(antiNodeLocations))
}

func readInput(filePath string) [][]string {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening %s", filePath)
		return nil
	}
	defer file.Close()

	grid := [][]string{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		row := strings.Split(line, "")

		grid = append(grid, row)
	}

	return grid
}

func getUniqueCharacters(grid [][]string) []string {
	uniqueChars := make(map[string]bool)

	ignoredChars := []string{"."}

	for _, row := range grid {
		for _, cell := range row {
			if !slices.Contains(ignoredChars, cell) {
				uniqueChars[cell] = true
			}
		}
	}

	var uniqueCharsSlice []string

	for char := range uniqueChars {
		uniqueCharsSlice = append(uniqueCharsSlice, char)
	}

	return uniqueCharsSlice
}

func pairs(options []int) [][]int {
	permutations := [][]int{}

	for i := 0; i < int(math.Pow(float64(len(options)), float64(2))); i++ {
		binary := strconv.FormatInt(int64(i), len(options))
		binaryStr := fmt.Sprintf("%02s", binary)
		binarySlice := strings.Split(binaryStr, "")

		intSlice := []int{}

		for _, str := range binarySlice {
			intVal, _ := strconv.Atoi(str)
			intSlice = append(intSlice, options[intVal])
		}

		if intSlice[0] != intSlice[1] {
			permutations = append(permutations, intSlice)
		}
	}

	return permutations
}

func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

func isInGrid(l shared.Coord, grid [][]string) bool {
	if l.X >= 0 && l.X < len(grid) && l.Y >= 0 && l.Y < len(grid[0]) {
		return true
	}
	return false
}

type coordList []shared.Coord
